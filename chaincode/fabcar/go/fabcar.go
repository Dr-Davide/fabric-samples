/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a gbtr
type SmartContract struct {
	contractapi.Contract
}

// GBtr describes basic details of what makes up a GB-FLex transaction

type GBtr struct {
	OBJ_UIT   string  `json:"obj_uit"`
	OBJ_QTY   float32 `json:"obj_qty"`
	TSENDER   string  `json:"sender"`
	TICK      int     `json:"tick"`
	TRECEIVER string  `json:"treceiver"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *GBtr
}

// InitLedger adds a base set of GBtr to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	gbtrs := []GBtr{

		GBtr{
			OBJ_UIT:   "init",
			OBJ_QTY:   0,
			TSENDER:   "init",
			TICK:      0,
			TRECEIVER: "init",
		},
	}

	for i, gbtr := range gbtrs {
		gbtrAsBytes, _ := json.Marshal(gbtr)
		err := ctx.GetStub().PutState("GBTR"+strconv.Itoa(i), gbtrAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateGBtr adds a new gbtr to the world state with given details
func (s *SmartContract) CreateGBtr(ctx contractapi.TransactionContextInterface, gbtrNumber string, obj_uit string, obj_qty float32, tsender string, tick int, treceiver string) error {
	gbtr := GBtr{
		OBJ_UIT:   obj_uit,
		OBJ_QTY:   obj_qty,
		TSENDER:   tsender,
		TICK:      tick,
		TRECEIVER: treceiver,
	}

	gbtrAsBytes, _ := json.Marshal(gbtr)

	return ctx.GetStub().PutState(gbtrNumber, gbtrAsBytes)
}

// QueryGBtr returns the gbtr stored in the world state with given id
func (s *SmartContract) QueryGBtr(ctx contractapi.TransactionContextInterface, gbtrNumber string) (*GBtr, error) {
	gbtrAsBytes, err := ctx.GetStub().GetState(gbtrNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if gbtrAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", gbtrNumber)
	}

	gbtr := new(GBtr)
	_ = json.Unmarshal(gbtrAsBytes, gbtr)

	return gbtr, nil
}

// QueryAllGBtrs returns all gbtr found in world state
func (s *SmartContract) QueryAllGBtrs(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		gbtr := new(GBtr)
		_ = json.Unmarshal(queryResponse.Value, gbtr)

		queryResult := QueryResult{Key: queryResponse.Key, Record: gbtr}
		results = append(results, queryResult)
	}

	return results, nil
}

/*
// ChangeCarOwner updates the owner field of car with given id in world state
func (s *SmartContract) ChangeCarOwner(ctx contractapi.TransactionContextInterface, carNumber string, newOwner string) error {
	car, err := s.QueryCar(ctx, carNumber)

	if err != nil {
		return err
	}

	car.Owner = newOwner

	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(carNumber, carAsBytes)
}
*/
func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
