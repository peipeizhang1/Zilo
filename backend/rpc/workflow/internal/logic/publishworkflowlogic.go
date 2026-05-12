package logic

import (
	"context"
	"errors"

	"zilo/backend/model"
	"zilo/backend/rpc/workflow/internal/svc"
	"zilo/backend/rpc/workflow/workflowpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishWorkflowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishWorkflowLogic {
	return &PublishWorkflowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishWorkflowLogic) PublishWorkflow(in *workflowpb.PublishReq) (*workflowpb.PublishResp, error) {
	workflow, err := l.svcCtx.WorkflowModel.FindOne(l.ctx, uint64(in.WorkflowId))
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("workflow not found")
		}
		return nil, err
	}
	if workflow.OwnerId != uint64(in.UserId) {
		return nil, errors.New("no permission to publish workflow")
	}
	latest, err := l.svcCtx.WorkflowVersionModel.FindLatestByWorkflowID(l.ctx, uint64(in.WorkflowId))
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("no draft version to publish")
		}
		return nil, err
	}
	latest.IsPublished = 1
	if err = l.svcCtx.WorkflowVersionModel.Update(l.ctx, latest); err != nil {
		return nil, err
	}
	workflow.LatestVersion = latest.Version
	if err = l.svcCtx.WorkflowModel.Update(l.ctx, workflow); err != nil {
		return nil, err
	}

	return &workflowpb.PublishResp{Version: latest.Version}, nil
}
