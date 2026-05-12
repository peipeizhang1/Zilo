package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WorkflowsModel = (*customWorkflowsModel)(nil)

type (
	// WorkflowsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWorkflowsModel.
	WorkflowsModel interface {
		workflowsModel
		FindByOwner(ctx context.Context, ownerID uint64, keyword string, page, pageSize int64) ([]*Workflows, int64, error)
		withSession(session sqlx.Session) WorkflowsModel
	}

	customWorkflowsModel struct {
		*defaultWorkflowsModel
	}
)

// NewWorkflowsModel returns a model for the database table.
func NewWorkflowsModel(conn sqlx.SqlConn) WorkflowsModel {
	return &customWorkflowsModel{
		defaultWorkflowsModel: newWorkflowsModel(conn),
	}
}

func (m *customWorkflowsModel) withSession(session sqlx.Session) WorkflowsModel {
	return NewWorkflowsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customWorkflowsModel) FindByOwner(ctx context.Context, ownerID uint64, keyword string, page, pageSize int64) ([]*Workflows, int64, error) {
	where := " where `owner_id` = ? "
	args := []interface{}{ownerID}
	if strings.TrimSpace(keyword) != "" {
		where += " and `name` like ? "
		args = append(args, "%"+strings.TrimSpace(keyword)+"%")
	}

	var total int64
	countQuery := fmt.Sprintf("select count(1) from %s %s", m.table, where)
	if err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []*Workflows{}, 0, nil
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	listQuery := fmt.Sprintf("select %s from %s %s order by `id` desc limit ? offset ?", workflowsRows, m.table, where)
	args = append(args, pageSize, offset)

	var list []*Workflows
	if err := m.conn.QueryRowsCtx(ctx, &list, listQuery, args...); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
