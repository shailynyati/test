package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

//Chaincode interface
type SimpleChaincode struct {
}

type UserRegistrationDetails struct {
	Ffid        string `json:"ffid"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	DOB         string `json:"DOB"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	Country     string `json:"country"`
	City        string `json:"city"`
	Zip         string `json:"zip"`
	CreatedBy   string `json:"createdby"`
	Title       string `json:"title"`
	Gender      string `json:"gender"`
	TotalPoints string `json:"totalPoints"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting User registration: %s", err)
	}
}

// To register User
func (t *SimpleChaincode) RegisterUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) < 2 {
		fmt.Println("Invalid number of args")
		return nil, errors.New("Expected at least two arguments for User registration")
	}
	user := UserRegistrationDetails{
		Ffid:        args[0],
		Firstname:   args[1],
		Lastname:    args[2],
		DOB:         args[3],
		Email:       args[4],
		Address:     args[5],
		Country:     args[6],
		City:        args[7],
		Zip:         args[8],
		CreatedBy:   args[9],
		Title:       args[10],
		Gender:      args[11],
		TotalPoints: args[12]}

	UserRegistrationBytes, err := json.Marshal(user)
	err = stub.PutState(args[0], UserRegistrationBytes)

	if err != nil {
		fmt.Println("Could not save UserRegistration to ledger", err)
		return nil, err
	}

	fmt.Println("Successfully saved User Registration")
	return nil, nil
}

//args[0] = id [<number>]
//args[1] = operator [add, delete]
//args[2] = points [<number>]
func (t *SimpleChaincode) AddDeletePoints(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var totalPoints int
	var pointsToModifyInt int

	operator := args[1]
	pointsToModify := args[2]

	userAsbytes, _ := t.getUser(stub, args)

	user := UserRegistrationDetails{}
	err := json.Unmarshal(userAsbytes, &user)

	if err != nil {
		return nil, err
	}

	totalPoints, _ = strconv.Atoi(user.TotalPoints)
	pointsToModifyInt, _ = strconv.Atoi(pointsToModify)

	if operator == "Add" {
		totalPoints = totalPoints + pointsToModifyInt
	}

	if operator == "Delete" {
		if totalPoints < pointsToModifyInt {
			return nil, errors.New("Points not sufficient to spend")
		}
		totalPoints = totalPoints - pointsToModifyInt
	}

	user.TotalPoints = strconv.Itoa(totalPoints)
	UserRegistrationBytes, _ := json.Marshal(user)
	err = stub.PutState(args[0], UserRegistrationBytes)

	if err != nil {
		return nil, err
	}

	return nil, err
}

func (t *SimpleChaincode) getPoints(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) < 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	user, err := t.getUser(stub, args)

	if err != nil {
		return nil, err
	}

	u := UserRegistrationDetails{}
	jsonResp := json.Unmarshal(user, &u)
	points := []byte(u.TotalPoints)
	return points, jsonResp
}

// Init resets all the things
//func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
//	if len(args) != 1 {
//		return nil, errors.New("Incorrect number of arguments. Expecting 1")
//	}
//	err := stub.PutState("User", []byte(args[0]))
//	if err != nil {
//		return nil, err
//	}
//	return nil, nil
//}

// Invoke is your entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "RegisterUser" {
		return t.RegisterUser(stub, args)
	}
	if function == "RegisterUserDetails" {
		return t.RegisterUser(stub, args)
	}

	if function == "AddDeletePoints" {
		return t.AddDeletePoints(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)
	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	//	if function == "read" { //read a variable
	//		return t.read(stub, args)
	//	}
	if function == "getUser" {
		return t.getUser(stub, args)
	}
	if function == "getPoints" {
		return t.getPoints(stub, args)
	}
	if function == "GetUserDetails" {
		return t.GetUserDetails(stub, args)
	}
	if function == "GetUserCount" {
		return t.GetUserCount(stub, args)
	}

	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// Get User - query function to read key/value pair

func (t *SimpleChaincode) getUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error

	if len(args) < 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}
	ffId := args[0] //keys to read from chaincode
	fmt.Print(ffId + " this is is the key ")
	userAsbytes, err := stub.GetState(ffId)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + ffId + "\"}"
		return nil, errors.New(jsonResp)
	}
	return userAsbytes, nil
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Check if table already exists
	_, err := stub.GetTable("UserTable")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}

	// Create application Table
	err = stub.CreateTable("UserTable", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "ffid", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "firstName", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "lastName", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "dob", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "email", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "address", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "country", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "city", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "zip", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "createdBy", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "title", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "gender", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "points", Type: shim.ColumnDefinition_STRING, Key: false},
	})

	if err != nil {
		return nil, errors.New("Failed creating UserTable.")
	}
	stub.PutState("userId", []byte("userId"))
	return nil, nil
}

//To store user details in User Table
func (t *SimpleChaincode) RegisterUserDetails(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	ffId := args[0]
	firstName := args[1]
	lastName := args[2]
	DOB := args[3]
	emailId := args[4]
	address := args[5]
	country := args[6]
	city := args[7]
	zip := args[8]
	createdBy := args[9]
	title := args[10]
	gender := args[11]
	points := args[12]

	ok, err := stub.InsertRow("UserTable", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: ffId}},
			&shim.Column{Value: &shim.Column_String_{String_: firstName}},
			&shim.Column{Value: &shim.Column_String_{String_: lastName}},
			&shim.Column{Value: &shim.Column_String_{String_: DOB}},
			&shim.Column{Value: &shim.Column_String_{String_: emailId}},
			&shim.Column{Value: &shim.Column_String_{String_: address}},
			&shim.Column{Value: &shim.Column_String_{String_: country}},
			&shim.Column{Value: &shim.Column_String_{String_: city}},
			&shim.Column{Value: &shim.Column_String_{String_: zip}},
			&shim.Column{Value: &shim.Column_String_{String_: createdBy}},
			&shim.Column{Value: &shim.Column_String_{String_: title}},
			&shim.Column{Value: &shim.Column_String_{String_: gender}},
			&shim.Column{Value: &shim.Column_String_{String_: points}},
		}})

	if err != nil {
		return nil, fmt.Errorf("insertTableOne operation failed. %s", err)
		panic(err)

	}
	if !ok {
		return []byte("Row with given key" + args[0] + " already exists"), errors.New("insertTableOne operation failed. Row with given key already exists")
	}
	return nil, nil
}

func (t *SimpleChaincode) GetUserDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting applicationid to query")
	}

	ffId := args[0]

	// Get the row pertaining to this userId
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: ffId}}
	columns = append(columns, col1)

	row, err := stub.GetRow("UserTable", columns)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get the data for the application " + ffId + "\"}"
		return nil, errors.New(jsonResp)
	}

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		jsonResp := "{\"Error\":\"Failed to get the data for the application " + ffId + "\"}"
		return nil, errors.New(jsonResp)
	}

	res2E := UserRegistrationDetails{}

	res2E.Ffid = row.Columns[0].GetString_()
	res2E.Firstname = row.Columns[1].GetString_()
	res2E.Lastname = row.Columns[2].GetString_()
	res2E.DOB = row.Columns[3].GetString_()
	res2E.Email = row.Columns[4].GetString_()
	res2E.Address = row.Columns[5].GetString_()
	res2E.City = row.Columns[6].GetString_()
	res2E.Country = row.Columns[7].GetString_()
	res2E.Zip = row.Columns[8].GetString_()
	res2E.CreatedBy = row.Columns[9].GetString_()
	res2E.Title = row.Columns[10].GetString_()
	res2E.Gender = row.Columns[11].GetString_()
	res2E.TotalPoints = row.Columns[12].GetString_()

	mapB, _ := json.Marshal(res2E)
	fmt.Println(string(mapB))

	return mapB, nil
}

//for counting users
type CountApplication struct {
	Count int `json:"count"`
}

// To count number of users
func (t *SimpleChaincode) GetUserCount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var columns []shim.Column
	contractCounter := 0

	rows, err := stub.GetRows("UserTable", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}

	for row := range rows {
		if len(row.Columns) != 0 {
			contractCounter++
		}
	}

	res2E := CountApplication{}
	res2E.Count = contractCounter
	mapB, _ := json.Marshal(res2E)
	fmt.Println(string(mapB))

	return mapB, nil
}
