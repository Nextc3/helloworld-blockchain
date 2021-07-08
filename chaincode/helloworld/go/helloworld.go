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

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

type Oi struct {
	Saudacao   string `json:"saudacao"`
	Despedida  string `json:"despedida"`
	Oidenovo string `json:"oidenovo"`
	Pessoa string `json:"pessoa"`
	
}
// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Oi
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	ois := []Oi{
		Oi{Saudacao: "Bom dia", Despedida: "Tchau", Oidenovo: "De novo?", Pessoa: "Marcola"},
		Oi{Saudacao: "Boa noite", Despedida: "Tô indo", Oidenovo: "já to enjoado", Pessoa: "BeiraMar"},
		Oi{Saudacao: "Boa tarde", Despedida: "Já fui", Oidenovo: "ai", Pessoa: "Marcinho VP"},
		Oi{Saudacao: "Olá", Despedida: "Xau", Oidenovo: "naaaum", Pessoa: "Escadinha"},
		Oi{Saudacao: "Satisfação", Despedida: "Fui nega", Oidenovo: "aff", Pessoa: "´Zé Pequeno"},
		Oi{Saudacao: "Saudações", Despedida: "Vou lá", Oidenovo: "oi", Pessoa: "Uê"},
		Oi{Saudacao: "E aí", Despedida: "...", Oidenovo: "vou lá", Pessoa: "Ravengar"},
		Oi{Saudacao: "oi sumida rs", Despedida: "Adeus", Oidenovo: "glub glub", Pessoa: "Bibi Perigosa"},
		Oi{Saudacao: "Tudo Nosso", Despedida: "_|_", Oidenovo: "Nada deles", Pessoa: "Lox Antrax"},
		Oi{Saudacao: "Plata o plomo?", Despedida: "Las mentiras", Oidenovo: "La guerra", Pessoa: "Pablo Escobar"},
	}

	for i, oi := range ois {
		oiAsBytes, _ := json.Marshal(oi)
		err := ctx.GetStub().PutState("OI"+strconv.Itoa(i), oiAsBytes)

		if err != nil {
			return fmt.Errorf("Falha ao colocar no estado mundial. %s", err.Error())
		}
	}

	return nil
}

// CreateCar adds a new car to the world state with given details
func (s *SmartContract) CreateCar(ctx contractapi.TransactionContextInterface, carNumber string, make string, model string, colour string, owner string) error {
	car := Car{
		Make:   make,
		Model:  model,
		Colour: colour,
		Owner:  owner,
	}

	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(carNumber, carAsBytes)
}

// QueryCar returns the car stored in the world state with given id
func (s *SmartContract) QueryCar(ctx contractapi.TransactionContextInterface, carNumber string) (*Car, error) {
	carAsBytes, err := ctx.GetStub().GetState(carNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if carAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", carNumber)
	}

	car := new(Car)
	_ = json.Unmarshal(carAsBytes, car)

	return car, nil
}

// QueryAllCars returns all cars found in world state
func (s *SmartContract) QueryAllCars(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
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

		car := new(Car)
		_ = json.Unmarshal(queryResponse.Value, car)

		queryResult := QueryResult{Key: queryResponse.Key, Record: car}
		results = append(results, queryResult)
	}

	return results, nil
}

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