package handlers

import (
	"aif/arweave"
	"aif/configs"
	"aif/utils/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// TransferToArweave .
func TransferToArweave(c *gin.Context) {

	start := time.Now()

	type tags struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	type newTransfer struct {
		Payload        string `json:"payload" binding:"required"`
		Tags           []tags `json:"tags"`
	}

	var nTransfer newTransfer

	err := c.BindJSON(&nTransfer)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload " + err.Error()})
		return
	}

	configuration, ok := c.MustGet("configuration").(*configs.ViperConfiguration)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the configuration"})
		return
	}

	var arTags []arweave.Tag
	for _, tag := range nTransfer.Tags {
		arTags = append(arTags, arweave.Tag{
			Name:  tag.Key,
			Value: tag.Value,
		})
	}

	// uploading it to Arweave
	txID, payloadLen, err := arweave.TransferDirectlyArweave(nTransfer.Payload, arTags, configuration)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Transfer to Arweave finished successfully. Tx ID %s", txID)

	c.JSON(http.StatusOK, gin.H{"id": txID, "payload_bytes": payloadLen, "duration": fmt.Sprintf("%s", time.Since(start))})

}
