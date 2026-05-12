package logic

import (
	"context"
	"errors"

	"zilo/backend/rpc/workflow/internal/svc"
	"zilo/backend/rpc/workflow/workflowpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListWorkflowsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListWorkflowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWorkflowsLogic {
	return &ListWorkflowsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListWorkflowsLogic) ListWorkflows(in *workflowpb.ListWorkflowsReq) (*workflowpb.ListWorkflowsResp, error) {
	if in.UserId <= 0 {
		return nil, errors.New("invalid user id")
	}
	page := in.Page
	pageSize := in.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	list, total, err := l.svcCtx.WorkflowModel.FindByOwner(l.ctx, uint64(in.UserId), in.Keyword, page, pageSize)
	if err != nil {
		return nil, err
	}
	items := make([]*workflowpb.WorkflowResp, 0, len(list))
	for _, item := range list {
		items = append(items, toWorkflowResp(item))
	}

	return &workflowpb.ListWorkflowsResp{
		Total: total,
		List:  items,
	}, nil
}
