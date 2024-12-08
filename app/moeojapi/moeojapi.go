package main

import (
	"flag"
	"fmt"

	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/config"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/handler"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/svc"

	"github.com/r27153733/fastgozero/core/conf"
	"github.com/r27153733/fastgozero/rest"
)

var configFile = flag.String("f", "etc/moeojapi.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
