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
	ipfsHash := c.Query("hash")
	if ipfsHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please provide the IPFS hash in the query ex: ?hash=Qmc5gCcjYypU7y28oCALwfSvxCBskLuPKWpK4qpterKC7z"})
		return
	}

	useCompression := c.Query("use_compression")
	if useCompression == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please provide the use_compression parameter in the query ex: ?use_compression=true"})
		return
	}

	err := utils.CheckValidityIPFSHash(ipfsHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	configuration, ok := c.MustGet("configuration").(*configs.ViperConfiguration)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the configuration"})
		return
	}

	log.Printf("starting retrieving from IPFS %s", ipfsHash)

	// get from IPFS
	out, statusCode, err := utils.GetRequest(configuration.Get("ipfsGateway") + ipfsHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("IPFS gateway returned status code %d", statusCode)

	// saving it to the filesystem
	f, err := os.Create(ipfsHash)
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

	// uploading it to Arweave
	txID, payloadLen, err := arweave.Transfer(ipfsHash, useCompression, ipfsHash, configuration)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Transfer to Arweave finished successfully. Tx ID %s", txID)

	err = cleanup(configuration, ipfsHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "transfer completed but I couldn't cleanup the file " + err.Error()})
		return
	}

	err = cleanup(configuration, ipfsHash+".zip")
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

// TODO: what to do with this ....oO
// checks for a more descriptie error
//func parseArweaveError(payload string) string {
//	if strings.Contains(payload, "208") {
//		return "transaction has already been submitted"
//	}
//	if strings.Contains(payload, "400") {
//		return "the transaction is invalid, couldn't be verified, or the arweave does not have suffucuent funds"
//	}
//	if strings.Contains(payload, "429") {
//		return "the request has exceeded the clients rate limit quota"
//	}
//	if strings.Contains(payload, "500") {
//		return "the nodes was unable to verify the transaction"
//	}
//	return payload
//}
