package main

import (
	"example.com/crud-example/internal/api"
	"example.com/crud-example/internal/dao"
)

func main() {
	dao := dao.DAO{}
	dao.Init()
	defer dao.Deinit()

	server := api.NewAPIServer(&dao)
	server.Run()
}
