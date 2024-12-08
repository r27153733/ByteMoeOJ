func (m *default{{.upperStartCamelObject}}Model) Update(ctx context.Context, {{if .containsIndexCache}}newData{{else}}data{{end}} *{{.upperStartCamelObject}}) error {
	{{if .withCache}}{{if .containsIndexCache}}data, err:=m.FindOne(ctx, newData.{{.upperStartCamelPrimaryKey}})
	if err!=nil{
		return err
	}

{{end}}	{{.keys}}
    _, {{if .containsIndexCache}}err{{else}}err:{{end}}= m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, {{.expressionValues}})
	}, {{.keyValues}}){{else}}query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
    _,err:=m.conn.ExecCtx(ctx, query, {{.expressionValues}}){{end}}
	return err
}

func (m *default{{.upperStartCamelObject}}Model) Upsert(ctx context.Context, {{if .containsIndexCache}}newData{{else}}data{{end}} *{{.upperStartCamelObject}}) (sql.Result,error) {
	{{if .withCache}}{{.keys}}
        ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
            query := fmt.Sprintf(`
                    insert into %s (%s)
                    values ($1, $2, $3, $4, $5, $6, $7, $8)
                    on conflict ({{.originalPrimaryKey}})
                    do update set %s
                `, m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet, {{.lowerStartCamelObject}}RowsWithPlaceHolder)

    		return conn.ExecCtx(ctx, query, {{.expressionValues}})
    	}, {{.keyValues}}){{else}}query := fmt.Sprintf(`
		insert into %s (%s)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
		on conflict ({{.originalPrimaryKey}})
		do update set %s
	`, m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet, {{.lowerStartCamelObject}}RowsWithPlaceHolder)

	ret, err := m.conn.ExecCtx(ctx, query, {{.expressionValues}}){{end}}
	return ret, err
}
