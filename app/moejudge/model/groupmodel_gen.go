// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.3

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/r27153733/fastgozero/core/stores/builder"
	"github.com/r27153733/fastgozero/core/stores/sqlx"
	"github.com/r27153733/fastgozero/core/stringx"

	"github.com/r27153733/ByteMoeOJ/lib/uuid"
)

var (
	groupFieldNames                  = builder.RawFieldNames(&Group{}, true)
	groupRows                        = strings.Join(groupFieldNames, ",")
	groupRowsExpectAutoSet           = strings.Join(stringx.Remove(groupFieldNames, "create_at", "create_time", "created_at", "update_at", "update_time", "updated_at"), ",")
	groupRowsExpectAutoSetAndIDArray = stringx.Remove(groupFieldNames, "id", "create_at", "create_time", "created_at", "update_at", "update_time", "updated_at")
	groupRowsExpectAutoSetAndID      = strings.Join(groupRowsExpectAutoSetAndIDArray, ",")
	groupRowsWithPlaceHolder         = builder.PostgreSqlJoin(stringx.Remove(groupFieldNames, "id", "create_at", "create_time", "created_at", "update_at", "update_time", "updated_at"))
)

type (
	groupModel interface {
		Insert(ctx context.Context, data *Group) (sql.Result, error)
		FindOne(ctx context.Context, id uuid.UUID) (*Group, error)
		FindOneLock(ctx context.Context, id uuid.UUID) (*Group, error)
		Update(ctx context.Context, data *Group) error
		Upsert(ctx context.Context, data *Group) (sql.Result, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	defaultGroupModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Group struct {
		Id        uuid.UUID `db:"id"`
		Title     string    `db:"title"`
		Content   string    `db:"content"`
		CreatedAt time.Time `db:"created_at"` // 创建时间
	}
)

func newGroupModel(conn sqlx.SqlConn) *defaultGroupModel {
	return &defaultGroupModel{
		conn:  conn,
		table: `"public"."group"`,
	}
}

func (m *defaultGroupModel) Delete(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf("delete from %s where id = $1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultGroupModel) FindOne(ctx context.Context, id uuid.UUID) (*Group, error) {
	query := fmt.Sprintf("select %s from %s where id = $1 limit 1", groupRows, m.table)
	var resp Group
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultGroupModel) FindOneLock(ctx context.Context, id uuid.UUID) (*Group, error) {
	query := fmt.Sprintf("select %s from %s where id = $1 limit 1 for update", groupRows, m.table)
	var resp Group
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultGroupModel) Insert(ctx context.Context, data *Group) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3)", m.table, groupRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Id, data.Title, data.Content)
	return ret, err
}

func (m *defaultGroupModel) Update(ctx context.Context, data *Group) error {
	query := fmt.Sprintf("update %s set %s where id = $1", m.table, groupRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Id, data.Title, data.Content)
	return err
}

func (m *defaultGroupModel) tableName() string {
	return m.table
}

func (m *defaultGroupModel) Upsert(ctx context.Context, data *Group) (sql.Result, error) {
	query := fmt.Sprintf(`
		insert into %s (%s)
		values ($1, $2, $3)
		on conflict (id)
		do update set %s
	`, m.table, groupRowsExpectAutoSet, groupRowsWithPlaceHolder)

	ret, err := m.conn.ExecCtx(ctx, query, data.Id, data.Title, data.Content)
	return ret, err
}
