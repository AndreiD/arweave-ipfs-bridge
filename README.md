# :zap: AIB :zap: Arweave IPFS Bridge 

A bridge to connect IPFS to Arweave

### Features:

- only 3 GET requests
- you easily split it into multiple services
- load balance it, use multiple wallets not just one
- easy to integrated with almost anything

### How to use it ?

Tested on ubuntu 19.04

- start and configure ipfs to your liking
- start hooverd and connect it to a wallet
- get the binary file aib
- copy the configuration.json file in the same directory (modify it to your liking)
- run ./aib **defaults on 0.0.0.0:5555/api**


## API Endpoints

#### GET /api/ping 
   
~~~~
pong
~~~~
   
#### GET /api/balance?wallet=YOUR_WALLET_HERE

> example: /balance?wallet=qGwglm54w6I9-CCcNSAjvWzqGNZfb0zAUNkXYVYN5LY

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
    "hash": "g9e6nzaiz74-RTCiXJwmOvQLtExT-wlx5oiC4ybqTtQ",
    "output": "Transaction g9e6nzaiz74-RTCiXJwmOvQLtExT-wlx5oiC4ybqTtQ dispatched to arweave.net:443 with response: 200.\n"
}
~~~~

Where:

- duration gives you the time it took to get it from IPFS and to upload it to Arweave
- hash represents the arweave transaction hash
- output: represents the output of hooverd

   
      
#### Why it's a REST API




### Bugs / Features / Questions

fell free to create an issue

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details.



## License

AIB is released under the MIT license.