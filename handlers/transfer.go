package handlers

import (
	"aif/arweave"
	"aif/configs"
	"aif/utils"
	"aif/utils/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
	"strings"
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
	out, err := exec.Command("ipfs", "get", ipfsHash).CombinedOutput()
	log.Printf("Output from ipfs cmd %s", string(out))
	if strings.Contains(string(out), "merkledag: not found") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file not found."})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("file retrieved successfully from IPFS")

	// uploading it to arweave
	txID, err := arweave.Transfer(ipfsHash, ipfsHash, configuration)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("transfer finished successfully")

	err = cleanup(configuration, ipfsHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "transfer completed but I couldn't cleanup the file " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": txID, "duration": string(time.Since(start))})
}

func cleanup(configuration *configs.ViperConfiguration, ipfsHash string) error {
	if configuration.GetBool("cleanup") {
		err := os.Remove(ipfsHash)
		if err != nil {
			return err
		}
		log.Println("cleanup: file deleted")
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
