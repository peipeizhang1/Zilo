package logic

import (
	"context"
	"errors"
	"time"

	"zilo/backend/model"
	"zilo/backend/rpc/workflow/internal/svc"
	"zilo/backend/rpc/workflow/workflowpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListWorkflowVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListWorkflowVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWorkflowVersionsLogic {
	return &ListWorkflowVersionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListWorkflowVersionsLogic) ListWorkflowVersions(in *workflowpb.ListVersionsReq) (*workflowpb.ListVersionsResp, error) {
	workflow, err := l.svcCtx.WorkflowModel.FindOne(l.ctx, uint64(in.WorkflowId))
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("workflow not found")
		}
		return nil, err
	}
	if workflow.OwnerId != uint64(in.UserId) {
		return nil, errors.New("no permission to access versions")
	}

	list, err := l.svcCtx.WorkflowVersionModel.FindByWorkflowID(l.ctx, uint64(in.WorkflowId))
	if err != nil {
		return nil, err
	}
	items := make([]*workflowpb.VersionItem, 0, len(list))
	for _, it := range list {
		items = append(items, &workflowpb.VersionItem{
			Id:          int64(it.Id),
			Version:     it.Version,
			IsPublished: it.IsPublished == 1,
			CreatedAt:   it.CreatedAt.Format(time.RFC3339),
		})
	}

	return &workflowpb.ListVersionsResp{List: items}, nil
}
