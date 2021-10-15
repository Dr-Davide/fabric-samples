/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"encoding/json"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		os.Exit(1)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			fmt.Printf("Failed to populate wallet contents: %s\n", err)
			os.Exit(1)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		os.Exit(1)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		fmt.Printf("Failed to get network: %s\n", err)
		os.Exit(1)
	}

	contract := network.GetContract("fabcar")

	fmt.Println("\nAll Transactions on the ledger:")
	result, err := contract.EvaluateTransaction("queryAllGBtrs")
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	//temp-json parsing

	type GBtr struct {
		OBJ_UIT   string  `json:"OBJ_UIT"`
		OBJ_QTY   float32 `json:"OBJ_QTY"`
		TSENDER   string  `json:"TSENDER"`
		TICK      int     `json:"TICK"`
		OBJECT    string  `json:"OBJECT"`
		TRECEIVER string  `json:"TRECEIVER"`
	}

	type GBtrs struct {
		Items []GBtr `json:"items"`
	}
	// we initialize our Users array
	var gbtr GBtrs

	// read file

	fmt.Println("Loading JSON")
	data, err := ioutil.ReadFile("gbTransactions/gbTransactions.json")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("JSON Loaded")

	// unmarshall it
	err = json.Unmarshal(data, &gbtr)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println(string(data))
	fmt.Printf("N. of elements: %d\n", len(gbtr.Items))

	for i := 0; i < len(gbtr.Items); i++ {
		fmt.Println("OBJ_UIT: " + gbtr.Items[i].OBJ_UIT)
		fmt.Printf("OBJ_QTY: %s\n", gbtr.Items[i].OBJ_QTY)
		fmt.Println("TSENDER: " + gbtr.Items[i].TSENDER)
		fmt.Println("TICK: %s\n", gbtr.Items[i].TICK)
		fmt.Println("OBJECT: " + gbtr.Items[i].OBJECT)
		fmt.Println("TRECEIVER: " + gbtr.Items[i].TRECEIVER)
	}

	// ----------------------

	fmt.Println("END OF TEST")
}

func populateWallet(wallet *gateway.Wallet) error {
	credPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	err = wallet.Put("appUser", identity)
	if err != nil {
		return err
	}
	return nil
}
