package model

import (
	"context"
	"github.com/r27153733/fastgozero/core/stores/sqlx"
)

type DBCtx struct {
	conn         sqlx.SqlConn
	Judge        JudgeModel
	Group        GroupModel
	GroupUser    GroupUserModel
	GroupProblem GroupProblemModel
	Problem      ProblemModel
}

func NewCtx(conn sqlx.SqlConn) DBCtx {
	return DBCtx{
		conn:         conn,
		Judge:        NewJudgeModel(conn),
		Group:        NewGroupModel(conn),
		GroupUser:    NewGroupUserModel(conn),
		GroupProblem: NewGroupProblemModel(conn),
		Problem:      NewProblemModel(conn),
	}
}

func (c DBCtx) TransactCtx(ctx context.Context, fn func(ctx context.Context, db DBCtx) error) error {
	return c.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, NewCtx(sqlx.NewSqlConnFromSession(session)))
	})
}
