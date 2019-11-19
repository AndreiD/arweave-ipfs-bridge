# :bridge_at_night: IPFS Arweave Bridge 

[![License](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/AndreiD/arweave-ipfs-bridge/blob/master/LICENSE)

A bridge to connect IPFS to Arweave

### Features

- only couple of requests
- handles compression / extractions for you (optional)
- easy to integrated with almost anything (logging, monitoring etc.)
- no ipfs running in the background needed, no hooverd needed

### How to use it (Ubuntu only)

Tested on ubuntu 19.04

- Choose an IPFS gateway to your liking (or create one yourself)
- If you use Ubuntu, get the binary file **iab**, if not you'll have to build/run it
- Copy the configuration.json file in the same directory (modify it to your liking)
- run ./iab  
- API is live on **http://localhost:5555**

* If you want to use it on other platforms, either use the Dockerfile or build it yourself

### Docker

checkout the Dockerfile

~~~~
docker build -t iab:1.0 .
docker run -p 5555:5555 iab:1.0
~~~~

### Build it

if you want to build it you need go >= 1.12
in the root directory run: 

~~~~
go build -o YOUR_BINARY_NAME
./YOUR_BINARY_NAME
~~~~

for installing go check: https://golang.org/doc/install
it should build & run without problems on macOS & Windows

#### Configuration file

~~~~
{
  "debug": true,
  "nodeURL": "https://arweave.net",
  "ipfsGateway": "https://ipfs.infura.io/ipfs/",
  "walletFile": "/path_to_your_wallet/arweave-wallet.json",
  "cleanup": true,
  "server": {
    "host": "0.0.0.0",
    "port": 5555
  }
}
~~~~

Where:

- **debug** runs with much more verbosity
- **nodeURL** how to connect to arweave
- **ipfsGateway** - your IPFS gateway. See below
- **walletFile** the json AR wallet file location
- **cleanup** deletes the files after they get retrieved
- **server's ip/port**: how to run this API service

about IPFS gateway:

choices: 
- http://127.0.0.1:8080/ipfs/
- https://ipfs.infura.io/ipfs/
- https://cloudflare-ipfs.com/ipfs/
- https://ipfs.io/ipfs/

What if you want to host your own gateway ?

make sure you crete an api that will just accept a ipfs hash after the "/"


## API Endpoints

#### :point_right: Ping
##### GET /api/ping 
   
used to check if the service is alive

output:
~~~~
pong
~~~~
   
#### :point_right: Check your AR Token balance
##### GET /api/balance

returns your balance of AR Tokens

output:
~~~~
{
    "ar": "0.9996836928",
    "winston": "999683692828"
}
~~~~
   
#### :point_right: Transfer from IPFS to Arweave
##### GET /api/transfer?hash=IPFS_HASH_HERE&use_compression=BOOLEAN

> example: /api/transfer?hash=QmUNXr47Bja3aHUMfhXX5mMWTFJKuoUGETcA48vHG7dhag&use_compression=true


output:
~~~~
{
    "duration": "953.239662ms",
    "id": "g9e6nzaiz74-RTCiXJwmOvQLtExT-wlx5oiC4ybqTtQ"
}
~~~~

Where:

query parameters:
- **hash** the IPFS hash
- **use_compression** if you want the file compressed before sending it to arweave. (zip is used)


response:
- **duration** gives you the time it took to get it from IPFS and to upload it to Arweave
- **id** represents the arweave transaction id
- **Attention** The transaction is NOT mined yet. You can get the status of a transaction by calling the API below

   
#### :point_right: Retrive a transaction or check if it's pending
##### GET /api/check_tx_arweave?transaction_id=TRANSACTION_ID

> example: /api/check_tx_arweave?transaction_id=bnRQhVkook_lPv8uxuDRcj-wC5R2nfVps-2qA6-81WU

output:
~~~~
the transaction details or it's status (ex: pending)
~~~~


### Other helper API calls that you might need

#### :point_right: Retrieve a file from IPFS

##### GET /api/ipfs?hash=IPFS_HASH

> example: /api/ipfs?hash=QmbRmU9vYwH9Hhn1eH1WEFVS9sugpGSdJrfqtuZ329EgZA

output:
~~~~
content of the file from IPFS
~~~~

#### :point_right: Retrieve a file from Arweave

##### GET /api/arweave?transaction_id=TRANSACTION_ID&decode=true

> example: /api/arweave?transaction_id=GyrTvuUBK9AgVLGBA8SsOHkUYmWApNqvJtMjJZZIvbQ&decode=BOOLEAN&extract=true

output:
~~~~
content of the file from Arweave
~~~~

Where:

- **decoded** (true or false) if you want it decoded or not
- **extract** (true or false) if you posted zipped, it will unzip it 

### Special thanks to:

https://github.com/Dev43/arweave-go -> for the transaction signing & transmitting code

### Bugs / Features / Questions

fell free to create an issue

### TODO://

- [x] compression
- [x] docker
- [x] golang-ci yml
- [ ] awaiting your idea


## Project Values

This project has a few principle-based goals that guide its development:

- **Limit dependencies.** Keep the package lightweight.

- **Pure Go.** This means no cgo or other external/system dependencies. This package should be able to stand on its own and cross-compile easily to any platform -- and that includes its library dependencies.

- **Idiomatic Go.** Keep interfaces small, variable names semantic, vet shows no errors, the linter (golangci lint) is generally quiet, etc.

- **Well-documented.** Use comments prudently; explain why non-obvious code is necessary (and use tests to enforce it). Keep the docs updated, and have examples where helpful.

- **Keep it efficient.** This often means keep it simple. Fast code is valuable.

- **Consensus.** Contributions should ideally be approved by multiple reviewers before being merged. Generally, avoid merging multi-chunk changes that do not go through at least one or two iterations/reviews. Except for trivial changes, PRs are seldom ready to merge right away.

- **Have fun contributing.** Coding is awesome!

I welcome contributions and appreciate your efforts! However, please open issues to discuss any changes before spending the time preparing a pull request. 
This will save time, reduce frustration, and help coordinate the work. Thank you!


## License

IPFS Arweave Bridge is released under the MIT license.
