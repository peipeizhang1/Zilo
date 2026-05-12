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

type ListWorkflowsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListWorkflowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWorkflowsLogic {
	return &ListWorkflowsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListWorkflowsLogic) ListWorkflows(req *types.ListWorkflowsReq) (resp *types.ListWorkflowsResp, err error) {
	userID, err := getUserIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	out, err := l.svcCtx.WorkflowRpc.ListWorkflows(l.ctx, &workflowrpc.ListWorkflowsReq{
		UserId:   userID,
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
	})
	if err != nil {
		return nil, err
	}
	items := make([]types.WorkflowResp, 0, len(out.List))
	for _, item := range out.List {
		items = append(items, *toWorkflowResp(item))
	}
	return &types.ListWorkflowsResp{
		Total: out.Total,
		List:  items,
	}, nil
}
