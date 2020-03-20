## kahf
[![Build Status](https://travis-ci.com/HalalChain/qitmeer.svg?token=DzCFNC6nhEqPc89sq1nd&branch=master)](https://travis-ci.com/HalalChain/qitmeer)

Kahf is wallet main program for Qitmeer

#### Build
```bash
go build
```
#### How to use

There has two ways to using the wallet

1. Command line

    command format
    
    ```bash
    ./kahf <cmd> [params]
    ```
    support commands
    
    ```bash
    ./kahf <help|tx|out|addr|token|status|peer> [params]
    ```
          
    examples
    
    ```bash
   
    # get helper
    ./kahf help
   
    # generate a wallet address
    ./kahf addr -g
    
    # get tokens from miner
    ./kahf token -addr [address] -amt 10 -g
    ```

2. Service mode

    run
    
    ```bash
    ./kahf --serv
    ```
    
    support commands
    
    ```bash
    ./kahf --serv <-help|-c>
    ```
    examples
    
    ```bash
    # clear database
    ./kahf --serv -c
   
    ```
    
    usage
    
    ```bash
    # get helper
    curl  -H 'Content-Type: application/json'  http://127.0.0.1:11360/api/helper |jq  
 
    # get utxo
    curl  -H 'Content-Type: application/json'  http://127.0.0.1:11360/api/utxo?addr=[address] |jq 
 
    # get tx records
    curl  -H 'Content-Type: application/json'  http://127.0.0.1:11360/api/tx?addr=[address] |jq 
 
    # get status of the address
    curl  -H 'Content-Type: application/json'  http://127.0.0.1:11360/api/status?addr=[address] |jq 
 
    ```
    
    `Note` The 'Content-Type: application/json' is required
    
1. Console model

    run
    
    ```bash
    ./kahf -console
    ```
   
    examples
    
    ```bash
    # get helper
    [wallet-cli]: help
    
    # get utxo
    [wallet-cli]: out -u [address]
    ```

#### Configuration instructions

***Config File Path:*** main/config.toml

```toml
version="testnet"
[miner]
address="TmbCBKbZF8PeSdj5Chm22T4hZRMJY5D8zyz"
key="c39fb9103419af8be42385f3d6390b4c0c8f2cb17cf24dd43a059c4045d1a419"
[rpc]
host="https://127.0.0.1:1234"
auth="admin:123"
syncBlockInterval=10
syncMemoryInterval=5
syncPeerInfoInterval=30
syncNodeInfoInterval=3600
syncLockedBlockInterval = 60
syncInvalidBlockInterval=60
[api]
listenPort="11360"
[auth]
secretKey=""
jwt="off"
[mysql]
user=""
password=""
dbname=""
prefix=""
maxidleconns=5
maxopenconns=5
[log]
# <console|file>
mode="console" 	

# <debug|info|warn|fail|error>	
level="debug"

buffersize=1000

```

`[miner]` The configuration under this node is used to get tokens from the miner‘s address, which are only used in private test environments

`[rpc]` The configuration under this node is used to access main node of qitmeer with rpc

`[api]` The configuration under this node is used to set api interface

`[auth]` The configuration under this node is used to set JWT permission validation

`[mysql]` The configuration under this node is used to set mysql

`[log]` The configuration under this node is used to set log