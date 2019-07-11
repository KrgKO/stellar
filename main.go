package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"

	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

func writeToFile(secret string, address string) error {
	filename := "./accounts"

	if _, err := os.Stat(filename); err != nil {
		log.Fatalf(filename, " is not exist")
		_, err = os.Create(filename)
		if err != nil {
			return err
		}
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = f.WriteString(fmt.Sprintf("%s|%s\r\n", secret, address))
	if err != nil {
		return err
	}

	defer f.Close()
	return nil
}

// GenerateAddress will return address that created from stellar SDK
func GenerateAddress() (string, error) {
	pair, err := keypair.Random()
	if err != nil {
		return "", err
	}

	secret := pair.Seed()
	address := pair.Address()

	log.Println("Secret:", secret)      // secret seed
	log.Println("Public Key:", address) // public key

	err = writeToFile(secret, address)
	if err != nil {
		return "", errors.Wrap(err, "Error while written file")
	}
	log.Println("The secret and public key has been written!")

	return address, nil
}

// LendXLM will lend XLM from bot
func LendXLM(address string) error {
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

// GetAccountDetails will get accout details
func GetAccountDetails(address string) error {
	account, err := horizon.DefaultTestNetClient.LoadAccount(address)
	if err != nil {
		return err
	}

	fmt.Println("Balance of account:", address)

	for _, balance := range account.Balances {
		log.Println(balance.Balance)
	}

	return nil
}

func main() {
	var userAddress string

	fmt.Printf("Input your address: ")
	fmt.Scanf("%s", &userAddress)

	if userAddress == "" {
		address, err := GenerateAddress()
		if err != nil {
			log.Fatal(err)
		}

		userAddress = address
	}

	err := LendXLM(userAddress)
	if err != nil {
		log.Fatal(err)
	}

	err = GetAccountDetails(userAddress)
	if err != nil {
		log.Fatal(err)
	}
}
