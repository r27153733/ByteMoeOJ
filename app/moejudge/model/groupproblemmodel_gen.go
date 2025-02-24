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
	groupProblemFieldNames                  = builder.RawFieldNames(&GroupProblem{}, true)
	groupProblemRows                        = strings.Join(groupProblemFieldNames, ",")
	groupProblemRowsExpectAutoSet           = strings.Join(stringx.Remove(groupProblemFieldNames, "create_at", "create_time", "created_at", "update_at", "update_time", "updated_at"), ",")
	groupProblemRowsExpectAutoSetAndIDArray = stringx.Remove(groupProblemFieldNames, "id", "create_at", "create_time", "created_at", "update_at", "update_time", "updated_at")
	groupProblemRowsExpectAutoSetAndID      = strings.Join(groupProblemRowsExpectAutoSetAndIDArray, ",")
	groupProblemRowsWithPlaceHolder         = builder.PostgreSqlJoin(stringx.Remove(groupProblemFieldNames, "id", "create_at", "create_time", "created_at", "update_at", "update_time", "updated_at"))
)

type (
	groupProblemModel interface {
		Insert(ctx context.Context, data *GroupProblem) (sql.Result, error)
		FindOne(ctx context.Context, id uuid.UUID) (*GroupProblem, error)
		FindOneLock(ctx context.Context, id uuid.UUID) (*GroupProblem, error)
		FindOneByProblemIdGroupId(ctx context.Context, problemId uuid.UUID, groupId uuid.UUID) (*GroupProblem, error)
		DeleteByProblemIdGroupId(ctx context.Context, problemId uuid.UUID, groupId uuid.UUID) error
		Update(ctx context.Context, data *GroupProblem) error
		Upsert(ctx context.Context, data *GroupProblem) (sql.Result, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	defaultGroupProblemModel struct {
		conn  sqlx.SqlConn
		table string
	}

	GroupProblem struct {
		Id        uuid.UUID `db:"id"`
		GroupId   uuid.UUID `db:"group_id"`
		ProblemId uuid.UUID `db:"problem_id"`
		CreatedAt time.Time `db:"created_at"` // 创建时间
	}
)

func newGroupProblemModel(conn sqlx.SqlConn) *defaultGroupProblemModel {
	return &defaultGroupProblemModel{
		conn:  conn,
		table: `"public"."group_problem"`,
	}
}

func (m *defaultGroupProblemModel) Delete(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf("delete from %s where id = $1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultGroupProblemModel) FindOne(ctx context.Context, id uuid.UUID) (*GroupProblem, error) {
	query := fmt.Sprintf("select %s from %s where id = $1 limit 1", groupProblemRows, m.table)
	var resp GroupProblem
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

func (m *defaultGroupProblemModel) FindOneLock(ctx context.Context, id uuid.UUID) (*GroupProblem, error) {
	query := fmt.Sprintf("select %s from %s where id = $1 limit 1 for update", groupProblemRows, m.table)
	var resp GroupProblem
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

func (m *defaultGroupProblemModel) FindOneByProblemIdGroupId(ctx context.Context, problemId uuid.UUID, groupId uuid.UUID) (*GroupProblem, error) {
	var resp GroupProblem
	query := fmt.Sprintf("select %s from %s where problem_id = $1 and group_id = $2 limit 1", groupProblemRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, problemId, groupId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultGroupProblemModel) DeleteByProblemIdGroupId(ctx context.Context, problemId uuid.UUID, groupId uuid.UUID) error {
	query := fmt.Sprintf("delete from %s where problem_id = $1 and group_id = $2", m.table)
	_, err := m.conn.ExecCtx(ctx, query, problemId, groupId)
	return err
}

func (m *defaultGroupProblemModel) Insert(ctx context.Context, data *GroupProblem) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3)", m.table, groupProblemRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Id, data.GroupId, data.ProblemId)
	return ret, err
}

func (m *defaultGroupProblemModel) Update(ctx context.Context, newData *GroupProblem) error {
	query := fmt.Sprintf("update %s set %s where id = $1", m.table, groupProblemRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.Id, newData.GroupId, newData.ProblemId)
	return err
}

func (m *defaultGroupProblemModel) tableName() string {
	return m.table
}

func (m *defaultGroupProblemModel) Upsert(ctx context.Context, data *GroupProblem) (sql.Result, error) {
	query := fmt.Sprintf(`
		insert into %s (%s)
		values ($1, $2, $3)
		on conflict (id)
		do update set %s
	`, m.table, groupProblemRowsExpectAutoSet, groupProblemRowsWithPlaceHolder)

	ret, err := m.conn.ExecCtx(ctx, query, data.Id, data.GroupId, data.ProblemId)
	return ret, err
}
