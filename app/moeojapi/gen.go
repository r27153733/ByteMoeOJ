package main

//go:generate fastgoctl api go --api api/moeoj.api --dir . --home ../.template

//go:generate fastgoctl model pg datasource --url=postgres://postgres:@127.0.0.1:5432/byte_moe_judge?sslmode=disable --table=users --dir ./model --home ../.template
