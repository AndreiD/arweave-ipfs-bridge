package utils

import (
	"aif/utils/log"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// Close error checking for defer close
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// GetRequest executes a generic GET request
func GetRequest(url string) ([]byte, int, error) {

	client := http.Client{Timeout: 180 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, -1, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	reqWithDeadline := req.WithContext(ctx)
	response, err := client.Do(reqWithDeadline)
	if err != nil {
		return nil, -1, err
	}

	data, err := ioutil.ReadAll(response.Body)

	return data, response.StatusCode, err

}

// EncodeToBase64 encodes a byte array to base64 raw url encoding
func EncodeToBase64(toEncode []byte) string {
	return base64.RawURLEncoding.EncodeToString(toEncode)
}

// DecodeString decodes from base64 raw url encoding to byte array
func DecodeString(toDecode string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(toDecode)
}

// CheckValidityIPFSHash - checks if the passed hash is correct size & starts with Qm
func CheckValidityIPFSHash(ipfsHash string) error {
	if len(ipfsHash) != 46 {
		return fmt.Errorf("it doesn't look like an IPFS hash. not 46 characters")
	}
	if !strings.HasPrefix(ipfsHash, "Qm") {
		return fmt.Errorf("it doesn't look like an IPFS hash. doesn't start with Qm")
	}
	return nil
}

// CheckValidityArweaveTxID - checks if an id is an Arweave ID...at least in length
func CheckValidityArweaveTxID(id string) error {
	if len(id) != 43 {
		return fmt.Errorf("this doesn't look like a transaction id from arweave. length is not 43")
	}
	return nil
}

// CheckFileExists - check if file exists
func CheckFileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		log.Errorf("can't verify if %s exists or not", filename)
		return false
	}
}
