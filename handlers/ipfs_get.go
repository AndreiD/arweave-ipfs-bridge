package handlers

import (
	"aif/configs"
	"aif/utils"
	"aif/utils/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// GetFromIPFS .
func GetFromIPFS(c *gin.Context) {

	start := time.Now()
	hash := c.Query("hash")
	if hash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please provide the hash in the query ex: ?hash=4dXDJacaPP_jbkOJorAmcd0eA7-oRkHjlFGilZe72bE"})
		return
	}

	err := utils.CheckValidityIPFSHash(hash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	configuration, ok := c.MustGet("configuration").(*configs.ViperConfiguration)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the configuration"})
		return
	}

	out, statusCode, err := utils.GetRequest(configuration.Get("ipfsGateway") + hash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("IPFS gateway returned status code %d", statusCode)

	log.Printf("done in %s", time.Since(start))
	c.String(http.StatusOK, string(out))
}
