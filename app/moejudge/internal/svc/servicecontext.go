package svc

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/helper/mq"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/internal/config"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/model"
	"github.com/r27153733/ByteMoeOJ/app/wasmexecclient/pb/wasm"
	"github.com/r27153733/ByteMoeOJ/app/wasmexecclient/wasmexecutor"
	"github.com/r27153733/fastgozero/core/stores/sqlx"
	"github.com/r27153733/fastgozero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	WasmExec wasmexecutor.WasmExecutor
	Producer mq.Producer
	DB       model.DBCtx
}

func NewServiceContext(c config.Config) *ServiceContext {
	zrpc.DontLogClientContentForMethod(wasm.WasmExecutor_Execute_FullMethodName)

	p := mq.NewSyncProducer(&c.Kafka)

	db, err := sql.Open("pgx", c.DataSourceName)
	if err != nil {
		panic(err)
	}
	conn := sqlx.NewSqlConnFromDB(db)

	return &ServiceContext{
		Config:   c,
		WasmExec: wasmexecutor.NewWasmExecutor(zrpc.MustNewClient(c.WasmExecRPC)),
		Producer: mq.NewKafkaProducer("MoeJudge", p),
		DB:       model.NewCtx(conn),
	}
}
