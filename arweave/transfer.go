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
func Transfer(ipfsHash string, useCompression bool, tags []Tag, configuration *configs.ViperConfiguration) (string, int, error) {

	buf := new(bytes.Buffer)

	if useCompression {
		log.Println("compression is activated")
		err := utils.ArchiveFile(ipfsHash)
		if err != nil {
			return "", -1, err
		}

		f, err := os.Open(ipfsHash + ".zip")
		if err != nil {
			return "", -1, err
		}
		defer utils.Close(f)
		_, err = buf.ReadFrom(f)
		if err != nil {
			return "", -1, err
		}

	} else {
		log.Warn("for files bigger than 200 bytes, please consider using compression")
		f, err := os.Open(ipfsHash)
		if err != nil {
			return "", -1, err
		}
		defer utils.Close(f)

		_, err = buf.ReadFrom(f)
		if err != nil {
			return "", -1, err
		}
	}

	ar, err := NewTransactor(configuration.Get("nodeURL"))
	if err != nil {
		return "", -1, err
	}

	arWallet := NewWallet()
	err = arWallet.LoadKeyFromFile(configuration.Get("walletFile"))
	if err != nil {
		return "", -1, err
	}

	log.Printf("creating a transaction with a payload of %d bytes", buf.Len())

	txBuilder, err := ar.CreateTransaction(context.Background(), ipfsHash, tags, arWallet, "0", buf.Bytes(), "")
	if err != nil {
		return "", -1, err
	}

	// sign the transaction
	txn, err := txBuilder.Sign(arWallet)
	if err != nil {
		return "", -1, err
	}

	// send the transaction
	resp, err := ar.SendTransaction(context.Background(), txn)
	if err != nil {
		return "", -1, err
	}

	log.Printf("arweave node responded %s", resp)

	return txn.Hash(), buf.Len(), nil
}

// TransferDirectlyArweave on arweave blockchain. Returns transaction hash, error
func TransferDirectlyArweave(payload string, tags []Tag, configuration *configs.ViperConfiguration) (string, int, error) {

	buf := new(bytes.Buffer)

	ar, err := NewTransactor(configuration.Get("nodeURL"))
	if err != nil {
		return "", -1, err
	}

	arWallet := NewWallet()
	err = arWallet.LoadKeyFromFile(configuration.Get("walletFile"))
	if err != nil {
		return "", -1, err
	}

	log.Printf("creating a transaction with a payload of %d bytes", buf.Len())

	txBuilder, err := ar.CreateTransactionArweave(context.Background(), tags, arWallet, "0", []byte(payload), "")
	if err != nil {
		return "", -1, err
	}

	// sign the transaction
	txn, err := txBuilder.Sign(arWallet)
	if err != nil {
		return "", -1, err
	}

	// send the transaction
	resp, err := ar.SendTransaction(context.Background(), txn)
	if err != nil {
		return "", -1, err
	}

	log.Printf("arweave node responded %s", resp)

	return txn.Hash(), buf.Len(), nil
}
