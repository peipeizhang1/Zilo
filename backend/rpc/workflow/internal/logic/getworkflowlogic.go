package logic

import (
	"context"
	"errors"

	"zilo/backend/model"
	"zilo/backend/rpc/workflow/internal/svc"
	"zilo/backend/rpc/workflow/workflowpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWorkflowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkflowLogic {
	return &GetWorkflowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetWorkflowLogic) GetWorkflow(in *workflowpb.WorkflowByIdReq) (*workflowpb.WorkflowResp, error) {
	item, err := l.svcCtx.WorkflowModel.FindOne(l.ctx, uint64(in.WorkflowId))
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("workflow not found")
		}
		return nil, err
	}
	if item.OwnerId != uint64(in.UserId) {
		return nil, errors.New("no permission to access workflow")
	}

	return toWorkflowResp(item), nil
}
