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

type PublishWorkflowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishWorkflowLogic {
	return &PublishWorkflowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishWorkflowLogic) PublishWorkflow(workflowID uint64, req *types.PublishReq) (resp *types.PublishResp, err error) {
	userID, err := getUserIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	out, err := l.svcCtx.WorkflowRpc.PublishWorkflow(l.ctx, &workflowrpc.PublishReq{
		UserId:     userID,
		WorkflowId: int64(workflowID),
		Note:       req.Note,
	})
	if err != nil {
		return nil, err
	}
	return &types.PublishResp{Version: out.Version}, nil
}
