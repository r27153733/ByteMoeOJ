package model

import (
	"context"
	"fmt"
	"github.com/r27153733/fastgozero/core/stores/sqlx"
)

var _ GroupUserModel = (*customGroupUserModel)(nil)

type (
	// GroupUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGroupUserModel.
	GroupUserModel interface {
		groupUserModel
		DeleteByGroupId(ctx context.Context, id string) error
		withSession(session sqlx.Session) GroupUserModel
	}

	customGroupUserModel struct {
		*defaultGroupUserModel
	}
)

// NewGroupUserModel returns a model for the database table.
func NewGroupUserModel(conn sqlx.SqlConn) GroupUserModel {
	return &customGroupUserModel{
		defaultGroupUserModel: newGroupUserModel(conn),
	}
}

func (m *customGroupUserModel) withSession(session sqlx.Session) GroupUserModel {
	return NewGroupUserModel(sqlx.NewSqlConnFromSession(session))
}

const (
	GroupUserRoleOwner = 100 - iota
	GroupUserRoleAdmin

	GroupUserRoleDefault = 0
)

func (m *defaultGroupUserModel) DeleteByGroupId(ctx context.Context, id string) error {
	query := fmt.Sprintf("delete from %s where group_id = $1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}
