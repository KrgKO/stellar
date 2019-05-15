# Stellar

Stellar blockchain is a one of popular blockchain which is public blockchain

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


## Getting start

1. Start stellar horizon using docker-compose more information: https://hub.docker.com/r/stellar/quickstart/
```
    docker-compose up -d
```