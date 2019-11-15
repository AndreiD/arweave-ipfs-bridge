package handlers

import (
	"aif/arweave"
	"aif/configs"
	"github.com/gin-gonic/gin"
	"math"
	"math/big"
	"net/http"
)

// GetBalance of tokens
func GetBalance(c *gin.Context) {

	address := c.Query("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please provide the address in the query ex: ?address=qGwglm54w6I9-CCcNSAjvWzqGNZfb0zAUNkXYVYN5LY"})
		return
	}

	configuration, ok := c.MustGet("configuration").(*configs.ViperConfiguration)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the configuration"})
		return
	}

	output, err := arweave.GetBalance(address, configuration)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fBalance := new(big.Float)
	fBalance.SetString(string(output))
	arBalance := new(big.Float).Quo(fBalance, big.NewFloat(math.Pow10(12)))

	c.JSON(http.StatusOK, gin.H{"winston": string(output), "ar": arBalance.String()})
}
