package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Ping responds with pong
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
