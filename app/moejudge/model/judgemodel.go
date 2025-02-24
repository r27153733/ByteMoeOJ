package model

import (
	"context"
	"fmt"
	"github.com/r27153733/ByteMoeOJ/lib/uuid"
	"github.com/r27153733/fastgozero/core/stores/sqlx"
	"strings"
)

var _ JudgeModel = (*customJudgeModel)(nil)

type (
	// JudgeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customJudgeModel.
	JudgeModel interface {
		judgeModel
		ListJudge(ctx context.Context, req *ListJudgeReq) ([]Judge, uint32, error)
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

type ListJudgeReq struct {
	JudgeId   *uuid.UUID `protobuf:"bytes,1,opt,name=judge_id,json=judgeId,proto3" json:"judge_id,omitempty"`       // ID（可选）
	UserId    *uuid.UUID `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`          // 用户 ID（可选）
	ProblemId *uuid.UUID `protobuf:"bytes,3,opt,name=problem_id,json=problemId,proto3" json:"problem_id,omitempty"` // 题目 ID（可选）
	GroupId   *uuid.UUID `protobuf:"bytes,4,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`       // 组 ID（可选）
	Lang      int32      `protobuf:"varint,5,opt,name=lang,proto3,enum=LangType" json:"lang,omitempty"`             // 语言（可选）
	Status    int32      `protobuf:"varint,6,opt,name=status,proto3" json:"status,omitempty"`                       // 判题状态（可选）
	Page      uint32     `protobuf:"varint,7,opt,name=page,proto3" json:"page,omitempty"`                           // 分页页码（可选）
	PageSize  uint32     `protobuf:"varint,8,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`   // 分页大小（可选）
}

func (m *defaultJudgeModel) ListJudge(ctx context.Context, req *ListJudgeReq) ([]Judge, uint32, error) {
	var conditions []string
	var args []interface{}
	paramCount := 1

	// 构建查询条件
	if req.JudgeId != nil {
		conditions = append(conditions, fmt.Sprintf("id = $%d", paramCount))
		args = append(args, req.JudgeId)
		paramCount++
	}
	if req.UserId != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", paramCount))
		args = append(args, req.UserId)
		paramCount++
	}
	if req.ProblemId != nil {
		conditions = append(conditions, fmt.Sprintf("problem_id = $%d", paramCount))
		args = append(args, req.ProblemId)
		paramCount++
	}
	if req.GroupId != nil {
		conditions = append(conditions, fmt.Sprintf("group_id = $%d", paramCount))
		args = append(args, req.GroupId)
		paramCount++
	}
	if req.Lang != -1 {
		conditions = append(conditions, fmt.Sprintf("lang = $%d", paramCount))
		args = append(args, req.Lang)
		paramCount++
	}
	if req.Status != -1 {
		conditions = append(conditions, fmt.Sprintf("status = $%d", paramCount))
		args = append(args, req.Status)
		paramCount++
	}

	// 构建 WHERE 子句
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 200
	}
	offset := (req.Page - 1) * req.PageSize

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM judge %s", whereClause)
	var total uint32
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 构建分页查询
	query := fmt.Sprintf(`
        SELECT %s 
        FROM judge 
        %s 
        ORDER BY created_at DESC 
        LIMIT $%d 
        OFFSET $%d`,
		judgeRows,
		whereClause,
		paramCount,
		paramCount+1,
	)

	// 添加分页参数
	args = append(args, req.PageSize, offset)

	var judges []Judge
	err = m.conn.QueryRowsCtx(ctx, &judges, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return judges, total, nil
}
