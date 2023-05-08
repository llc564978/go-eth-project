# Ethereum Blockchain Service
This project is an Ethereum blockchain service implemented in Golang, consisting of an API service and a blockchain data indexer service. The API service provides 3 basic API endpoints for querying blockchain block and transaction data. The blockchain data indexer service is responsible for fetching the blockchain data and storing it in a PostgreSQL database.

## API Service
The API service provides the following endpoints:

* [GET] /blocks?limit=n Returns the latest n blocks (without all transaction hashes)
* [GET] /blocks/:id Returns a single block by block id (including all transaction hashes)
* [GET] /transaction/:txHash Returns transaction data with event logs

## Blockchain Data Indexer Service
The blockchain data indexer service fetches blockchain data using the web3 API through RPC and stores the data in a PostgreSQL database. It scans the blockchain data in parallel, starting from block n, and continues scanning until the latest block is reached. It continuously scans new and old blocks on the chain.

## Getting Started
### Prerequisites

* Golang
* PostgreSQL

### Setting Up
1. Clone the repository:
```
git clone https://github.com/llc564978/go-eth-project.git
```

2. Change to the project directory:
```
cd go-eth-project
```

3. Install dependencies:

```
go mod download
```

4. Set up the PostgreSQL database and update the database connection parameters in the models/db.go file.

5. Replace the rpcURL variable in the services/indexer.go file with your Ethereum RPC URL.

## Running the Application
### Build the project:

```
go build -o main
```

### Run the application:

```
go run .
```

### Run the Docker-compose:

```
docker-compose up
```

The API service will be running on http://localhost:8080.