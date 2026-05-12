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

type GetExecutionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetExecutionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetExecutionLogic {
	return &GetExecutionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetExecutionLogic) GetExecution(executionID uint64) (resp *types.ExecutionResp, err error) {
	userID, err := getUserIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	out, err := l.svcCtx.WorkflowRpc.GetExecution(l.ctx, &workflowrpc.ExecutionReq{
		UserId:      userID,
		ExecutionId: int64(executionID),
	})
	if err != nil {
		return nil, err
	}
	return &types.ExecutionResp{
		Id:         out.Id,
		WorkflowId: out.WorkflowId,
		VersionId:  out.VersionId,
		Status:     out.Status,
		InputJson:  out.InputJson,
		OutputJson: out.OutputJson,
		ErrorMsg:   out.ErrorMsg,
		StartedAt:  out.StartedAt,
		EndedAt:    out.EndedAt,
		DurationMs: out.DurationMs,
	}, nil
}
