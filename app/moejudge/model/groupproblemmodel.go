package model

import (
	"context"
	"fmt"
	"github.com/r27153733/fastgozero/core/stores/sqlx"
)

var _ GroupProblemModel = (*customGroupProblemModel)(nil)

type (
	// GroupProblemModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGroupProblemModel.
	GroupProblemModel interface {
		groupProblemModel
		DeleteByGroupId(ctx context.Context, id string) error
		withSession(session sqlx.Session) GroupProblemModel
	}

	customGroupProblemModel struct {
		*defaultGroupProblemModel
	}
)

// NewGroupProblemModel returns a model for the database table.
func NewGroupProblemModel(conn sqlx.SqlConn) GroupProblemModel {
	return &customGroupProblemModel{
		defaultGroupProblemModel: newGroupProblemModel(conn),
	}
}

func (m *customGroupProblemModel) withSession(session sqlx.Session) GroupProblemModel {
	return NewGroupProblemModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultGroupProblemModel) DeleteByGroupId(ctx context.Context, id string) error {
	query := fmt.Sprintf("delete from %s where group_id = $1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}
