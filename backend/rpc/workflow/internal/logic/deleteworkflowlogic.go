package logic

import (
	"context"
	"errors"

	"zilo/backend/model"
	"zilo/backend/rpc/workflow/internal/svc"
	"zilo/backend/rpc/workflow/workflowpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteWorkflowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWorkflowLogic {
	return &DeleteWorkflowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteWorkflowLogic) DeleteWorkflow(in *workflowpb.WorkflowByIdReq) (*workflowpb.Empty, error) {
	item, err := l.svcCtx.WorkflowModel.FindOne(l.ctx, uint64(in.WorkflowId))
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("workflow not found")
		}
		return nil, err
	}
	if item.OwnerId != uint64(in.UserId) {
		return nil, errors.New("no permission to delete workflow")
	}
	if err = l.svcCtx.WorkflowModel.Delete(l.ctx, uint64(in.WorkflowId)); err != nil {
		return nil, err
	}

	return &workflowpb.Empty{}, nil
}
