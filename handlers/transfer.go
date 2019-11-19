package handlers

import (
	"aif/arweave"
	"aif/configs"
	"aif/utils"
	"aif/utils/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

// TransferIPFSToArweave .
func TransferIPFSToArweave(c *gin.Context) {

	start := time.Now()

	type tags struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	type newTransfer struct {
		IPFSHash       string `json:"ipfs_hash" binding:"required"`
		UseCompression bool   `json:"use_compression" `
		Tags           []tags `json:"tags"`
	}

	var nTransfer newTransfer

	err := c.BindJSON(&nTransfer)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload " + err.Error()})
		return
	}

	err = utils.CheckValidityIPFSHash(nTransfer.IPFSHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	configuration, ok := c.MustGet("configuration").(*configs.ViperConfiguration)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the configuration"})
		return
	}

	log.Printf("starting retrieving from IPFS %s", nTransfer.IPFSHash)

	// get from IPFS
	out, statusCode, err := utils.GetRequest(configuration.Get("ipfsGateway") + nTransfer.IPFSHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("IPFS gateway returned status code %d", statusCode)

	// saving it to the filesystem
	f, err := os.Create(nTransfer.IPFSHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = f.Write(out)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.Close(f)

	log.Println("file retrieved successfully from IPFS")

	var arTags []arweave.Tag
	for _, tag := range nTransfer.Tags {
		arTags = append(arTags, arweave.Tag{
			Name:  tag.Key,
			Value: tag.Value,
		})
	}

	// uploading it to Arweave
	txID, payloadLen, err := arweave.Transfer(nTransfer.IPFSHash, nTransfer.UseCompression, arTags, configuration)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Transfer to Arweave finished successfully. Tx ID %s", txID)

	err = cleanup(configuration, nTransfer.IPFSHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "transfer completed but I couldn't cleanup the file " + err.Error()})
		return
	}

	err = cleanup(configuration, nTransfer.IPFSHash+".zip")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "transfer completed but I couldn't cleanup the file " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": txID, "payload_bytes": payloadLen, "duration": fmt.Sprintf("%s", time.Since(start))})

}

// deletes a file if the configuration is set to cleanup = true
func cleanup(configuration *configs.ViperConfiguration, filename string) error {
	if configuration.GetBool("cleanup") {
		if utils.CheckFileExists(filename) {
			err := os.Remove(filename)
			if err != nil {
				return err
			}
			log.Println("cleanup: file deleted")
		}
	}
	return nil
}
