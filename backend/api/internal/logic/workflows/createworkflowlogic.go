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

type CreateWorkflowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWorkflowLogic {
	return &CreateWorkflowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateWorkflowLogic) CreateWorkflow(req *types.CreateWorkflowReq) (resp *types.WorkflowResp, err error) {
	userID, err := getUserIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	out, err := l.svcCtx.WorkflowRpc.CreateWorkflow(l.ctx, &workflowrpc.CreateWorkflowReq{
		UserId:      userID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return toWorkflowResp(out), nil
}
