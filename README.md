# Authenticated Subgraph Matching in Hybrid-Storage Blockchains

## Introduction

Subgraph matching is to find all the subgraphs isomorphic with a given query graph in a data graph. The data owner could delegate their data graphs to the Tamper-proof blockchain due to the limited storage and computation. In the blockchain, for improving the scalability, it is a good choice to store the raw data in a off-chain storage service provider(SP) and only maintains the digest of the raw data on-chain by the smart contracts. However, the SP is untrusted, it may return uncorrect or uncompleteness results. To support integrity-assured data query in such scenario, authenticated subgraph matching can be applied.  To our knowledge, there is no work to enable the blockchain to support graph data queries. 

In this paper, we study the novel approach to support authenticated subgraph matching query for the large graph kept off-chain. We first design the ADS structure as MELTree for labels and keep the digests of the roots on-chain. We propose the verification object (VO) construction algorithm as AMatching, for subgraph matching queries in order to ensure the completeness and soundness of the results. To further reduce the cost, we propose the approach as AMatching* based on two-way search including forward search and reverse search. Forward search aims to identify the isomorphic
subgraphs. When an edge or a vertex cannot generate isomorphic subgraphs, AMatching* aggregates all the related edges and vertices for un isomorphic subgraphs without enumerating the corresponding partial matching. Moreover, we further optimize the on-chain storage cost by proposing MVPTree , in which it organizes the structures for vertices, and only needs to keep one root digest on-chain for verification. Experimental results show that, the proposed algorithms and the optimizations improve the performance by 1-2 orders of magnitude.

## Environment

Ethereum blockchain platform, Golang 1.15.5.

1. devDependencies

   ```
   CentOS:
   yum install git wget bzip2 vim gcc-c++ ntp epel-release
   nodejs cmake -y
   yum update
   Ubuntu:
   sudo apt install make
   sudo apt install g++
   sudo apt-get install libltdl-dev
   apt install -y build-essential
   ```

2. go-ethereum install

   ```
   git clone https://github.com/ethereum/go-ethereum
   cd go-ethereum
   make geth
   # Copy the finished go-ethereum /buil/bin/geth executable to /usr/local/bin
   ```

3. Create Genesis Block

    Create a new genesis.json file and write the following information

   ```
   {
   "config": {
   "chainId": 981106,
   "homesteadBlock": 0,
   "eip150Block": 0,
   "eip150Hash":
   "0x0000000000000000000000000000000000000000000000000000000000
   000000",
   "eip155Block": 0,
   "eip158Block": 0,
   "byzantiumBlock": 0,
   "constantinopleBlock": 0,
   "petersburgBlock": 0,
   "ethash": {}
   
   },
   "nonce": "0x0",
   "timestamp": "0x284d29c0",
   "extraData":
   "0x0000000000000000000000000000000000000000000000000000000000
   000000",
   "gasLimit": "0x47b760",
   "difficulty": "0x80000",
   "mixHash":
   "0x0000000000000000000000000000000000000000000000000000000000
   000000",
   "coinbase": "0x0000000000000000000000000000000000000000",
   "alloc": {
   "0000000000000000000000000000000000000000": {
   "balance": "0x1"
   }
   },
   "number": "0x0",
   "gasUsed": "0x0",
   "parentHash":
   "0x0000000000000000000000000000000000000000000000000000000000
   000000"
   }
   ```

## Test

1. Enabling the Network

   ```
   geth --datadir data init genesis.json
   geth --datadir data --networkid 981106 --http --
   http.corsdomain "*" --http.port 8545 --http.addr 0.0.0.0 --
   nodiscover console --allow-insecure-unlock
   ```

2. Open a new shell to mine

   ```
   geth attach ipc:geth.ipc
   miner.start()
   ```

3. Open a new shell to test

   ```
   cd src
   go run main.go
   ```

4. Remember to turn off the network after use, press ctrl+d or type 'exit'.
