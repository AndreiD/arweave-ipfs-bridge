package arweave

import (
	"aif/configs"
	"github.com/Dev43/arweave-go/batchchunker"
	"github.com/Dev43/arweave-go/transactor"
	"os"
)

// Transfer on Arweave blockchain. Returns transaction hash,  error
func Transfer(fileName string, configuration *configs.ViperConfiguration) ([]string, error) {

	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	// creates a new batch
	ar, err := transactor.NewTransactor(configuration.Get("nodeURL"))
	if err != nil {
		return nil, err
	}
	newB := batchchunker.NewBatch(ar, myWallet, f, info.Size())

	// sends all the transactions
	list, err := newB.SendBatchTransaction()
	if err != nil {
		return nil, err
	}

	return list, nil
}
