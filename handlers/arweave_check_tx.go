package handlers

import (
	"aif/configs"
	"aif/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CheckTxArweave - get the transaction info or status
func CheckTxArweave(c *gin.Context) {

	txID := c.Query("transaction_id")
	if txID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please provide the transaction_id in the query ex: ?transaction_id=4dXDJacaPP_jbkOJorAmcd0eA7-oRkHjlFGilZe72bE"})
		return
	}

	err := utils.CheckValidityArweaveTxID(txID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	configuration, ok := c.MustGet("configuration").(*configs.ViperConfiguration)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the configuration"})
		return
	}

	output, status, err := utils.GetRequest(configuration.Get("nodeURL") + "/tx/" + txID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch status {
	case 200:
		c.Data(http.StatusOK, "application/json; charset=utf-8", output)
		return
	case 202:
		c.JSON(http.StatusAccepted, gin.H{"status": "pending"})
		return
	case 400:
		c.JSON(http.StatusBadRequest, gin.H{"status": "transaction ID is not valid"})
		return
	case 404:
		c.JSON(http.StatusNotFound, gin.H{"status": "transaction ID could not be found"})
		return
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"status": "something went terribly wrong"})
		return
	}

}
