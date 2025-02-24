package model

import (
	"context"
	"fmt"
	"github.com/r27153733/ByteMoeOJ/lib/uuid"
	"github.com/r27153733/fastgozero/core/stores/sqlx"
	"strings"
)

var _ GroupModel = (*customGroupModel)(nil)

type (
	// GroupModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGroupModel.
	GroupModel interface {
		groupModel
		ListGroup(ctx context.Context, req *ListGroupReq) ([]Group, uint32, error)
		withSession(session sqlx.Session) GroupModel
	}

	customGroupModel struct {
		*defaultGroupModel
	}
)

// NewGroupModel returns a model for the database table.
func NewGroupModel(conn sqlx.SqlConn) GroupModel {
	return &customGroupModel{
		defaultGroupModel: newGroupModel(conn),
	}
}

func (m *customGroupModel) withSession(session sqlx.Session) GroupModel {
	return NewGroupModel(sqlx.NewSqlConnFromSession(session))
}

type ListGroupReq struct {
	GroupId *uuid.UUID `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	UserId  *uuid.UUID `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	MinRole int32      `protobuf:"varint,3,opt,name=min_role,json=minRole,proto3" json:"min_role,omitempty"`
}

func (m *defaultGroupModel) ListGroup(ctx context.Context, req *ListGroupReq) ([]Group, uint32, error) {
	var conditions []string
	var args []interface{}
	paramCount := 1 // PostgreSQL 参数计数器

	// 构建查询条件
	if req.GroupId != nil {
		conditions = append(conditions, fmt.Sprintf("g.id = $%d", paramCount))
		args = append(args, req.GroupId)
		paramCount++
	}
	if req.UserId != nil {
		conditions = append(conditions, fmt.Sprintf("gm.user_id = $%d", paramCount))
		args = append(args, req.UserId)
		paramCount++
	}
	if req.MinRole >= 0 {
		conditions = append(conditions, fmt.Sprintf("gm.role >= $%d", paramCount))
		args = append(args, req.MinRole)
		paramCount++
	}

	// 构建 WHERE 子句
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// 查询总数
	countQuery := fmt.Sprintf(`
        SELECT COUNT(g.id) 
        FROM "group" g 
        LEFT JOIN group_user gm ON g.id = gm.group_id 
        %s`, whereClause)

	var total uint32
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 构建分页查询
	query := fmt.Sprintf(`
        SELECT 
            g.id, g.title, g.content,
            g.created_at
        FROM "group" g 
        LEFT JOIN group_user gm ON g.id = gm.group_id 
        %s 
        GROUP BY g.id
        ORDER BY g.created_at DESC 
        LIMIT $%d OFFSET $%d`,
		whereClause,
		paramCount,
		paramCount+1,
	)

	// 添加分页参数
	pageSize := 20 // 默认每页大小
	offset := 0    // 默认偏移量
	args = append(args, pageSize, offset)

	var groups []Group
	err = m.conn.QueryRowsCtx(ctx, &groups, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}
