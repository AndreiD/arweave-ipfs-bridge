package arweave

import (
	"context"
	"encoding/json"
	"errors"
)

// original from: https://github.com/Dev43/arweave-go

// ClientCaller is the base interface needed to create a Transactor
type ClientCaller interface {
	TxAnchor(ctx context.Context) (string, error)
	LastTransaction(ctx context.Context, address string) (string, error)
	GetReward(ctx context.Context, data []byte) (string, error)
	Commit(ctx context.Context, data []byte) (string, error)
}

// Transactor type, allows one to create transactions
type Transactor struct {
	Client ClientCaller
}

// NewTransactor creates a new arweave transactor. You need to pass in a context and a url
// If sending an empty string, the default url is https://arweave.net
func NewTransactor(nodeURL string) (*Transactor, error) {
	c, err := Dial(nodeURL)
	if err != nil {
		return nil, err
	}

	return &Transactor{
		Client: c,
	}, nil
}

// CreateTransaction creates a brand new transaction
func (tr *Transactor) CreateTransaction(ctx context.Context, ipfsHash string, w WalletSigner, amount string, data []byte, target string) (*Transaction, error) {
	lastTx, err := tr.Client.TxAnchor(ctx)
	if err != nil {
		return nil, err
	}

	price, err := tr.Client.GetReward(ctx, data)
	if err != nil {
		return nil, err
	}

	// Non encoded transaction fields
	tx := NewTransaction(
		ipfsHash,
		lastTx,
		w.PubKeyModulus(),
		amount,
		target,
		data,
		price,
	)

	return tx, nil
}

// SendTransaction formats the transactions (base64url encodes the necessary fields)
// marshalls the Json and sends it to the arweave network
func (tr *Transactor) SendTransaction(ctx context.Context, tx *Transaction) (string, error) {
	if len(tx.Signature()) == 0 {
		return "", errors.New("transaction is missing signature")
	}
	serialised, err := json.Marshal(tx)
	if err != nil {
		return "", err
	}
	return tr.Client.Commit(ctx, serialised)
}
