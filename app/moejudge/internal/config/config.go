package config

import (
	"github.com/r27153733/ByteMoeOJ/app/moejudge/helper/mq"
	"github.com/r27153733/fastgozero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	DataSourceName string
	WasmExecRPC    zrpc.RpcClientConf
	Kafka          mq.KafkaConfig
}
