package main

import (
	"github.com/narinjtp/bookstore_users-api/app"
	"github.com/narinjtp/bookstore_users-api/logger"
)

func main(){
	app.StartApplication()
}

func init(){
	logger.Info("main init")
}