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

type GetExecutionLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetExecutionLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetExecutionLogsLogic {
	return &GetExecutionLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetExecutionLogsLogic) GetExecutionLogs(in *workflowpb.ExecutionReq) (*workflowpb.ExecutionLogsResp, error) {
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
		return nil, errors.New("no permission to access logs")
	}
	logs, err := l.svcCtx.ExecutionLogModel.FindByExecutionID(l.ctx, uint64(in.ExecutionId))
	if err != nil {
		return nil, err
	}
	items := make([]*workflowpb.ExecutionLogItem, 0, len(logs))
	for _, item := range logs {
		startedAt := ""
		if item.StartedAt.Valid {
			startedAt = item.StartedAt.Time.Format(time.RFC3339)
		}
		endedAt := ""
		if item.EndedAt.Valid {
			endedAt = item.EndedAt.Time.Format(time.RFC3339)
		}
		items = append(items, &workflowpb.ExecutionLogItem{
			NodeId:     item.NodeId,
			NodeType:   item.NodeType,
			Status:     item.Status,
			InputJson:  item.InputJson.String,
			OutputJson: item.OutputJson.String,
			ErrorMsg:   item.ErrorMsg.String,
			StartedAt:  startedAt,
			EndedAt:    endedAt,
			DurationMs: item.DurationMs,
		})
	}

	return &workflowpb.ExecutionLogsResp{List: items}, nil
}
