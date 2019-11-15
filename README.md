# :zap: IAB :zap: IPFS Arweave Bridge 

[![License](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/AndreiD/arweave-ipfs-bridge/blob/master/LICENSE)

A bridge to connect IPFS to Arweave

### Features

- only 3 GET requests
- you easily split it into multiple services
- load balance it, use multiple wallets not just one
- easy to integrated with almost anything

### How to use it

Tested on ubuntu 19.04

- start and configure ipfs to your liking
- start hooverd and connect it to a wallet
- get the binary file ***iab***
- copy the configuration.json file in the same directory (modify it to your liking)
- run ./iab **defaults on 0.0.0.0:5555**

### Build it

if you want to build it you need go >= 1.12
in the root directory run: 

~~~~
go build -o YOUR_BINARY_NAME
~~~~

for installing go check: https://golang.org/doc/install
it should run without problems also on macOS & Windows

### Configuration file

~~~~
{
  "debug": true,
  "nodeURL": "https://arweave.net",
  "walletFile": "/path_to_your_wallet/arweave-wallet.json",
  "cleanup": true,
  "server": {
    "host": "0.0.0.0",
    "port": 5555
  }
}
~~~~

Where:

- debug: runs with more verbosity
- nodeURL: how to connect to arweave ?
- walletFile: the json AR wallet file
- cleanup: deletes the files after they get retrived
- server host/port: how to run this service


## API Endpoints

#### GET /api/ping 
   
used to check if the service is alive

~~~~
pong
~~~~
   
#### GET /api/balance

return your balance of AR Tokens

~~~~
{
    "ar": "0.9996836928",
    "winston": "999683692828"
}
~~~~
   
#### GET /api/transfer?hash=IPFS_HASH_HERE

> example: /api/transfer?hash=QmUNXr47Bja3aHUMfhXX5mMWTFJKuoUGETcA48vHG7dhag

~~~~
{
    "duration": "953.239662ms",
    "id": "g9e6nzaiz74-RTCiXJwmOvQLtExT-wlx5oiC4ybqTtQ"
}
~~~~

Where:

- duration gives you the time it took to get it from IPFS and to upload it to Arweave
- id represents the arweave transaction id

**Attention**

The transaction ID is not mined yet. You can get the status of a transaction by calling the API below

   
#### GET /api/check_tx_arweave?transaction_id=TRANSACTION_ID

> example: /api/check_tx_arweave?transaction_id=bnRQhVkook_lPv8uxuDRcj-wC5R2nfVps-2qA6-81WU

~~~~
the transaction details or it's status (ex: pending)
~~~~


### Other helper API calls that you might need

#### GET /api/ipfs?hash=IPFS_HASH

> example: /api/ipfs?hash=QmbRmU9vYwH9Hhn1eH1WEFVS9sugpGSdJrfqtuZ329EgZA

~~~~
content of the file from IPFS
~~~~

#### GET /api/arweave?transaction_id=TRANSACTION_ID&decode=true

> example: /api/arweave?transaction_id=GyrTvuUBK9AgVLGBA8SsOHkUYmWApNqvJtMjJZZIvbQ&decode=true

~~~~
content of the file from Arweave
~~~~

Where:

- decoded: if you want it decoded or not

### Special thanks to:

https://github.com/Dev43/arweave-go -> for the transaction signing & transmitting code

### Bugs / Features / Questions

fell free to create an issue

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details.



## License

AIB is released under the MIT license.