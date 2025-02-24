package model

import (
	"context"
	"fmt"
	"github.com/r27153733/fastgozero/core/stores/sqlx"
	"strings"
)

var _ ProblemLangModel = (*customProblemLangModel)(nil)

type (
	// ProblemLangModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProblemLangModel.
	ProblemLangModel interface {
		problemLangModel
		BatchInsert(ctx context.Context, data []ProblemLang) error
		withSession(session sqlx.Session) ProblemLangModel
	}

	customProblemLangModel struct {
		*defaultProblemLangModel
	}
)

// NewProblemLangModel returns a model for the database table.
func NewProblemLangModel(conn sqlx.SqlConn) ProblemLangModel {
	return &customProblemLangModel{
		defaultProblemLangModel: newProblemLangModel(conn),
	}
}

func (m *customProblemLangModel) withSession(session sqlx.Session) ProblemLangModel {
	return NewProblemLangModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultProblemLangModel) BatchInsert(ctx context.Context, data []ProblemLang) error {
	if len(data) == 0 {
		return nil
	}

	// 构建批量插入的 values 占位符，每行 7 个字段
	valueStrings := make([]string, 0, len(data))
	valueArgs := make([]interface{}, 0, len(data)*7)

	for i, d := range data {
		// 构建如 ($1, $2, $3, $4, $5, $6, $7) 的占位符
		n := i * 7
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			n+1, n+2, n+3, n+4, n+5, n+6, n+7))

		valueArgs = append(valueArgs,
			d.Id,
			d.ProblemId,
			d.Lang,
			d.InitCode,
			d.Template,
			d.TimeLimit,
			d.MemoryLimit,
		)
	}

	// 构建完整的插入语句
	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		m.table,
		problemLangRowsExpectAutoSet,
		strings.Join(valueStrings, ","))

	// 执行批量插入
	_, err := m.conn.ExecCtx(ctx, stmt, valueArgs...)
	return err
}
