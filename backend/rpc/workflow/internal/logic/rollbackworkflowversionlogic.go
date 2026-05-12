package logic

import (
	"context"
	"errors"

	"zilo/backend/model"
	"zilo/backend/rpc/workflow/internal/svc"
	"zilo/backend/rpc/workflow/workflowpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RollbackWorkflowVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRollbackWorkflowVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RollbackWorkflowVersionLogic {
	return &RollbackWorkflowVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RollbackWorkflowVersionLogic) RollbackWorkflowVersion(in *workflowpb.RollbackReq) (*workflowpb.Empty, error) {
	if in.Version <= 0 {
		return nil, errors.New("invalid version")
	}

	workflow, err := l.svcCtx.WorkflowModel.FindOne(l.ctx, uint64(in.WorkflowId))
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("workflow not found")
		}
		return nil, err
	}
	if workflow.OwnerId != uint64(in.UserId) {
		return nil, errors.New("no permission to rollback workflow")
	}

	target, err := l.svcCtx.WorkflowVersionModel.FindOneByWorkflowIdVersion(l.ctx, uint64(in.WorkflowId), in.Version)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("target version not found")
		}
		return nil, err
	}

	target.IsPublished = 1
	if err = l.svcCtx.WorkflowVersionModel.Update(l.ctx, target); err != nil {
		return nil, err
	}

	workflow.LatestVersion = target.Version
	if err = l.svcCtx.WorkflowModel.Update(l.ctx, workflow); err != nil {
		return nil, err
	}

	return &workflowpb.Empty{}, nil
}
