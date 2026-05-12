package logic

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"zilo/backend/model"
	"zilo/backend/rpc/workflow/internal/svc"
	"zilo/backend/rpc/workflow/workflowpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExecuteWorkflowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewExecuteWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExecuteWorkflowLogic {
	return &ExecuteWorkflowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ExecuteWorkflowLogic) ExecuteWorkflow(in *workflowpb.ExecuteReq) (*workflowpb.ExecuteResp, error) {
	workflow, err := l.svcCtx.WorkflowModel.FindOne(l.ctx, uint64(in.WorkflowId))
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("workflow not found")
		}
		return nil, err
	}
	if workflow.OwnerId != uint64(in.UserId) {
		return nil, errors.New("no permission to execute workflow")
	}

	version, err := l.svcCtx.WorkflowVersionModel.FindOneByWorkflowIdVersion(l.ctx, uint64(in.WorkflowId), workflow.LatestVersion)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("no published version to execute")
		}
		return nil, err
	}

	// 1) 落 execution 头记录（pending）
	pending, err := l.svcCtx.ExecutionModel.Insert(l.ctx, &model.WorkflowExecutions{
		WorkflowId:  uint64(in.WorkflowId),
		VersionId:   version.Id,
		TriggerType: "manual",
		Status:      "running",
		InputJson:   sql.NullString{String: in.InputJson, Valid: in.InputJson != ""},
		StartedAt:   sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return nil, err
	}
	executionID, err := pending.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 2) 调 Eino 引擎跑 DSL
	result, runErr := l.svcCtx.Engine.Execute(l.ctx, version.DslJson, in.InputJson)
	// 如果 result 为空说明在解析阶段就失败了，直接落终态
	if result == nil {
		if updateErr := l.markExecutionFailed(uint64(executionID), runErr.Error()); updateErr != nil {
			l.Errorf("update failed execution err: %v", updateErr)
		}
		return nil, runErr
	}

	// 3) 写每个节点的执行日志（即使整体失败也要保留已执行节点的日志）
	for _, nl := range result.NodeLogs {
		_, insErr := l.svcCtx.ExecutionLogModel.Insert(l.ctx, &model.WorkflowExecutionLogs{
			ExecutionId: uint64(executionID),
			NodeId:      nl.NodeID,
			NodeType:    nl.NodeType,
			Status:      nl.Status,
			InputJson:   sql.NullString{String: nl.Input, Valid: nl.Input != ""},
			OutputJson:  sql.NullString{String: nl.Output, Valid: nl.Output != ""},
			ErrorMsg:    sql.NullString{String: nl.Error, Valid: nl.Error != ""},
			StartedAt:   sql.NullTime{Time: nl.StartedAt, Valid: true},
			EndedAt:     sql.NullTime{Time: nl.EndedAt, Valid: true},
			DurationMs:  nl.DurationMs,
		})
		if insErr != nil {
			l.Errorf("insert node log err: %v", insErr)
		}
	}

	// 4) 更新 execution 终态
	finalExec := &model.WorkflowExecutions{
		Id:          uint64(executionID),
		WorkflowId:  uint64(in.WorkflowId),
		VersionId:   version.Id,
		TriggerType: "manual",
		Status:      result.Status,
		InputJson:   sql.NullString{String: in.InputJson, Valid: in.InputJson != ""},
		OutputJson:  sql.NullString{String: result.OutputJSON, Valid: result.OutputJSON != ""},
		ErrorMsg:    sql.NullString{String: result.ErrorMsg, Valid: result.ErrorMsg != ""},
		StartedAt:   sql.NullTime{Time: result.StartedAt, Valid: true},
		EndedAt:     sql.NullTime{Time: result.EndedAt, Valid: true},
		DurationMs:  result.DurationMs,
	}
	if updateErr := l.svcCtx.ExecutionModel.Update(l.ctx, finalExec); updateErr != nil {
		l.Errorf("update execution err: %v", updateErr)
	}

	if runErr != nil {
		return &workflowpb.ExecuteResp{
			ExecutionId: executionID,
			Status:      result.Status,
		}, runErr
	}

	return &workflowpb.ExecuteResp{
		ExecutionId: executionID,
		Status:      result.Status,
	}, nil
}

func (l *ExecuteWorkflowLogic) markExecutionFailed(execID uint64, errMsg string) error {
	exec, err := l.svcCtx.ExecutionModel.FindOne(l.ctx, execID)
	if err != nil {
		return err
	}
	exec.Status = "failed"
	exec.ErrorMsg = sql.NullString{String: errMsg, Valid: true}
	exec.EndedAt = sql.NullTime{Time: time.Now(), Valid: true}
	return l.svcCtx.ExecutionModel.Update(l.ctx, exec)
}
