package utils

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Close error checking for defer close
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// GetRequest executes a GET request
func GetRequest(url string) ([]byte, error) {

	client := http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	reqWithDeadline := req.WithContext(ctx)
	response, err := client.Do(reqWithDeadline)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(response.Body)

	return data, err

}


