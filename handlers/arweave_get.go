package handlers

import (
	"aif/configs"
	"aif/utils"
	"aif/utils/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// GetFromArweave .
func GetFromArweave(c *gin.Context) {

	start := time.Now()
	txID := c.Query("transaction_id")
	if txID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please provide the transaction_id in the query ex: ?transaction_id=4dXDJacaPP_jbkOJorAmcd0eA7-oRkHjlFGilZe72bE"})
		return
	}

	decode := c.Query("decode")
	if decode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please provide the decode in the query ex: ?decode=true"})
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

	output, _, err := utils.GetRequest(configuration.Get("nodeURL") + "/tx/" + txID + "/data")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if decode == "true" {
		decoded, err := utils.DecodeString(string(output))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.String(http.StatusOK, string(decoded))
		return
	}

	log.Printf("done in %s", time.Since(start))
	c.String(http.StatusOK, string(output))
}
