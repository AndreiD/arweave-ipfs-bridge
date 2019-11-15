package arweave

import (
	"aif/configs"
	"context"
	"github.com/Dev43/arweave-go/api"
)

// GetBalance of Arweave tokens
func GetBalance(address string, configuration *configs.ViperConfiguration) (string, error) {
	client, err := api.Dial(configuration.Get("nodeURL"))
	if err != nil {
		return "", err
	}
	balance, err := client.GetBalance(context.Background(), address)
	if err != nil {
		return "", err
	}

	return balance, nil
}
