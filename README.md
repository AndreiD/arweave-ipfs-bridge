# :bridge_at_night: IPFS Arweave Bridge 

[![License](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/AndreiD/arweave-ipfs-bridge/blob/master/LICENSE)

A bridge to connect IPFS to Arweave

<p align="center">
  <p align="center">Click on the image to enlarge it</p>
   <img alt="how it looks" height="500" src="https://raw.githubusercontent.com/AndreiD/arweave-ipfs-bridge/master/assets/postman_example.png">
 </p>

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
- API will run on **http://localhost:5555**

* If you want to use it on other platforms, either use the Dockerfile or build it yourself

### How it does the transfer

##### POST /api/transfer

> example: /api/transfer


output:
~~~~
{
    "duration": "1.145424997s",
    "id": "Qq2IE63kxOaZp7laJapQqyu1PUPPktGoSSyKTLI9iP4",
    "payload_bytes": 14
}
~~~~

POST a raw json body containing

{
    "ipfs_hash": "THE_IPFS_HASH",
    "use_compression": BOOLEAN,
    "tags": [
        {
            "key": "key1",
            "value": "value1"
        },...
    ]
}

Where: 
- **hash** the IPFS hash
- **use_compression** if you want the file compressed before sending it to arweave. (zip is used)
- **tags** (optional) if you want to add your own tags. (Note: the IPFS-Add tag is added automatically)


example:
~~~~
{
    "ipfs_hash": "Qmc5gCcjYypU7y28oCALwfSvxCBskLuPKWpK4qpterKC7z",
    "use_compression": false,
    "tags": [
        {
            "key": "key1",
            "value": "value1"
        },
        {
            "key": "key2",
            "value": "value2"
        }
    ]
}
~~~~

response:
- **duration** gives you the time it took to get it from IPFS and to upload it to Arweave
- **id** represents the arweave transaction id
- **Attention** The transaction is NOT mined yet. You can get the status of a transaction by calling the API below

(it's up to the user how often they pool for the transaction status)

Response:

~~~~
{
    "duration": "1.145424997s",
    "id": "Qq2IE63kxOaZp7laJapQqyu1PUPPktGoSSyKTLI9iP4",
    "payload_bytes": 14
}
~~~~


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
- **cleanup** deletes the files after they get retrieved. use false if you want to debug things
- **server's ip/port**: how to run this API service

about IPFS gateway:

choices: 
- http://127.0.0.1:8080/ipfs/
- https://ipfs.infura.io/ipfs/
- https://cloudflare-ipfs.com/ipfs/
- https://ipfs.io/ipfs/

What if you want to host your own gateway ?

make sure you crete an api that will just accept a ipfs hash after the "/"


## Other API Endpoints

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
   

#### :point_right: Retrieve a transaction or check if it's pending
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




#### :point_right: Post content directly to Arweave

##### POST /api/transfer_arweave

with:
~~~~
{
    "payload": "CXR4SUQsIHBheWxvYWRMZW4sIGVyciA6PSBhcndlYXZlLlRyYW5zZmVyRGlyZWN0bHlBcndlYXZlKG5UcmFuc2Zlci5QYXlsb2FkLCBhclRhZ3MsIGNvbmZpZ3VyYXRpb24pCglpZiBlcnIgIT0gbmlsIHsKCQljLkpTT04oaHR0cC5TdGF0dXNCYWRSZXF1ZXN0LCBnaW4uSHsiZXJyb3IiOiBlcnIuRXJyb3IoKX0pCgkJcmV0dXJuCgl9CgoJbG9nLlByaW50ZigiVHJhbnNmZXIgdG8gQXJ3ZWF2ZSBmaW5pc2hlZCBzdWNjZXNzZnVsbHkuIFR4IElEICVzIiwgdHhJRCkKCgljLkpTT04oaHR0cC5TdGF0dXNPSywgZ2luLkh7ImlkIjogdHhJRCwgInBheWxvYWRfYnl0ZXMiOiBwYXlsb2FkTGVuLCAiZHVyYXRpb24iOiBmbXQuU3ByaW50ZigiJXMiLCB0aW1lLlNpbmNlKHN0YXJ0KSl9KQ==",
    "tags": [
        {
            "key": "key1",
            
           "value": "value1"
        },
        {
            "key2": "value2"
        }
    ]
}
~~~~

Where:

- **decoded** (true or false) if you want it decoded or not
- **extract** (true or false) if you posted zipped, it will unzip it 




### Special thanks to:

https://github.com/Dev43/arweave-go -> for the transaction signing & sending code

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
