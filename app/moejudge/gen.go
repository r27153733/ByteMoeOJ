package main

//go:generate fastgoctl rpc protoc ./pb/moejudge.proto --go_out=./ --go-grpc_out=./ --zrpc_out=. --client=true --home ../.template

//go:generate fastgoctl model pg datasource --url=postgres://postgres:@127.0.0.1:5432/byte_moe_judge?sslmode=disable --table=group,group_problem,group_user,judge,judge_case,problem,problem_data,problem_lang --dir ./model --home ../.template

//go:generate go install github.com/r27153733/ByteMoeOJ/app/moecompile@latest
