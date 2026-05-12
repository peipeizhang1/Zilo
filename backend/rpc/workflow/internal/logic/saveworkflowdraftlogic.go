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

type SaveWorkflowDraftLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveWorkflowDraftLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveWorkflowDraftLogic {
	return &SaveWorkflowDraftLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SaveWorkflowDraftLogic) SaveWorkflowDraft(in *workflowpb.SaveDraftReq) (*workflowpb.Empty, error) {
	if strings.TrimSpace(in.DslJson) == "" {
		return nil, errors.New("dslJson is required")
	}
	workflow, err := l.svcCtx.WorkflowModel.FindOne(l.ctx, uint64(in.WorkflowId))
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("workflow not found")
		}
		return nil, err
	}
	if workflow.OwnerId != uint64(in.UserId) {
		return nil, errors.New("no permission to modify workflow")
	}
	latest, err := l.svcCtx.WorkflowVersionModel.FindLatestByWorkflowID(l.ctx, uint64(in.WorkflowId))
	nextVersion := int64(1)
	if err == nil {
		nextVersion = latest.Version + 1
	} else if !errors.Is(err, model.ErrNotFound) {
		return nil, err
	}
	_, err = l.svcCtx.WorkflowVersionModel.Insert(l.ctx, &model.WorkflowVersions{
		WorkflowId:  uint64(in.WorkflowId),
		Version:     nextVersion,
		DslJson:     in.DslJson,
		IsPublished: 0,
		CreatedBy:   uint64(in.UserId),
	})
	if err != nil {
		return nil, err
	}

	return &workflowpb.Empty{}, nil
}
