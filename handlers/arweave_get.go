package handlers

import (
	"aif/configs"
	"aif/utils"
	"aif/utils/log"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
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

	extract := c.Query("extract")
	if extract == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please provide the extract in the query ex: ?extract=false"})
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

	if extract == "true" {
		// save output to file
		f, err := os.Create(txID + ".zip")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if decode == "true" {
			output, err = utils.DecodeString(string(output))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
		_, err = f.Write(output)
		utils.Close(f)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		defer func() {
			err := cleanup(configuration, txID+".zip")
			if err != nil {
				log.Error(err)
			}
		}()

		archivedFilename, err := utils.UnarchiveIt(txID+".zip", ".")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//read the unzipped file
		b, err := ioutil.ReadFile(archivedFilename)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		defer func() {
			err = cleanup(configuration, txID+".zip")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "I couldn't cleanup the file " + err.Error()})
				return
			}
			err = cleanup(configuration, archivedFilename)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "I couldn't cleanup the file " + err.Error()})
				return
			}
		}()

		log.Printf("done in %s", time.Since(start))
		c.String(http.StatusOK, string(b))
		return
	}

	if decode == "true" {
		decoded, err := utils.DecodeString(string(output))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = cleanup(configuration, txID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "I couldn't cleanup the file " + err.Error()})
			return
		}

		log.Printf("done in %s", time.Since(start))
		c.String(http.StatusOK, string(decoded))
		return
	}

	// if there's no compression & no decoding

	log.Printf("done in %s", time.Since(start))
	c.String(http.StatusOK, string(output))
}
