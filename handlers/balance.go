package handlers

import (
	"aif/arweave"
	"aif/configs"
	"aif/utils"
	"github.com/gin-gonic/gin"
	"math"
	"math/big"
	"net/http"
)

// GetBalance of tokens of the wallet in use
func GetBalance(c *gin.Context) {

	configuration, ok := c.MustGet("configuration").(*configs.ViperConfiguration)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the configuration"})
		return
	}

	// get my address
	arWallet := arweave.NewWallet()
	err := arWallet.LoadKeyFromFile(configuration.Get("walletFile"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	output, _, err := utils.GetRequest(configuration.Get("nodeURL") + "/wallet/" + arWallet.Address() + "/balance")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fBalance := new(big.Float)
	fBalance.SetString(string(output))
	arBalance := new(big.Float).Quo(fBalance, big.NewFloat(math.Pow10(12)))

	c.JSON(http.StatusOK, gin.H{"winston": string(output), "ar": arBalance.String()})
}
