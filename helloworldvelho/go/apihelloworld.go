package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)
/*
Chaincode
ype Oi struct {
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

	

	return ctx.GetStub().PutState(oiNumber, oiAsBytes)
}

// QueryOi returns the Oi stored in the world state with given id
func (s *SmartContract) QueryOi(ctx contractapi.TransactionContextInterface, oiNumber string) (*Oi, error) {
	

	return oi, nil
}

// QueryAllOis returns all Ois found in world state
func (s *SmartContract) QueryAllOis(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	
	return results, nil
}

// ChangeOiPessoa atualiza o campo Pessoa da Oi com id fornecido no estado mundial
func (s *SmartContract) ChangeOiPessoa(ctx contractapi.TransactionContextInterface, oiNumber string, newPessoa string) error {
	
	return ctx.GetStub().PutState(oiNumber, oiAsBytes)
}

*/
// "Person type" (tipo um objeto)
type Person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

var people []Person

// GetPeople mostra todos os contatos da variável people
func GetPeople(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}

// GetPerson mostra apenas um contato
func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}

// CreatePerson cria um novo contato
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

// DeletePerson deleta um contato
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
}

// função principal para executar a api
func main() {
	router := mux.NewRouter()
	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	router.HandleFunc("/contato", GetPeople).Methods("GET")
	router.HandleFunc("/contato/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/contato/{id}", CreatePerson).Methods("POST")
	router.HandleFunc("/contato/{id}", DeletePerson).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}