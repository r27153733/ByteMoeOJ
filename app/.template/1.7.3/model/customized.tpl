
func (m *default{{.upperStartCamelObject}}Model) Upsert(ctx context.Context, data *{{.upperStartCamelObject}}) (sql.Result,error) {
	{{if .withCache}}{{.keys}}
        ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
            query := fmt.Sprintf(`
                    insert into %s (%s)
                    values ({{.expression}})
                    on conflict ({{.originalPrimaryKey}})
                    do update set %s
                `, m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet, {{.lowerStartCamelObject}}RowsWithPlaceHolder)

    		return conn.ExecCtx(ctx, query, {{.expressionValues}})
    	}, {{.keyValues}}){{else}}query := fmt.Sprintf(`
		insert into %s (%s)
		values ({{.expression}})
		on conflict ({{.originalPrimaryKey}})
		do update set %s
	`, m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet, {{.lowerStartCamelObject}}RowsWithPlaceHolder)

	ret, err := m.conn.ExecCtx(ctx, query, {{.expressionValues}}){{end}}
	return ret, err
}
