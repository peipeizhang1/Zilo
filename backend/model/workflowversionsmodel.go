package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WorkflowVersionsModel = (*customWorkflowVersionsModel)(nil)

type (
	// WorkflowVersionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWorkflowVersionsModel.
	WorkflowVersionsModel interface {
		workflowVersionsModel
		FindByWorkflowID(ctx context.Context, workflowID uint64) ([]*WorkflowVersions, error)
		FindLatestByWorkflowID(ctx context.Context, workflowID uint64) (*WorkflowVersions, error)
		withSession(session sqlx.Session) WorkflowVersionsModel
	}

	customWorkflowVersionsModel struct {
		*defaultWorkflowVersionsModel
	}
)

// NewWorkflowVersionsModel returns a model for the database table.
func NewWorkflowVersionsModel(conn sqlx.SqlConn) WorkflowVersionsModel {
	return &customWorkflowVersionsModel{
		defaultWorkflowVersionsModel: newWorkflowVersionsModel(conn),
	}
}

func (m *customWorkflowVersionsModel) withSession(session sqlx.Session) WorkflowVersionsModel {
	return NewWorkflowVersionsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customWorkflowVersionsModel) FindByWorkflowID(ctx context.Context, workflowID uint64) ([]*WorkflowVersions, error) {
	query := fmt.Sprintf("select %s from %s where `workflow_id` = ? order by `version` desc", workflowVersionsRows, m.table)
	var list []*WorkflowVersions
	if err := m.conn.QueryRowsCtx(ctx, &list, query, workflowID); err != nil {
		if err == sqlx.ErrNotFound {
			return []*WorkflowVersions{}, nil
		}
		return nil, err
	}
	return list, nil
}

func (m *customWorkflowVersionsModel) FindLatestByWorkflowID(ctx context.Context, workflowID uint64) (*WorkflowVersions, error) {
	query := fmt.Sprintf("select %s from %s where `workflow_id` = ? order by `version` desc limit 1", workflowVersionsRows, m.table)
	var item WorkflowVersions
	err := m.conn.QueryRowCtx(ctx, &item, query, workflowID)
	switch err {
	case nil:
		return &item, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
