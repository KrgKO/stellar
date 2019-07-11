package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/pkg/errors"

	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
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

	_, err = f.WriteString(fmt.Sprintf("%s %s %s\r\n", secret, address, network.TestNetworkPassphrase))
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

// GenerateTransaction will create stellar transaction and commit on testnet
func GenerateTransaction(source string, destination string, amount string) (string, error) {
	// Retrieve account information
	client := horizonclient.DefaultTestNetClient
	accountRequest := horizonclient.AccountRequest{AccountID: source}
	sourceRetrieved, err := client.AccountDetail(accountRequest)
	if err != nil {
		return "", err
	}

	// Create operation
	accountOperation := txnbuild.CreateAccount{
		Destination: destination,
		Amount:      amount,
	}

	// Generate transaction
	tx := txnbuild.Transaction{
		SourceAccount: &sourceRetrieved,
		Operations:    []txnbuild.Operation{&accountOperation},
		Timebounds:    txnbuild.NewTimeout(300),
		Network:       network.TestNetworkPassphrase,
	}

	var byte32 [32]byte

	// filename := "./accounts"
	// f, err := os.Open(filename)
	// if err != nil {
	// 	log.Fatalln((err))
	// }
	// defer f.Close()
	// scanner := bufio.NewScanner(f)
	// for scanner.Scan() {
	// 	s := strings.Split(scanner.Text(), " ")
	// 	if source == s[1] {
	// 		for k, v := range []byte(s[0]) {
	// 			byte32[k] = v
	// 		}
	// 		log.Println(s[0], s[1], []byte(s[0]), byte32)
	// 	}
	// }

	// TODO: Find way to retrieve keypair
	newKeypair, err := keypair.FromRawSeed(byte32)
	if err != nil {
		return "", err
	}

	log.Println(newKeypair)
	// Sign transaction
	txBase64, err := tx.BuildSignEncode(newKeypair)
	if err != nil {
		return "", err
	}

	fmt.Println("Transaction:", txBase64)

	response, err := client.SubmitTransactionXDR(txBase64)
	if err != nil {
		horizonError := err.(*horizonclient.Error)
		log.Printf("Response error: %#v\n", horizonError.Response)
		return "", horizonError
	}

	return response.TransactionSuccessToString(), nil
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
		result, err := GenerateTransaction(source, destination, amount)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Transaction response:", result)
	}
}
