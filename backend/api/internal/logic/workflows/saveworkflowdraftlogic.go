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

type SaveWorkflowDraftLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveWorkflowDraftLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveWorkflowDraftLogic {
	return &SaveWorkflowDraftLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveWorkflowDraftLogic) SaveWorkflowDraft(workflowID uint64, req *types.SaveDraftReq) error {
	userID, err := getUserIDFromCtx(l.ctx)
	if err != nil {
		return err
	}
	_, err = l.svcCtx.WorkflowRpc.SaveWorkflowDraft(l.ctx, &workflowrpc.SaveDraftReq{
		UserId:     userID,
		WorkflowId: int64(workflowID),
		DslJson:    req.DslJson,
	})
	return err
}
