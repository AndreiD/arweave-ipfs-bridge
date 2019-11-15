package main

import (
	"aif/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitializeRouter initialises all the routes in the app
func InitializeRouter() {

	api := router.Group("/api")

	// pong
	api.GET("/ping", handlers.Ping)

	// transfer a file from IFPS to Arweave
	api.GET("/transfer", handlers.TransferIPFSToArweave)

	// check the balance of arweave tokens
	api.GET("/balance", handlers.GetBalance)


	// In case no route is found
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "api endpoint not found"})
	})

}
