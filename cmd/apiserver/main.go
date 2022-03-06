package main

import (
	"github.com/webdevolegkuprianov/server_rest_grpc/app/apiserver"
	logger "github.com/webdevolegkuprianov/server_rest_grpc/app/apiserver/logger"
	"github.com/webdevolegkuprianov/server_rest_grpc/app/apiserver/model"
)

func main() {

	config, err := model.NewConfig()
	if err != nil {
		logger.ErrorLogger.Println(err)
	}

	if err := apiserver.Start(config); err != nil {
		logger.ErrorLogger.Println(err)
	}

}
