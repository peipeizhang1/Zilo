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

type ListWorkflowVersionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListWorkflowVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWorkflowVersionsLogic {
	return &ListWorkflowVersionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListWorkflowVersionsLogic) ListWorkflowVersions(workflowID uint64) (resp *types.VersionsResp, err error) {
	userID, err := getUserIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	out, err := l.svcCtx.WorkflowRpc.ListWorkflowVersions(l.ctx, &workflowrpc.ListVersionsReq{
		UserId:     userID,
		WorkflowId: int64(workflowID),
	})
	if err != nil {
		return nil, err
	}
	items := make([]types.VersionItem, 0, len(out.List))
	for _, it := range out.List {
		items = append(items, types.VersionItem{
			Id:          it.Id,
			Version:     it.Version,
			IsPublished: it.IsPublished,
			CreatedAt:   it.CreatedAt,
		})
	}
	return &types.VersionsResp{List: items}, nil
}
