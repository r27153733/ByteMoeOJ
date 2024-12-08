package model

import (
	"github.com/r27153733/fastgozero/core/stores/sqlx"
)

var _ JudgeModel = (*customJudgeModel)(nil)

type (
	// JudgeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customJudgeModel.
	JudgeModel interface {
		judgeModel
		withSession(session sqlx.Session) JudgeModel
	}

	customJudgeModel struct {
		*defaultJudgeModel
	}
)

// NewJudgeModel returns a model for the database table.
func NewJudgeModel(conn sqlx.SqlConn) JudgeModel {
	return &customJudgeModel{
		defaultJudgeModel: newJudgeModel(conn),
	}
}

func (m *customJudgeModel) withSession(session sqlx.Session) JudgeModel {
	return NewJudgeModel(sqlx.NewSqlConnFromSession(session))
}
