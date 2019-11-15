package handlers

import (
	"aif/configs"
	"aif/utils/log"
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "please provide the ipfs hash in the query ex: ?hash=Qmc5gCcjYypU7y28oCALwfSvxCBskLuPKWpK4qpterKC7z"})
		return
	}

	err := checkValidityIPFSHash(ipfsHash)
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
	payload := []byte("@" + ipfsHash)
	out, err = postToArweave("http://localhost:1908/raw", ipfsHash, bytes.NewBuffer(payload))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("got output from hooverd %s", string(out))

	if !strings.Contains(string(out), "response: 200") {
		c.JSON(http.StatusBadRequest, gin.H{"error": parseArweaveError(string(out))})
		return
	}

	if len(string(out)) < 56 {
		c.JSON(http.StatusBadRequest, gin.H{"error": string(out)})
		return
	}

	hashArweave := string(out)[12:55]

	log.Println("transfer finished successfully")

	err = cleanup(configuration, ipfsHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "transfer completed but I couldn't cleanup the file " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"output": string(out), "hash": hashArweave, "duration": fmt.Sprintf("%s", time.Since(start))})
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

// checks for a more descriptie error
func parseArweaveError(payload string) string {
	if strings.Contains(payload, "208") {
		return "transaction has already been submitted"
	}
	if strings.Contains(payload, "400") {
		return "the transaction is invalid, couldn't be verified, or the wallet does not have suffucuent funds"
	}
	if strings.Contains(payload, "429") {
		return "the request has exceeded the clients rate limit quota"
	}
	if strings.Contains(payload, "500") {
		return "the nodes was unable to verify the transaction"
	}
	return payload
}

// check if the passed hash is correct size & starts with Qm
func checkValidityIPFSHash(ipfsHash string) error {
	if len(ipfsHash) != 46 {
		return fmt.Errorf("it doesn't look like an IPFS hash. not 46 characters")
	}
	if !strings.HasPrefix(ipfsHash, "Qm") {
		return fmt.Errorf("it doesn't look like an IPFS hash. doesn't start with Qm")
	}
	return nil
}

func postToArweave(url string, ipfsHash string, payload io.Reader) ([]byte, error) {

	client := http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req.Header.Set("Content-Type", "text/html")
	req.Header.Set("IPFS-Add", ipfsHash)

	reqWithDeadline := req.WithContext(ctx)
	response, err := client.Do(reqWithDeadline)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
