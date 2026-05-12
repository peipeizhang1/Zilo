package workflows

import (
	"context"

	"zilo/backend/api/internal/pkg/userctx"
	"zilo/backend/api/internal/types"
	"zilo/backend/rpc/workflow/workflowrpc"
)

func getUserIDFromCtx(ctx context.Context) (int64, error) {
	return userctx.GetUserIDFromCtx(ctx)
}

func toWorkflowResp(item *workflowrpc.WorkflowResp) *types.WorkflowResp {
	return &types.WorkflowResp{
		Id:            item.Id,
		Name:          item.Name,
		Description:   item.Description,
		Status:        item.Status,
		LatestVersion: item.LatestVersion,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
	}
}
