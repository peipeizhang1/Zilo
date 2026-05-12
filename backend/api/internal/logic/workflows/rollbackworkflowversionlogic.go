package workflows

import (
	"context"

	"zilo/backend/api/internal/svc"
	"zilo/backend/rpc/workflow/workflowrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RollbackWorkflowVersionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRollbackWorkflowVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RollbackWorkflowVersionLogic {
	return &RollbackWorkflowVersionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RollbackWorkflowVersionLogic) RollbackWorkflowVersion(workflowID uint64, version int64) error {
	userID, err := getUserIDFromCtx(l.ctx)
	if err != nil {
		return err
	}

	_, err = l.svcCtx.WorkflowRpc.RollbackWorkflowVersion(l.ctx, &workflowrpc.RollbackReq{
		UserId:     userID,
		WorkflowId: int64(workflowID),
		Version:    version,
	})
	return err
}
