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
	problemDataFieldNames                  = builder.RawFieldNames(&ProblemData{}, true)
	problemDataRows                        = strings.Join(problemDataFieldNames, ",")
	problemDataRowsExpectAutoSet           = strings.Join(stringx.Remove(problemDataFieldNames, "create_at", "create_time", "created_at", "update_at", "update_time", "updated_at"), ",")
	problemDataRowsExpectAutoSetAndIDArray = stringx.Remove(problemDataFieldNames, "id", "create_at", "create_time", "created_at", "update_at", "update_time", "updated_at")
	problemDataRowsExpectAutoSetAndID      = strings.Join(problemDataRowsExpectAutoSetAndIDArray, ",")
	problemDataRowsWithPlaceHolder         = builder.PostgreSqlJoin(stringx.Remove(problemDataFieldNames, "id", "create_at", "create_time", "created_at", "update_at", "update_time", "updated_at"))
)

type (
	problemDataModel interface {
		Insert(ctx context.Context, data *ProblemData) (sql.Result, error)
		FindOne(ctx context.Context, id uuid.UUID) (*ProblemData, error)
		FindOneLock(ctx context.Context, id uuid.UUID) (*ProblemData, error)
		Update(ctx context.Context, data *ProblemData) error
		Upsert(ctx context.Context, data *ProblemData) (sql.Result, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	defaultProblemDataModel struct {
		conn  sqlx.SqlConn
		table string
	}

	ProblemData struct {
		Id              uuid.UUID `db:"id"`
		ProblemId       uuid.UUID `db:"problem_id"`
		Input           string    `db:"input"`
		Output          string    `db:"output"`
		OutputHash      int64     `db:"output_hash"`
		OutputTokenHash int64     `db:"output_token_hash"`
		OutputLen       int64     `db:"output_len"`
		CreatedAt       time.Time `db:"created_at"` // 创建时间
	}
)

func newProblemDataModel(conn sqlx.SqlConn) *defaultProblemDataModel {
	return &defaultProblemDataModel{
		conn:  conn,
		table: `"public"."problem_data"`,
	}
}

func (m *defaultProblemDataModel) Delete(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf("delete from %s where id = $1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultProblemDataModel) FindOne(ctx context.Context, id uuid.UUID) (*ProblemData, error) {
	query := fmt.Sprintf("select %s from %s where id = $1 limit 1", problemDataRows, m.table)
	var resp ProblemData
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

func (m *defaultProblemDataModel) FindOneLock(ctx context.Context, id uuid.UUID) (*ProblemData, error) {
	query := fmt.Sprintf("select %s from %s where id = $1 limit 1 for update", problemDataRows, m.table)
	var resp ProblemData
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

func (m *defaultProblemDataModel) Insert(ctx context.Context, data *ProblemData) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7)", m.table, problemDataRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Id, data.ProblemId, data.Input, data.Output, data.OutputHash, data.OutputTokenHash, data.OutputLen)
	return ret, err
}

func (m *defaultProblemDataModel) Update(ctx context.Context, data *ProblemData) error {
	query := fmt.Sprintf("update %s set %s where id = $1", m.table, problemDataRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Id, data.ProblemId, data.Input, data.Output, data.OutputHash, data.OutputTokenHash, data.OutputLen)
	return err
}

func (m *defaultProblemDataModel) tableName() string {
	return m.table
}

func (m *defaultProblemDataModel) Upsert(ctx context.Context, data *ProblemData) (sql.Result, error) {
	query := fmt.Sprintf(`
		insert into %s (%s)
		values ($1, $2, $3, $4, $5, $6, $7)
		on conflict (id)
		do update set %s
	`, m.table, problemDataRowsExpectAutoSet, problemDataRowsWithPlaceHolder)

	ret, err := m.conn.ExecCtx(ctx, query, data.Id, data.ProblemId, data.Input, data.Output, data.OutputHash, data.OutputTokenHash, data.OutputLen)
	return ret, err
}
