package arweave

import (
	"aif/configs"
	"aif/utils"
	"aif/utils/log"
	"bytes"
	"context"
	"os"
)

// original from: https://github.com/Dev43/arweave-go

// Transfer on arweave blockchain. Returns transaction hash, error
func Transfer(ipfsHash string, fileName string, configuration *configs.ViperConfiguration) (string, error) {

	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer utils.Close(f)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(f)
	if err != nil {
		return "", err
	}

	// creates a new batch
	ar, err := NewTransactor(configuration.Get("nodeURL"))
	if err != nil {
		return "", err
	}

	arWallet := NewWallet()
	err = arWallet.LoadKeyFromFile(configuration.Get("walletFile"))
	if err != nil {
		return "", err
	}

	txBuilder, err := ar.CreateTransaction(ipfsHash, context.Background(), arWallet, "0", buf.Bytes(), "")
	if err != nil {
		return "", err
	}

	// sign the transaction
	txn, err := txBuilder.Sign(arWallet)
	if err != nil {
		return "", err
	}

	// send the transaction
	resp, err := ar.SendTransaction(context.Background(), txn)
	if err != nil {
		return "", err
	}

	log.Printf("arweave node responded %s", resp)

	return txn.Hash(), nil
}
