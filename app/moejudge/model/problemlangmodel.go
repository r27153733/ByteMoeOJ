package model

import "github.com/r27153733/fastgozero/core/stores/sqlx"

var _ ProblemLangModel = (*customProblemLangModel)(nil)

type (
	// ProblemLangModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProblemLangModel.
	ProblemLangModel interface {
		problemLangModel
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
