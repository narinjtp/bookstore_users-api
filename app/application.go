package app

import (
	"github.com/gin-gonic/gin"
	"github.com/narinjtp/bookstore_users-api/logger"
	"github.com/narinjtp/bookstore_users-api/test"
)

var(
	router = gin.Default()
)
func StartApplication(){
	mapUrls()
	test.Test()
	logger.Info("about to start the application...")
	router.Run(":8079")
}