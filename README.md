
# ETH Defender

## Prerequisites

1.  **Docker**: Ensure Docker is installed and running.
2.  **Docker-compose**: Ensure Docker Compose is installed.
3.  **Deployment of Smart Contract**: Deploy `assets/contract/VulnerableContract.sol` using Remix, Truffle, or a similar tool.
4.  **Configuration**: Set up Ethereum and email SMTP information in `config.yml`.

## Component Description

### Detector

-   **Function**: Exploit Detection and Notification Microservice.
-   **Checks**:
    1.  **Deposit**: Detects deposits exceeding a threshold (currently set to 0.0001 ETH for testing).
    2.  **Withdrawal**: Detects withdrawals exceeding three times within one minutes.

### Preventer

-   **Function**: Front-running Prevention Microservice.
-   **Checks**: Monitors pending transactions; if the transaction gas price is high, sends a pause transaction to the contract.

### Reporter

-   **Functions**:
    
    1.  **Transaction Recording**: Records every transaction to the database (for system performance monitoring).
    2.  **Paused Event Listening**: Listens for "Paused" events from the contract, tests the pausing mechanism, and records it in the database. (For checking effectiveness of of the pausing mechanism)
    3.  **Reporting and Analytics APIs**:
        -   Analyzes smart contract files with Slither, creates reports, and stores them in the database.
        -   Lists smart contract reports.
        -   Lists transaction records.
        -   Lists smart contract paused records.
    
    -   **API Documentation**: Check APIs at [http://localhost:8080/docs/index.html](http://localhost:8080/docs/index.html).

## How To Run

### Start Server

Execute the following command in the repository directory:

`make setup` 

### Testing

#### Tester
```
make go-build-test
```
#### Detector

-   **Deposit Test**:

`./bin/ethtester dt1`

    -   Send more than 0.0001 ETH to the smart contract and verify that an email notification is received. 
-   **Withdrawal Test**:

`./bin/ethtester dt2`

    -   Withdraw ETH from the smart contract and verify that an email notification is received.

#### Preventer

-   **Test Execution**:

`./bin/ethtester pt`

    -   The Preventer container listens for any pending transactions. Upon detection, it will check for suspicious conditions(the gas price is bigger than 1 Gwei for the testing) and send a pause transaction to the smart contract.
-   **Verification**:
    -   Check the paused call on the smart contract using Remix or similar tools. The value should be `true`.

#### Reporter

-   **Report Creation**:

You can use `assets/postman/ETH defender.postman_collection.json`

**Slither Report**:
        
        curl --location 'http://localhost:8080/api/v1/report/slither' \
        --form 'contract_file=@"/solidy/file/path"' \
        --form 'contract_name="contract_name"'
        
**List Reports**:
        
        curl --location 'http://localhost:8080/api/v1/report/slither'
        curl --location 'http://localhost:8080/api/v1/report/pause'
        curl --location 'http://localhost:8080/api/v1/report/transaction'

## Remove Docker Containers and Images

-   **Stop and Remove Containers**:
    
    `make down` 
    
-   **Remove Images**:
    
    `make clean-image` 
    
-   **Perform Both Actions**:
    
    `make clean-all` 

## Restart and Rebuild Docker Images and Containers

Execute:

`make rebuild`
