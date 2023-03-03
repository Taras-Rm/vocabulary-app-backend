package server

import "github.com/gin-gonic/gin"

func NewServer() *gin.Engine {
	engine := gin.Default()
	return engine
}
