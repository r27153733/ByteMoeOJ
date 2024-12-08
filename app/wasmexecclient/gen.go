package main

//go:generate fastgoctl rpc protoc ../wasmexec/src/pb/wasmexecutor.proto --proto_path=../wasmexec/src/pb/ --go_out=./pb --go-grpc_out=./pb --zrpc_out=. --client=true --home ../.template
//go:generate rm -rf etc && rm -rf internal && rm wasmexecutor.go
