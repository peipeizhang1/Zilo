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

type GetWorkflowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkflowLogic {
	return &GetWorkflowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWorkflowLogic) GetWorkflow(workflowID uint64) (resp *types.WorkflowResp, err error) {
	userID, err := getUserIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	out, err := l.svcCtx.WorkflowRpc.GetWorkflow(l.ctx, &workflowrpc.WorkflowByIdReq{
		UserId:     userID,
		WorkflowId: int64(workflowID),
	})
	if err != nil {
		return nil, err
	}
	return toWorkflowResp(out), nil
}
