package handlers

import (
	"aif/utils"
	"github.com/gin-gonic/gin"
	"math"
	"math/big"
	"net/http"
)

// GetBalance of tokens
func GetBalance(c *gin.Context) {

	wallet := c.Query("wallet")
	if wallet == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please provide the wallet in the query ex: ?wallet=qGwglm54w6I9-CCcNSAjvWzqGNZfb0zAUNkXYVYN5LY"})
		return
	}

	output, err := utils.GetRequest("http://arweave.net/wallet/" + wallet + "/balance")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fBalance := new(big.Float)
	fBalance.SetString(string(output))
	arBalance := new(big.Float).Quo(fBalance, big.NewFloat(math.Pow10(12)))

	c.JSON(http.StatusOK, gin.H{"winston": string(output), "ar": arBalance.String()})
}
