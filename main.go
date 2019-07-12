package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/stellar/go/build"
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

	_, err = f.WriteString(fmt.Sprintf("%s %s\r\n", secret, address))
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

// RetrieveSeed will retrieve secret by public key
func RetrieveSeed(sourceAddress string) (string, error) {
	var seed string

	f, err := os.Open("./accounts")
	if err != nil {
		return "", err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " ")
		if sourceAddress == s[1] {
			seed = s[0]
		}
	}

	return seed, nil
}

// GenerateTransaction will create stellar transaction and commit on testnet
func GenerateTransaction(source string, destination string, amount string) (string, error) {

	sourceSecret, err := RetrieveSeed(source)
	if err != nil {
		return "", err
	}

	// Build transaction
	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{sourceSecret},
		build.AutoSequence{horizon.DefaultTestNetClient},
		build.Payment(
			build.Destination{destination},
			build.NativeAmount{amount},
		),
	)
	if err != nil {
		return "", err
	}

	txSigned, err := tx.Sign(sourceSecret)
	if err != nil {
		return "", err
	}

	txSignedBase64, err := txSigned.Base64()
	if err != nil {
		return "", err
	}

	response, err := horizon.DefaultTestNetClient.SubmitTransaction(txSignedBase64)
	if err != nil {
		return "", err
	}

	log.Println("Response transaction: ", response)
	return response.Hash, nil
}

func main() {
	var source string
	var destination string
	var amount string
	var flag string

	fmt.Printf("Create new account (y/n): ")
	fmt.Scanf("%s", &flag)

	if flag == "y" {
		address, err := GenerateAddress()
		if err != nil {
			log.Fatal(err)
		}

		source = address

		// Can work only new account
		err = LendXLM(source)
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}

	fmt.Printf("Input source address: ")
	fmt.Scanf("%s", &source)

	if source == "" {
		log.Fatalln("Source is required !")
		os.Exit(1)
	}

	err := GetAccountDetails(source)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Input destination address: ")
	fmt.Scanf("%s", &destination)

	if destination == "" {
		log.Fatalln("Destination is required !")
		os.Exit(1)
	}

	err = GetAccountDetails(destination)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("===========================")
	fmt.Printf("Input amount to transfer: ")
	fmt.Scanf("%s", &amount)

	amountToInt, err := strconv.Atoi(amount)
	if err != nil {
		log.Fatalln(err)
	}

	if amountToInt > 0 {
		txHash, err := GenerateTransaction(source, destination, amount)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Transaction Hash:", txHash)
	}
}
