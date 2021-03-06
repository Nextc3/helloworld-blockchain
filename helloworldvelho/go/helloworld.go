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

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		fmt.Printf("Falha em criar wallet: %s\n", err)
		os.Exit(1)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			fmt.Printf("Falha em preencher a wallet com contéudo: %s\n", err)
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
		fmt.Printf("Falha em conectar gateway: %s\n", err)
		os.Exit(1)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		fmt.Printf("Falhar em pegar a network: %s\n", err)
		os.Exit(1)
	}
	

	contract := network.GetContract("helloworld")


	//mostrar o conteúdo da Ledger e pegar todos os OI's
	result, err := contract.EvaluateTransaction("queryAllOis")
	if err != nil {
		fmt.Printf("Falha ao avaliar a transação: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))
	/*
	Oi := Oi{
		Saudacao:  saudacao,
		Despedida: despedida,
		Oidenovo:  oidenovo,
		Pessoa:    pessoa,
	}

	*/

	result, err = contract.SubmitTransaction("createOi", "Cheguei otário", "Tô indo fdp", "Que cu", "MarianaArrombada")
	if err != nil {
		fmt.Printf("Falhou a SUBMIT (altera estado da ledger) transação: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	result, err = contract.EvaluateTransaction("queryOi", "OI6")
	if err != nil {
		fmt.Printf("Falhou a EVALUATE (consulta sem alterar estado da ledger) transação: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	_, err = contract.SubmitTransaction("changeOiPessoa", "OI6", "Val Bandeira")
	if err != nil {
		fmt.Printf("Falhou a SUBMIT (altera estado da ledger) transação: %s\n", err)
		os.Exit(1)
	}

	result, err = contract.EvaluateTransaction("queryOi", "OI6")
	if err != nil {
		fmt.Printf("Falhou a EVALUATE (consulta sem alterar estado da ledger) transação: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))
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
