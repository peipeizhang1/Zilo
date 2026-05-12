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

type ExecuteWorkflowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExecuteWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExecuteWorkflowLogic {
	return &ExecuteWorkflowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExecuteWorkflowLogic) ExecuteWorkflow(workflowID uint64, req *types.ExecuteReq) (resp *types.ExecuteResp, err error) {
	userID, err := getUserIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	out, err := l.svcCtx.WorkflowRpc.ExecuteWorkflow(l.ctx, &workflowrpc.ExecuteReq{
		UserId:     userID,
		WorkflowId: int64(workflowID),
		InputJson:  req.InputJson,
	})
	if err != nil {
		return nil, err
	}
	return &types.ExecuteResp{
		ExecutionId: out.ExecutionId,
		Status:      out.Status,
	}, nil
}
