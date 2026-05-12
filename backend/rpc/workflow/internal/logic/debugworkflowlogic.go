package logic

import (
	"context"

	"zilo/backend/rpc/workflow/internal/svc"
	"zilo/backend/rpc/workflow/workflowpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DebugWorkflowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDebugWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DebugWorkflowLogic {
	return &DebugWorkflowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DebugWorkflowLogic) DebugWorkflow(in *workflowpb.ExecuteReq) (*workflowpb.ExecuteResp, error) {
	exe := NewExecuteWorkflowLogic(l.ctx, l.svcCtx)
	resp, err := exe.ExecuteWorkflow(in)
	if err != nil {
		return nil, err
	}
	resp.Status = "debug_success"
	return resp, nil
}
