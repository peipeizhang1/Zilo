package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WorkflowExecutionLogsModel = (*customWorkflowExecutionLogsModel)(nil)

type (
	// WorkflowExecutionLogsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWorkflowExecutionLogsModel.
	WorkflowExecutionLogsModel interface {
		workflowExecutionLogsModel
		FindByExecutionID(ctx context.Context, executionID uint64) ([]*WorkflowExecutionLogs, error)
		withSession(session sqlx.Session) WorkflowExecutionLogsModel
	}

	customWorkflowExecutionLogsModel struct {
		*defaultWorkflowExecutionLogsModel
	}
)

// NewWorkflowExecutionLogsModel returns a model for the database table.
func NewWorkflowExecutionLogsModel(conn sqlx.SqlConn) WorkflowExecutionLogsModel {
	return &customWorkflowExecutionLogsModel{
		defaultWorkflowExecutionLogsModel: newWorkflowExecutionLogsModel(conn),
	}
}

func (m *customWorkflowExecutionLogsModel) withSession(session sqlx.Session) WorkflowExecutionLogsModel {
	return NewWorkflowExecutionLogsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customWorkflowExecutionLogsModel) FindByExecutionID(ctx context.Context, executionID uint64) ([]*WorkflowExecutionLogs, error) {
	query := fmt.Sprintf("select %s from %s where `execution_id` = ? order by `id` asc", workflowExecutionLogsRows, m.table)
	var list []*WorkflowExecutionLogs
	if err := m.conn.QueryRowsCtx(ctx, &list, query, executionID); err != nil {
		if err == sqlx.ErrNotFound {
			return []*WorkflowExecutionLogs{}, nil
		}
		return nil, err
	}
	return list, nil
}
