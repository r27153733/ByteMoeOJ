package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/r27153733/fastgozero/core/stores/sqlx"
	"strings"
)

var _ JudgeCaseModel = (*customJudgeCaseModel)(nil)

type (
	// JudgeCaseModel is an interface to be customized, add more methods here,
	// and implement the added methods in customJudgeCaseModel.
	JudgeCaseModel interface {
		judgeCaseModel
		BatchInsert(ctx context.Context, dataList []JudgeCase) (sql.Result, error)
		withSession(session sqlx.Session) JudgeCaseModel
	}

	customJudgeCaseModel struct {
		*defaultJudgeCaseModel
	}
)

// NewJudgeCaseModel returns a model for the database table.
func NewJudgeCaseModel(conn sqlx.SqlConn) JudgeCaseModel {
	return &customJudgeCaseModel{
		defaultJudgeCaseModel: newJudgeCaseModel(conn),
	}
}

func (m *customJudgeCaseModel) withSession(session sqlx.Session) JudgeCaseModel {
	return NewJudgeCaseModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultJudgeCaseModel) BatchInsert(ctx context.Context, dataList []JudgeCase) (sql.Result, error) {
	// 构建批量插入的查询语句
	l := len(judgeCaseRowsExpectAutoSetAndIDArray)
	values := make([]interface{}, 0, len(dataList))
	placeholders := make([]string, 0, len(dataList))
	for i := 0; i < len(dataList); i++ {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
			i*l+1, i*l+2, i*l+3, i*l+4, i*l+5, i*l+6))
		values = append(values, dataList[i].JudgeId, dataList[i].ProblemDataId, dataList[i].Status, dataList[i].TimeUsed, dataList[i].MemoryUsed, dataList[i].Reason)
	}

	// 拼接完整的查询语句
	query := fmt.Sprintf("insert into %s (%s) values %s", m.table, judgeCaseRowsExpectAutoSetAndID, strings.Join(placeholders, ","))

	// 执行批量插入
	ret, err := m.conn.ExecCtx(ctx, query, values...)
	return ret, err
}
