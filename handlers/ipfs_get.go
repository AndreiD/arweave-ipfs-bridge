package handlers

import (
	"aif/configs"
	"aif/utils"
	"aif/utils/log"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
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

	out, err := exec.Command("ipfs", "get", hash).CombinedOutput()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("output from IPFS cmd %s", string(out))

	if strings.Contains(string(out), "merkledag: not found") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file not found."})
		return
	}

	content, err := ioutil.ReadFile(hash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file not found."})
		return
	}

	// cleanup
	err = cleanup(configuration, hash)
	if err != nil {
		log.Error(err)
	}

	log.Printf("done in %s", time.Since(start))
	c.String(http.StatusOK, string(content))
}
