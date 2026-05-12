package logic

import (
	"context"
	"errors"
	"strings"

	"zilo/backend/model"
	"zilo/backend/rpc/workflow/internal/svc"
	"zilo/backend/rpc/workflow/workflowpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateWorkflowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWorkflowLogic {
	return &UpdateWorkflowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateWorkflowLogic) UpdateWorkflow(in *workflowpb.UpdateWorkflowReq) (*workflowpb.WorkflowResp, error) {
	item, err := l.svcCtx.WorkflowModel.FindOne(l.ctx, uint64(in.WorkflowId))
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("workflow not found")
		}
		return nil, err
	}
	if item.OwnerId != uint64(in.UserId) {
		return nil, errors.New("no permission to update workflow")
	}

	if strings.TrimSpace(in.Name) != "" {
		item.Name = strings.TrimSpace(in.Name)
	}
	if in.Description != "" {
		item.Description = strings.TrimSpace(in.Description)
	}
	if err = l.svcCtx.WorkflowModel.Update(l.ctx, item); err != nil {
		return nil, err
	}
	item, err = l.svcCtx.WorkflowModel.FindOne(l.ctx, uint64(in.WorkflowId))
	if err != nil {
		return nil, err
	}

	return toWorkflowResp(item), nil
}
