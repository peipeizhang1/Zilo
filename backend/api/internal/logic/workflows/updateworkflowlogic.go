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

type UpdateWorkflowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWorkflowLogic {
	return &UpdateWorkflowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateWorkflowLogic) UpdateWorkflow(workflowID uint64, req *types.UpdateWorkflowReq) (resp *types.WorkflowResp, err error) {
	userID, err := getUserIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	out, err := l.svcCtx.WorkflowRpc.UpdateWorkflow(l.ctx, &workflowrpc.UpdateWorkflowReq{
		UserId:      userID,
		WorkflowId:  int64(workflowID),
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return toWorkflowResp(out), nil
}
