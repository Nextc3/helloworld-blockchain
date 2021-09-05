/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	//"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

/*
type ServerConfig struct {
	CCID    string
	Address string
}

*/
// SmartContract  fornece funções para gerenciar um hello-world
type SmartContract struct {
	contractapi.Contract
}

//Struct que estabelece um ativo
type Oi struct {
	Saudacao  string `json:"saudacao"`
	Despedida string `json:"despedida"`
	Oidenovo  string `json:"oidenovo"`
	Pessoa    string `json:"pessoa"`
}

// Estrutura QueryResult usada para lidar com o resultado da consulta
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Oi
}

// InitLedger adds a base set of Oi's to the ledger
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
		//coloca um identificador nos OI's exemplo OI1, OI2, OI3
		err := ctx.GetStub().PutState("OI"+strconv.Itoa(i), oiAsBytes)

		if err != nil {
			return fmt.Errorf("Falha ao colocar no estado mundial.(Colocar na ledger) %s", err.Error())
		}
	}

	return nil
}

// Criar Oi adiciona uma nova Oi ao estado mundial com detalhes fornecidos
func (s *SmartContract) CreateOi(ctx contractapi.TransactionContextInterface, oiNumber string, saudacao string, despedida string, oidenovo string, pessoa string) error {
	Oi := Oi{
		Saudacao:  saudacao,
		Despedida: despedida,
		Oidenovo:  oidenovo,
		Pessoa:    pessoa,
	}

	oiAsBytes, _ := json.Marshal(Oi)

	return ctx.GetStub().PutState(oiNumber, oiAsBytes)
}

// QueryOi returns the Oi stored in the world state with given id
func (s *SmartContract) QueryOi(ctx contractapi.TransactionContextInterface, oiNumber string) (*Oi, error) {
	oiAsBytes, err := ctx.GetStub().GetState(oiNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if oiAsBytes == nil {
		return nil, fmt.Errorf("%s não existe", oiNumber)
	}

	oi := new(Oi)
	_ = json.Unmarshal(oiAsBytes, oi)

	return oi, nil
}

//Consulta se Oi existe 
func (s *SmartContract) ExisteOi(ctx contractapi.TransactionContextInterface, oiNumber string) (bool, error)  {
	oiAsBytes, err := ctx.GetStub().GetState(oiNumber)
	if err != nil {
		return false, fmt.Errorf("falhou em ler o estado do bagulho: %v", err)
	}

	return oiAsBytes != nil, nil
}

func (s *SmartContract) DeleteOi(ctx contractapi.TransactionContextInterface, oiNumber string) error {
	exists, err := s.ExisteOi(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("o ativo %s não ecsiste", id)
	}

	return ctx.GetStub().DelState(oiNumber)

}

// QueryAllOis returns all Ois found in world state
func (s *SmartContract) QueryAllOis(ctx contractapi.TransactionContextInterface) ([]*QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []*QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		oi := new(Oi)
		_ = json.Unmarshal(queryResponse.Value, oi)

		queryResult := QueryResult{Key: queryResponse.Key, Record: oi}
		results = append(results, &queryResult)
	}

	return results, nil
}

// ChangeOiPessoa atualiza o campo Pessoa da Oi com id fornecido no estado mundial
func (s *SmartContract) ChangeOiPessoa(ctx contractapi.TransactionContextInterface, oiNumber string, newPessoa string) error {
	oi, err := s.QueryOi(ctx, oiNumber)

	if err != nil {
		return err
	}

	oi.Pessoa = newPessoa

	oiAsBytes, _ := json.Marshal(oi)

	return ctx.GetStub().PutState(oiNumber, oiAsBytes)
}

func main() {
	// See chaincode.env.example
	/*
		config := ServerConfig{
			CCID:    os.Getenv("CHAINCODE_ID"),
			Address: os.Getenv("CHAINCODE_SERVER_ADDRESS"),
		} */

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Erro em criar helloworld chaincode: %s", err.Error())
		return
	}
	/*
		server := &shim.ChaincodeServer{
			CCID:    config.CCID,
			Address: config.Address,
			CC:      chaincode,
			TLSProps: shim.TLSProperties{
				Disabled: true,
			},
		}

		if err := server.Start(); err != nil {
			fmt.Printf("Erro em estartar helloworld chaincode: %s", err.Error())
		}

	*/

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Erro em criar helloworld chaincode: %s", err.Error())
	}
}
