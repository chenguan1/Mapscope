package server

import "github.com/gin-gonic/gin"

// ping
func handlePing(c *gin.Context)  {
	res := NewRes()
	res.Done(c,"pong")
}