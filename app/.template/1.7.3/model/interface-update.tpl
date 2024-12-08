Update(ctx context.Context, {{if .containsIndexCache}}newData{{else}}data{{end}} *{{.upperStartCamelObject}}) error
Upsert(ctx context.Context, data *{{.upperStartCamelObject}}) (sql.Result,error)