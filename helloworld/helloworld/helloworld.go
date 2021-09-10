package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
	log.Println("============ minha primeira aplicação em golang ============")

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Erro em setar DISCOVERY_AS_LOCALHOST como variável de ambiente: %v", err)
	}
	fmt.Println("Chegou aqui depois de discovery")

	wallet, err := gateway.NewFileSystemWallet("gateway")
	if err != nil {
		log.Fatalf("Falhou em criar wallet: %v", err)
	}
	fmt.Println("Chegou depois do gateway")


	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Falhou em popular a wallet: %v", err)
		}
	}

	fmt.Println("Pegou a wallet")

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
		log.Fatalf("Falhou em conectar com o gateway: %v", err)
	}
	defer gw.Close()

	fmt.Println("Antes de pegar o canal")

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Fatalf("Falhou em pegar a network: %v", err)
	}

	fmt.Println("Depois de pegar o canal")

	contract := network.GetContract("helloworld")

	log.Println("--> Transação de Submit: InitLedger, função cria o conjunto inicial de ativos no razão")
	result, err := contract.SubmitTransaction("InitLedger")
	if err != nil {

		log.Fatalf("Falhou em InitLedger SUBMIT (altera estado da ledger) %v", err)
	}
	log.Println(string(result))

	log.Println("--> Transação Evaluate: QueryAllOis, função que retorna todos os ativos na ledger")
	result, err = contract.EvaluateTransaction("QueryAllOis")
	if err != nil {
		log.Fatalf("Falhou a EVALUATE (consulta sem alterar estado da ledger) transação: %v", err)
	}
	log.Println(string(result))
	/*
		Saudacao:  saudacao,
		Despedida: despedida,
		Oidenovo:  oidenovo,
		Pessoa:	pessoa
	*/

	log.Println("--> Transação de Submit: CreateOi, cria ativos com ID(OINUMEROQUALQUER), saudação, despedida, oidenovo, e pessoa")
	result, err = contract.SubmitTransaction("CreateOi", "OI11", "Cheguei otário", "Tô indo fdp", "Que cu", "MarianaArrombada")
	if err != nil {
		log.Fatalf("Falhou a SUBMIT (altera estado da ledger) transação: %v", err)
	}
	log.Println(string(result))

	log.Println("--> Transação Evaluate: QueryOi, função retorna um ativo com OIID")
	result, err = contract.EvaluateTransaction("QueryOi", "OI6")
	if err != nil {
		log.Fatalf("Falhou em Transação Evaluate: %v\n", err)
	}
	log.Println(string(result))

	log.Println("--> Transação Evaluate: ExisteOi, função que retorna um boleano se achou o ativo na ledger")
	result, err = contract.EvaluateTransaction("ExisteOi", "OI1")
	if err != nil {
		log.Fatalf("Falhou em ExisteOi Transação Evaluate: %v\n", err)
	}
	log.Println(string(result))

	log.Println("--> Transação de Submit: ChangeOiPessoa OI1, transfere para um novo dono Val Bandeira")
	_, err = contract.SubmitTransaction("ChangeOiPessoa", "OI1", "Val Bandeira")
	if err != nil {
		log.Fatalf("Falhou em ChangeOiPessoa Transação de Submit: %v", err)
	}

	log.Println("--> Transação Evaluate: QueryOi, function returns 'OI1' attributes(não sei pra que porra)")
	result, err = contract.EvaluateTransaction("QueryOi", "OI1")
	if err != nil {
		log.Fatalf("Falhou em QueryOi Transação Evaluate: %v", err)
	}
	log.Println(string(result))
	log.Println("============ fim da minha primeira aplicação em golang ============")
}

func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
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

	certPath := filepath.Join(credPath, "signcerts", "User1@org1.example.com-cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}
	fmt.Println("pegou as credenciais em de cert.pem")
	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	fmt.Println("pegou a chave")
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	fmt.Println("pegou arquivo da chave")
	fmt.Println(string(files[0].Name()))
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}

