package model

import (
	"context"
	"fmt"
	"github.com/r27153733/ByteMoeOJ/lib/uuid"
	"github.com/r27153733/fastgozero/core/stores/sqlx"
	"strings"
)

func init() {
	problemDataFieldNamesDeleteOutput := make([]string, 0, len(problemDataRows)-1)
	for i := 0; i < len(problemDataFieldNames); i++ {
		if problemDataFieldNames[i] != "output" {
			problemDataFieldNamesDeleteOutput = append(problemDataFieldNamesDeleteOutput, problemDataFieldNames[i])
		}
	}
	problemDataRowsDeleteOutput = strings.Join(problemDataFieldNames, ",")
}

var problemDataRowsDeleteOutput string

var _ ProblemDataModel = (*customProblemDataModel)(nil)

type (
	// ProblemDataModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProblemDataModel.
	ProblemDataModel interface {
		problemDataModel
		BatchInsert(ctx context.Context, data []ProblemData) error
		FindJudgeDataHash(ctx context.Context, problemID uuid.UUID) ([]ProblemData, error)
		withSession(session sqlx.Session) ProblemDataModel
	}

	customProblemDataModel struct {
		*defaultProblemDataModel
	}
)

// NewProblemDataModel returns a model for the database table.
func NewProblemDataModel(conn sqlx.SqlConn) ProblemDataModel {
	return &customProblemDataModel{
		defaultProblemDataModel: newProblemDataModel(conn),
	}
}

func (m *customProblemDataModel) withSession(session sqlx.Session) ProblemDataModel {
	return NewProblemDataModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultProblemDataModel) FindJudgeDataHash(ctx context.Context, problemID uuid.UUID) ([]ProblemData, error) {
	query := fmt.Sprintf("select %s from %s where problem_id = $1", problemDataRowsDeleteOutput, m.table)
	var resp []ProblemData
	err := m.conn.QueryRowsCtx(ctx, &resp, query, problemID)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultProblemDataModel) BatchInsert(ctx context.Context, data []ProblemData) error {
	if len(data) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(data))
	valueArgs := make([]interface{}, 0, len(data)*7)

	for i, d := range data {
		n := i * 7
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			n+1, n+2, n+3, n+4, n+5, n+6, n+7))

		valueArgs = append(valueArgs,
			d.Id,
			d.ProblemId,
			d.Input,
			d.Output,
			d.OutputHash,
			d.OutputTokenHash,
			d.OutputLen,
		)
	}

	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		m.table,
		problemDataRowsExpectAutoSet,
		strings.Join(valueStrings, ","))

	// 执行批量插入
	_, err := m.conn.ExecCtx(ctx, stmt, valueArgs...)
	return err
}
