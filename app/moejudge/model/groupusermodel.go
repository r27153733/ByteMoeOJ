package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/r27153733/ByteMoeOJ/lib/uuid"
	"github.com/r27153733/fastgozero/core/stores/sqlx"
)

var _ GroupUserModel = (*customGroupUserModel)(nil)

type (
	// GroupUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGroupUserModel.
	GroupUserModel interface {
		groupUserModel
		UpsertByUserIdGroupId(ctx context.Context, data *GroupUser) (sql.Result, error)
		DeleteByGroupId(ctx context.Context, id uuid.UUID) error
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

func (m *defaultGroupUserModel) DeleteByGroupId(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf("delete from %s where group_id = $1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultGroupUserModel) UpsertByUserIdGroupId(ctx context.Context, data *GroupUser) (sql.Result, error) {
	query := fmt.Sprintf(`
		insert into %s (%s)
		values ($1, $2, $3, $4)
		on conflict (user_id, group_id)
		do update set %s
	`, m.table, groupUserRowsExpectAutoSet, groupUserRowsWithPlaceHolder)

	ret, err := m.conn.ExecCtx(ctx, query, data.Id, data.GroupId, data.UserId, data.Role)
	return ret, err
}
