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

	// get a file from IPFS by it's hash
	api.GET("/ipfs", handlers.GetFromIPFS)

	// get a file from AR by it's transaction ID
	api.GET("/arweave", handlers.GetFromArweave)

	// transfer a file from IFPS to AR
	api.POST("/transfer", handlers.TransferIPFSToArweave)

	// transfer a file directly to AR
	api.POST("/transfer_arweave", handlers.TransferToArweave)

	// checks if a tx has been mined
	api.GET("/check_tx_arweave", handlers.CheckTxArweave)

	// check the balance of AR tokens
	api.GET("/balance", handlers.GetBalance)

	// In case no route is found
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "api endpoint not found"})
	})

}
