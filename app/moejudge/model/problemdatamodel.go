package model

import (
	"context"
	"fmt"
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
		FindJudgeDataHash(ctx context.Context, problemID string) ([]ProblemData, error)
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

func (m *defaultProblemDataModel) FindJudgeDataHash(ctx context.Context, problemID string) ([]ProblemData, error) {
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
