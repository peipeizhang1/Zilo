// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package workflows

import (
	"context"

	"zilo/backend/api/internal/svc"
	"zilo/backend/api/internal/types"
	"zilo/backend/rpc/workflow/workflowrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetExecutionLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetExecutionLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetExecutionLogsLogic {
	return &GetExecutionLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetExecutionLogsLogic) GetExecutionLogs(executionID uint64) (resp *types.ExecutionLogsResp, err error) {
	userID, err := getUserIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	out, err := l.svcCtx.WorkflowRpc.GetExecutionLogs(l.ctx, &workflowrpc.ExecutionReq{
		UserId:      userID,
		ExecutionId: int64(executionID),
	})
	if err != nil {
		return nil, err
	}
	items := make([]types.ExecutionLogItem, 0, len(out.List))
	for _, item := range out.List {
		items = append(items, types.ExecutionLogItem{
			NodeId:     item.NodeId,
			NodeType:   item.NodeType,
			Status:     item.Status,
			InputJson:  item.InputJson,
			OutputJson: item.OutputJson,
			ErrorMsg:   item.ErrorMsg,
			StartedAt:  item.StartedAt,
			EndedAt:    item.EndedAt,
			DurationMs: item.DurationMs,
		})
	}
	return &types.ExecutionLogsResp{List: items}, nil
}
