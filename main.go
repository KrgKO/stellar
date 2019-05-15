package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/stellar/go/keypair"
)

// GenerateAddress will return address that created from stellar SDK
func GenerateAddress() (string, error) {
	pair, err := keypair.Random()
	if err != nil {
		return "", err
	}

	secret := pair.Seed()
	address := pair.Address()

	log.Println("Secret: ", secret)      // secret seed
	log.Println("Public Key: ", address) // public key

	// TODO: write to file

	return address, nil
}

// CreateAccount will lend XRP from bot
func CreateAccount(address string) error {
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + address)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return nil
}

func main() {
	address, err := GenerateAddress()
	if err != nil {
		log.Fatal(err)
	}

	err = CreateAccount(address)
	if err != nil {
		log.Fatal(err)
	}
}
