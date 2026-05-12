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

type CreateWorkflowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWorkflowLogic {
	return &CreateWorkflowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateWorkflowLogic) CreateWorkflow(in *workflowpb.CreateWorkflowReq) (*workflowpb.WorkflowResp, error) {
	if in.UserId <= 0 {
		return nil, errors.New("invalid user id")
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, errors.New("name is required")
	}

	result, err := l.svcCtx.WorkflowModel.Insert(l.ctx, &model.Workflows{
		OwnerId:       uint64(in.UserId),
		Name:          name,
		Description:   strings.TrimSpace(in.Description),
		Status:        1,
		LatestVersion: 0,
	})
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	item, err := l.svcCtx.WorkflowModel.FindOne(l.ctx, uint64(id))
	if err != nil {
		return nil, err
	}

	return toWorkflowResp(item), nil
}
