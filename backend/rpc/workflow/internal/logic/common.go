package logic

import (
	"time"
	"zilo/backend/model"
	"zilo/backend/rpc/workflow/workflowpb"
)

func toWorkflowResp(item *model.Workflows) *workflowpb.WorkflowResp {
	return &workflowpb.WorkflowResp{
		Id:            int64(item.Id),
		Name:          item.Name,
		Description:   item.Description,
		Status:        item.Status,
		LatestVersion: item.LatestVersion,
		CreatedAt:     item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     item.UpdatedAt.Format(time.RFC3339),
	}
}
