package arweave

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// original from: https://github.com/Dev43/arweave-go

// Client struct
type Client struct {
	client *http.Client
	url    string
}

// Dial creates a new arweave client
func Dial(url string) (*Client, error) {
	return &Client{client: new(http.Client), url: url}, nil
}

// TxAnchor .
func (c *Client) TxAnchor(ctx context.Context) (string, error) {
	body, err := c.get(ctx, "tx_anchor")
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// LastTransaction requests the last transaction of an account
func (c *Client) LastTransaction(ctx context.Context, address string) (string, error) {
	body, err := c.get(ctx, fmt.Sprintf("wallet/%s/last_tx", address))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// GetReward requests the current network reward
func (c *Client) GetReward(ctx context.Context, data []byte) (string, error) {
	body, err := c.get(ctx, fmt.Sprintf("price/%d", len(data)))
	if err != nil {
		return "", err
	}
	return string(body), nil

}

// Commit sends a transaction to the weave with a context
func (c *Client) Commit(ctx context.Context, data []byte) (string, error) {
	body, err := c.post(ctx, "tx", data)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func getResponse(resp io.ReadCloser, returnedError error) ([]byte, error) {
	if resp != nil {
		defer resp.Close()
	}
	if returnedError != nil {
		return handleHTTPError(resp, returnedError)
	}

	b, err := ioutil.ReadAll(resp)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func handleHTTPError(resp io.Reader, returnedError error) ([]byte, error) {
	if resp != nil {
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(resp); err == nil {
			return nil, fmt.Errorf("%v %v", returnedError, buf.String())
		}
	}
	return nil, returnedError
}

func (c *Client) requestWithContext(ctx context.Context, method string, url string, body []byte) (io.ReadCloser, error) {
	req, err := http.NewRequest(method, url, ioutil.NopCloser(bytes.NewReader(body)))
	if err != nil {
		return nil, err
	}
	reqWithContext := req.WithContext(ctx)
	reqWithContext.ContentLength = int64(len(body))
	if method == "POST" {
		reqWithContext.Header.Set("Content-type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.Body, errors.New(resp.Status)
	}
	return resp.Body, nil
}

func (c *Client) post(ctx context.Context, endpoint string, body []byte) ([]byte, error) {
	resp, err := c.requestWithContext(ctx, "POST", c.formatURL(endpoint), body)
	return getResponse(resp, err)
}

func (c *Client) get(ctx context.Context, endpoint string) ([]byte, error) {
	resp, err := c.requestWithContext(ctx, "GET", c.formatURL(endpoint), nil)
	return getResponse(resp, err)
}

func (c *Client) formatURL(endpoint string) string {
	return fmt.Sprintf("%s/%s", c.url, endpoint)
}
