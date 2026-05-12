// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package workflows

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"zilo/backend/api/internal/svc"
	"zilo/backend/rpc/workflow/workflowrpc"
)

type DeleteWorkflowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWorkflowLogic {
	return &DeleteWorkflowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteWorkflowLogic) DeleteWorkflow(workflowID uint64) error {
	userID, err := getUserIDFromCtx(l.ctx)
	if err != nil {
		return err
	}
	_, err = l.svcCtx.WorkflowRpc.DeleteWorkflow(l.ctx, &workflowrpc.WorkflowByIdReq{
		UserId:     userID,
		WorkflowId: int64(workflowID),
	})
	return err
}
