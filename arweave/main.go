package arweave

import (
	"aif/configs"
	"aif/utils/log"
	"github.com/Dev43/arweave-go/wallet"
)

var myWallet *wallet.Wallet

// Initialize initializes a new arweave
func Initialize(config *configs.ViperConfiguration) {
	myWallet = wallet.NewWallet()
	err := myWallet.LoadKeyFromFile(config.Get("walletFile"))
	if err != nil {
		log.Panic(err)
	}
}
