package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ WorkflowExecutionsModel = (*customWorkflowExecutionsModel)(nil)

type (
	// WorkflowExecutionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWorkflowExecutionsModel.
	WorkflowExecutionsModel interface {
		workflowExecutionsModel
		withSession(session sqlx.Session) WorkflowExecutionsModel
	}

	customWorkflowExecutionsModel struct {
		*defaultWorkflowExecutionsModel
	}
)

// NewWorkflowExecutionsModel returns a model for the database table.
func NewWorkflowExecutionsModel(conn sqlx.SqlConn) WorkflowExecutionsModel {
	return &customWorkflowExecutionsModel{
		defaultWorkflowExecutionsModel: newWorkflowExecutionsModel(conn),
	}
}

func (m *customWorkflowExecutionsModel) withSession(session sqlx.Session) WorkflowExecutionsModel {
	return NewWorkflowExecutionsModel(sqlx.NewSqlConnFromSession(session))
}
