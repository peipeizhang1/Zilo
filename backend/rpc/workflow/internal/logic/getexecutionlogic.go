package logic

import (
	"context"
	"errors"
	"time"

	"zilo/backend/model"
	"zilo/backend/rpc/workflow/internal/svc"
	"zilo/backend/rpc/workflow/workflowpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetExecutionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetExecutionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetExecutionLogic {
	return &GetExecutionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetExecutionLogic) GetExecution(in *workflowpb.ExecutionReq) (*workflowpb.ExecutionResp, error) {
	exec, err := l.svcCtx.ExecutionModel.FindOne(l.ctx, uint64(in.ExecutionId))
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("execution not found")
		}
		return nil, err
	}
	workflow, err := l.svcCtx.WorkflowModel.FindOne(l.ctx, exec.WorkflowId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("workflow not found")
		}
		return nil, err
	}
	if workflow.OwnerId != uint64(in.UserId) {
		return nil, errors.New("no permission to access execution")
	}
	startedAt := ""
	if exec.StartedAt.Valid {
		startedAt = exec.StartedAt.Time.Format(time.RFC3339)
	}
	endedAt := ""
	if exec.EndedAt.Valid {
		endedAt = exec.EndedAt.Time.Format(time.RFC3339)
	}

	return &workflowpb.ExecutionResp{
		Id:         int64(exec.Id),
		WorkflowId: int64(exec.WorkflowId),
		VersionId:  int64(exec.VersionId),
		Status:     exec.Status,
		InputJson:  exec.InputJson.String,
		OutputJson: exec.OutputJson.String,
		ErrorMsg:   exec.ErrorMsg.String,
		StartedAt:  startedAt,
		EndedAt:    endedAt,
		DurationMs: exec.DurationMs,
	}, nil
}
