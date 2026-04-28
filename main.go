package main

import (
	"example.com/crud-example/api"
	"example.com/crud-example/dao"
)

func main() {
	dao := dao.DAO{}
	dao.Init()
	defer dao.Deinit()

	server := api.NewAPIServer(&dao)
	server.Run()
}
