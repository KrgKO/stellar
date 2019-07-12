# Stellar

Stellar blockchain is a one of popular blockchain which is public blockchain

More easier: https://www.stellar.org/laboratory/
More stellargo document: https://godoc.org/github.com/stellar/go

## Components

1. Stellar network - contains many collection of stellar core
2. Stellar core - backbone of stellar network and validate, agree by using stellar consensus protocol
3. Stellar horizon - RESTful HTTP API server

- Stellar consenesus protocol - using **Federated Byzantine** which is adaptation of byzantine -> majority node must be agree it contains 2 main mechanism
    1. federated voting
    2. federated leader selection - select pseudorandom leaders
**ref:** http://www.scs.stanford.edu/~dm/blog/simplified-scp.html

## Install Golang package
To install golang package use command `go get github.com/stellar/go`

If got an error install other dependencies `go get ./...`

Or if want to use go version 1.11+ `go mod tidy`

## Getting start

1. Start stellar horizon using docker-compose more information: https://hub.docker.com/r/stellar/quickstart/
```
    docker-compose up -d
```
2. Run golang source code
```
    go run main.go
```

## Account

The stellar account contains public key (Gxxxxxxxxxxxxxx) and secret (Sxxxxxxxxxxxxx). They saved in the ledger. Also support multi-signature

- Keypair - contains address (public key), seed (secret)
- Public key use for source and destination definition
- Sign transaction use public key and secret

Reference: https://www.stellar.org/developers/guides/concepts/accounts.html

## Asset

Asset can be currencies

## Fees

Stellar network requires fee on transactions. in unit stroops (0.00001 lumen)

## Ledger

Contains list of accounts and balances. The first ledger is genesis ledger

## Transaction

Transaction are commands that modify the ledger state. Metaphor: Ledger is database, Transaction is SQL commands

Memo: can be used to store data but not much

Reference: https://www.stellar.org/developers/guides/concepts/transactions.html