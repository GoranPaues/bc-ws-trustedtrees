/*
Copyright 2018 brondera.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main


import (
	"fmt"
	"strconv"
	"math/rand"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
// TODO: Rename
type SimpleChaincode struct {
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data, so be careful to avoid a scenario where you
// inadvertently clobber your ledger's data!
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response  {
        fmt.Println("########### oow18_example Init ###########")
	_, args := stub.GetFunctionAndParameters()

	// Only create init entities if the user intends to do it.
	if len(args) > 0 {
		fmt.Println("User put in Init params")

		var A, B string    // Entities
		var Aval, Bval int // Asset holdings
		var err error
	
		if len(args) != 4 {
			return shim.Error("Incorrect number of Init arguments. Expecting 4")
		}
	
		// Initialize the chaincode
		A = args[0]
		Aval, err = strconv.Atoi(args[1])
		if err != nil {
			return shim.Error("Expecting integer value for asset holding")
		}
		B = args[2]
		Bval, err = strconv.Atoi(args[3])
		if err != nil {
			return shim.Error("Expecting integer value for asset holding")
		}
		fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)
	
		// Write the state to the ledger
		err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
		if err != nil {
			return shim.Error(err.Error())
		}
	
		err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
		if err != nil {
			return shim.Error(err.Error())
		}

	}

	return shim.Success(nil)
}

// Query ledger
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
		return shim.Error("Unknown supported call")
}

// Invoke is called per transaction on the chaincode. 
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
        fmt.Println("########### oow18_example Invoke ###########")
	function, args := stub.GetFunctionAndParameters()
	
	if function != "invoke" {
                return shim.Error("Unknown function call")
	}

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting at least 2")
	}

	if args[0] == "delete" {
		// Deletes an entity from its state
		// TODO: remove/modify
		return t.delete(stub, args)
	}

	if args[0] == "query" {
		// queries an entity state
		return t.query(stub, args)
	}
	if args[0] == "move" {
		// Moves an asset between entities
		// TODO: remove
		return t.move(stub, args)
	}
	if args[0] == "create" {
		// creates an entity with asset
		// TODO: remove
		return t.create(stub, args)
	}
	if args[0] == "createTreeOrder" {
		// Create a tree order from username and money
		return t.createTreeOrder(stub, args)
	}
	return shim.Error("Unknown action, check the first argument, must be one of 'delete', 'query', create, or 'move'")
}

// Create a tree order with username, ID, and append ORDER to username
func (t *SimpleChaincode) createTreeOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var owner string // Name of the "account holder" from input
	var cash int // Amount of money in the tree order from input
	var tree int // Number of trees to be planted
	var orderid string // Generate something at random
	var change string
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3, function followed by 1 name and 1 value")
	}

	owner = args[1]
	cash, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}

	// Create orderID from name + random int n, 0 <= n < 1000
	// TODO: verify ID does not exist
	// TODO: Maybe using ownername in this context is not that clever after all.
	orderid = owner + strconv.Itoa(rand.Intn(1000))

	// Convert cash to tree order
	// "Send" uneven/exceeded amount in return
	price := 10
	tree = cash / price
	change = strconv.Itoa(cash % price)
	fmt.Printf("Uneven amount! " + change + " will be sent in return\n")

	// Put orderid and no of trees on the ledger
	// Append string "ORDER_" to key string for trackability
	fmt.Printf("Attempt to put state: orderid = %d, tree = %d\n", orderid, strconv.Itoa(tree))
	err = stub.PutState("ORDER_" + orderid, []byte(strconv.Itoa(tree)))
	if err != nil {
		return shim.Error(err.Error())
	}
	
	// Put owner and orderID on the ledger
	// TODO: Add a wait thing to get tree order confirmation first
	// Append string "OWNER_" to key string for trackability
	fmt.Printf("Attempt to put state: owner = %d, orderid = %d\n", owner, orderid)
	err = stub.PutState("OWNER_" + owner, []byte(orderid))
	if err != nil {
		return shim.Error(err.Error())
	}

	//TODO: Return something informational to the rest response (log only, atm)
	jsonResp := "{\"Your ORDERID: \":\" " + orderid + "\",\"No of trees to be planted\":\" " + strconv.Itoa(tree) + "\"}"
	fmt.Printf("CreateTreeOrder Response:%s\n", jsonResp)
	return shim.Success(nil)
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var qt string // QueryType
	var id string // name or orderid
	var Avalbytes []byte
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting querytype: 'name' or 'id' and string to query")
	}

	qt = args[1]
	id = args[2]

	if qt == "name" {
		// name-strings are of type OWNER_name
		id = "OWNER_" + id

		// Get the state from the ledger
		Avalbytes, err = stub.GetState(id)
	}

	if qt == "id" {
		// name-strings are of type ORDER_id
		id = "ORDER_" + id

		// Get the state from the ledger
		Avalbytes, err = stub.GetState(id)
	}

	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + id + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + id + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Param A\":\"" + id + "\",\"Param B\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

// creates an entity with asset
func (t *SimpleChaincode) create(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// must be an invoke
	var X string    // Entity
	var Xval int // Asset holding
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3, function followed by 1 name and 1 value")
	}

	X = args[1]
	Xval, err = strconv.Atoi(args[2])

	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	fmt.Printf("X = %d, Xval = %d\n", X, Xval)

	// Write the state to the ledger
	err = stub.PutState(X, []byte(strconv.Itoa(Xval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) move(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// must be an invoke
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4, function followed by 2 names and 1 value")
	}

	A = args[1]
	B = args[2]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	Aval = Aval - X
	Bval = Bval + X
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

        return shim.Success(nil);
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[1]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// main function starts up the chaincode in the container during instantiate
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
