package main

import (
	"flag"
	"fmt"

	"github.com/r27153733/ByteMoeOJ/app/moejudge/internal/config"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/internal/server"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/pb"

	"github.com/r27153733/fastgozero/core/conf"
	"github.com/r27153733/fastgozero/core/service"
	"github.com/r27153733/fastgozero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/moejudge.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	c.Middlewares.StatConf = zrpc.StatConf{
		IgnoreContentMethods: []string{"666"},
	}

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterMoeJudgeSvcServer(grpcServer, server.NewMoeJudgeSvcServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
