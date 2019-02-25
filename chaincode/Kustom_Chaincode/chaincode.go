package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type Product struct {
	Nombre string `json:"Nombre"`
	Precio int    `json:"Precio"`
}

type User struct {
	Nombre string `json:"Nombre"`
	Saldo  int    `json:"Saldo"`
}

//Init biatch
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

//Invoke invokes invocable functions
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	var result string
	var error error

	function, args := APIstub.GetFunctionAndParameters()

	if function == "initLedger" {
		result, error = s.initLedger(APIstub, args)
	} else if function == "querySaldo" {
		result, error = s.querySaldo(APIstub, args)
	} else if function == "compraProd" {
		result, error = s.compraProd(APIstub, args)
	} else {
		return shim.Error("Funcion no reconocida")
	}
	if error != nil {
		return shim.Error(error.Error())
	}
	return shim.Success([]byte(result))
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	prod := []Product{
		Product{Nombre: "Negrita", Precio: 300},
	}
	usr := []User{
		User{Nombre: "Pedro", Saldo: 1000},
	}

	proAsBytes, _ := json.Marshal(prod)
	usrAsBytes, _ := json.Marshal(usr)
	APIstub.PutState("PROD0", proAsBytes)
	APIstub.PutState("USR0", usrAsBytes)

	return "", nil
}

//Consulta saldo de usuario
func (s *SmartContract) querySaldo(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Se esperaba un argumentos")
	}
	usr := User{}
	usrAsBytes, _ := APIstub.GetState(args[0])
	error := json.Unmarshal(usrAsBytes, &usr)
	return "Saldo es: " + string(usr.Saldo), error
}

//Compra productos, recibe como argumento usuario y productos
func (s *SmartContract) compraProd(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Se esperaban dos argumentos, Usuario, Producto")
	}
	usr := User{}
	prod := Product{}
	usrAsBytes, _ := APIstub.GetState(args[0])
	prodAsBytes, _ := APIstub.GetState(args[1])
	error := json.Unmarshal(usrAsBytes, &usr)
	error = json.Unmarshal(prodAsBytes, &prod)
	usr.Saldo = usr.Saldo - prod.Precio
	usrAsBytes, _ = json.Marshal(usr)
	APIstub.PutState(args[0], usrAsBytes)

	return "Producto comprado, nuevo saldo es " + string(usr.Saldo), error
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
