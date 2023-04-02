package main

/**
 * Shows how to use the history
 **/

import (
	// For printing messages on console
	"fmt"

	// April 2020, Updated to Fabric 2.0 Shim
	"github.com/hyperledger/fabric-chaincode-go/shim"

	peer "github.com/hyperledger/fabric-protos-go/peer"


	// KV Interface
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"

	// JSON Encoding
	"encoding/json"
	// "strings"
	"strconv"
)

// VehicleChaincode Represents our chaincode object
type VehicleChaincode struct {
}

// PkrToken structure manages the state
type PkrToken struct {
	Symbol      string `json:"symbol"`
	TotalSupply uint64 `json:"totalSupply"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
}
type Wallet struct {
	CNIC    string `json:"cnic"`
	Name    string `json:"name"`
	Balance string `json:"balance"`
}

// OwnerPrefix is used for creating the key for balances
const OwnerPrefix = "CNIC."
const FbrCnic = "6666"

// Vehicle Represents our car asset
type Vehicle struct {
	DocType          string `json:"docType"`
	EngineNo         string `json:"engineNo"`
	ChassisNo        string `json:"chassisNo"`
	Year             string   `json:"year"`
	Type             string `json:"type"` //bike or car
	Make             string `json:"make"`
	Model            string `json:"model"`
	CompanyName      string `json:"companyName"`
	ManufacturerCnic string `json:"manufacturerCnic"`
	OwnerCNIC        string `json:"ownerCNIC"`
	Sold             string `json:"soldPrice"`
	LaunchingPrice   string `json:"launchingPrice"`
	TransferDate     string `json:"transferDate"`
}

// DocType Represents the object type
const DocType = "VehicleAsset"

// vehicle docType
const objectType = "chassisNo~engineNo~companyName"

// Init Implements the Init method
// Init Implements the Init method
// Receives 4 parameters =  [0] Symbol [1] TotalSupply   [2] Description  [3] Owner
func (history *VehicleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {

	// Simply print a message
	fmt.Println("Init executed")
	_, args := stub.GetFunctionAndParameters()

	// Check if we received the right number of arguments
	if len(args) < 4 {
		return shim.Error("Failed - incorrect number of parameters!! ")
	}
	symbol := string(args[0])
	// Get total supply & check if it is > 0
	//second argument specify the base of conversion
	//3rd argument specify the size in bits

	totalSupply, err := strconv.ParseUint(string(args[1]), 10, 64)

	if err != nil || totalSupply == 0 {
		return shim.Error("Total Supply MUST be a number > 0 or there is error!! ")
	}

	// Creator name cannot be zero length
	if len(args[3]) == 0 {
		return errorResponse("Creator identity cannot be 0 length!!!", 3) //passing three to print also error code helps in debugging
	}
	creator := string(args[3])

	// Create an instance of the token struct
	var Pkr = PkrToken{Symbol: symbol, TotalSupply: totalSupply, Description: string(args[2]), Creator: creator}

	// Convert to JSON and store token in the state
	jsonPkr, _ := json.Marshal(Pkr)
	//jsonPkr is first converted to json and then to byte array
	stub.PutState("Pkrtoken", []byte(jsonPkr))

	// Maintain the balances in the state db
	// In the begining all tokens are owned by the creator of the token
	key := OwnerPrefix + "3520299610969"
	fmt.Println("Key=", key)
	//key setting=CNIC.35202-99610969
	var wallet = Wallet{CNIC: "3520299610969", Name: "M.Adil", Balance: string(args[1])}
	jsonwallet, _ := json.Marshal(wallet)
	err = stub.PutState(key, []byte(jsonwallet))
	if err != nil {
		return errorResponse(err.Error(), 4)
	}

	return shim.Success([]byte(jsonPkr))
}

// Invoke method
func (history *VehicleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	// Get the function name and parameters
	funcName, args := stub.GetFunctionAndParameters()

	fmt.Println("Invoke executed : ", funcName, " args=", args)

	switch {

	case funcName == "totalSupply":
		return totalSupply(stub)
	case funcName == "balanceOf":
		return balanceOf(stub, args)
	case funcName == "transfer":
		return transfer(stub, args)
	case funcName == "getStateRangeOnKey" :
			// Query with GetStateByPartialCompositeKey
			//getting vehicle current state
		return history.GetVehicleByPartialCompositeKey(stub, args)
		
	}

	//manufacturing use case
	if funcName == "Manufacture" {
		return history.Manufacture(stub, args)
	}else if funcName == "TransferOwnership" {
		// Invoke this function to transfer ownership of vehicle
		return history.TransferOwnership(stub, args)
	}else if funcName == "GetVehicleHistory" {
			// Query this function to get txn history for specific vehicle
		return history.GetVehicleHistory(stub, args)
	
	}else if funcName == "GetVehiclesByCNIC" {

		// Get all vehicle by year - just another example of query
		return history.GetVehiclesByCNIC(stub, args)
	}

	// if funcName == "GetVehicleByVin" {
	// 	// Returns the vehicle's current state
	// 	return history.GetVehicleByVin(stub, args)

	// } 


	// } else if funcName == "GetVehiclesOwners" {
	// 	// To be implemented in the exercise
	// 	// return history.GetVehiclesOwners(stub, args)
	// }

	// This is not good
	return shim.Error(("Bad Function Name = !!!"))
}
func addData(stub shim.ChaincodeStubInterface, vehicle Vehicle) peer.Response  {

	jsonvehicle, _ := json.Marshal(vehicle)
	balanceIndexKey, _ := stub.CreateCompositeKey(objectType, []string{vehicle.ChassisNo, vehicle.EngineNo, vehicle.CompanyName})
	stub.PutState(balanceIndexKey, jsonvehicle)
	return successResponse("successfully manufacture!!")
}
func (token *VehicleChaincode) GetVehicleByPartialCompositeKey(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	// Create the query key
	qryKey, errkey := stub.CreateCompositeKey(objectType, args)
	if errkey != nil {
		fmt.Printf("Error in creating key =" + errkey.Error())
		return shim.Error(errkey.Error())
	}
	
	fmt.Printf("Composite Key=%s\n", qryKey)

	var resultJSON = "["
	// Get the data
	dat, _ := stub.GetState(qryKey)


	// Set the data in result string
	resultJSON += string(dat)
	resultJSON += "]"

	return shim.Success([]byte(resultJSON)) 	
	// Print statements used in dev mode
	// fmt.Printf("==== Exec qry with:  ")
	// fmt.Println(args)
	// // Gets the state by partial query key
	// QryIterator, err := stub.GetStateByPartialCompositeKey(objectType, args)
	// if err != nil {
	// 	fmt.Printf("Error in getting by range=" + err.Error())
	// 	return shim.Error(err.Error())
	// }
	// var resultJSON = "["
	// counter := 0
	// // Iterate to read the keys returned
	// for QryIterator.HasNext() {
	// 	// Hold pointer to the query result
	// 	var resultKV *queryresult.KV
	// 	var err error

	// 	// Get the next element
	// 	resultKV, err = QryIterator.Next()
	// 	if err != nil {
	// 		fmt.Println("Err=" + err.Error())
	// 		return shim.Error(err.Error())
	// 	}

	// 	// Split the composite key and send it as part of the result set
	// 	key, arr, _ := stub.SplitCompositeKey(resultKV.GetKey())
	// 	fmt.Println(key)
	// 	resultJSON += " [" + strings.Join(arr, "//") + "] "
	// 	counter++

	// }
	// // Closing
	// QryIterator.Close()

	// resultJSON += "]"
	// resultJSON = "Counter=" + strconv.Itoa(counter) + "  " + resultJSON
	// fmt.Println("Done.")
	// return shim.Success([]byte(resultJSON))
}
//manufacture vehicle
//args order,:manufacturerCnic,ENgine No,Chassis no,company name,year ,type,model,Ownercnic,sold,launchingPrice,transferDate

func (history *VehicleChaincode) Manufacture(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// Check the number of args
	if len(args) < 11 {
		return shim.Error("Cannot Procede manufacturing incomplete details !!!")
	}
	//check the wallet if manufacturer has money to pay 16 % tax or not of launching price
	OwnerID := args[0]
	//contains cnic
	bytes, err := stub.GetState(OwnerPrefix + OwnerID)
	if err != nil {
		return errorResponse(err.Error(), 7)
	}
	//convert the balance in json
	var manufacturerwallet Wallet
	_ = json.Unmarshal(bytes, &manufacturerwallet)

	// balance:= strconv.ParseInt(manufacturerwallet.Balance, 10, 64)
	launchingPrice,_:= strconv.ParseUint(args[9], 10, 64)

	TaxAmount:= strconv.FormatUint(uint64(float64(launchingPrice)*0.16), 10)
	argg:=[]string{args[0], FbrCnic, TaxAmount, "Fbr"}
	transferResult := transfer(stub, argg)
	//if there is error in transfering money stop manufacture

	if transferResult.Status != 200 {
		return transferResult
	}

	//else tax has been paid
	Newvehicle := Vehicle{DocType: DocType, ManufacturerCnic: args[0], EngineNo: args[1], ChassisNo: args[2], CompanyName: args[3], Year: args[4], Type: args[5], Model: args[6], OwnerCNIC: args[7], Sold: args[8], LaunchingPrice: args[9], TransferDate: args[10]}
	return addData(stub, Newvehicle)
	// skey:=[]string{Newvehicle.ChassisNo, Newvehicle.EngineNo, Newvehicle.CompanyName}
	// balanceIndexKey, _ := stub.CreateCompositeKey(objectType, skey)
	// return successResponse("Vehicle Manufactured successfully!!!!"+balanceIndexKey)

}

// GetVehicleHistory gets the history of the vehicle by chassis no engine and company name

// args[0] = ChassisNo args[1 ]=Engine No args[2]=company name
func (history *VehicleChaincode) GetVehicleHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	// Check the number of args
	if len(args) < 3 {
		return shim.Error("MUST provide engine no chassis no and company name !!!")
	}
	qryKey, errkey := stub.CreateCompositeKey(objectType, args)
	if errkey != nil {
		fmt.Printf("Error in creating key =" + errkey.Error())
		return shim.Error(errkey.Error())
	}
	
	fmt.Printf("Composite Key=%s\n", qryKey)
	// Get the history for the key i.e., VIN#
	historyQueryIterator, err := stub.GetHistoryForKey(qryKey)

	// In case of error - return error
	if err != nil {
		return shim.Error("Error in fetching history !!!" + err.Error())
	}

	// Local variable to hold the history record
	var resultModification *queryresult.KeyModification
	counter := 0
	resultJSON := "["

	// Start a loop with check for more rows
	for historyQueryIterator.HasNext() {

		// Get the next record
		resultModification, err = historyQueryIterator.Next()

		if err != nil {
			return shim.Error("Error in reading history record!!!" + err.Error())
		}

		// Append the data to local variable
		data := "{\"txn\":" + resultModification.GetTxId()
		data += " , \"value\": " + string(resultModification.GetValue()) + "}  "
		if counter > 0 {
			data = ", " + data
		}
		resultJSON += data

		counter++
	}

	// Close the iterator
	historyQueryIterator.Close()

	// finalize the return string
	resultJSON += "]"
	resultJSON = "{ \"counter\": " + strconv.Itoa(counter) + ", \"txns\":" + resultJSON + "}"

	// return success
	return shim.Success([]byte(resultJSON))
}


// ***working hrere incomplete

// TransferOwnership gets the asset information
// Transfer the ownership of the vehicle from owner1 to owner2
// args[0]=vin   args[1]=current owner  args[2]=new owner args[3]=transfer date
// args[1] used for validation
func (history *VehicleChaincode) TransferOwnership(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 6 {
		return shim.Error("Must provide from ,to :owners cnic ,vehicle chassino engine no and company name,date")
	}



	// Create the query key
	qryKey, errkey := stub.CreateCompositeKey(objectType,[]string{args[2],args[3],args[4]})
	if errkey != nil {
		fmt.Printf("Error in creating key =" + errkey.Error())
		return shim.Error(errkey.Error())
	}
	
	fmt.Printf("Composite Key=%s\n", qryKey)

	// Get the data
	dat, err := stub.GetState(qryKey)

	if err!=nil{
		return errorResponse(err.Error(),888)
	}

	var vehicle Vehicle
	_=json.Unmarshal(dat,&vehicle)
	var a =args[0]
	if vehicle.OwnerCNIC!=a{
		fmt.Printf("You are not the owner")
		return errorResponse(vehicle.OwnerCNIC+" "+ a+"You are not the owner",889)

	}
	// #if he is the owner just change the vehilce cnic and update it 
	var newOwner = OwnerPrefix+args[1]
	vehicle.OwnerCNIC=newOwner


	vehicle.TransferDate = args[5]
	jsonVehicle, _ := json.Marshal(vehicle)

	stub.PutState(qryKey, jsonVehicle)
	// Set the data in result string

	return shim.Success([]byte("Vehicle Record Updated!!! " + string(jsonVehicle)))

}

// // GetVehicleByVin gets the asset information
// func (history *VehicleChaincode) GetVehicleByVin(stub shim.ChaincodeStubInterface, args []string) peer.Response {
// 	// Check on args
// 	if len(args) < 1 {
// 		return shim.Error("MUST provide Vin Number in args[0] !!!")
// 	}

// 	// Get the data
// 	vehicle, _ := stub.GetState(args[0])

// 	return shim.Success([]byte(vehicle))
// }

// GetVehiclesByCNIC gets all the vehiclesownercinc=cnic
// Another sample to show the use of Rich Queries
// To make this work you need to create an index :)
// Not using Pagination - so results restricted to a max of totalQueryLimit
func (history *VehicleChaincode) GetVehiclesByCNIC(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args[0]) < 1 {
		return shim.Error("Please provide a valid year !!!")
	}
	qry := `{
		"selector": {
		   "ownerCNIC": "`

	qry += args[0]
	qry += `"		  
		   }
	 }`
	// qry := `{
	// 	"selector": {
	// 	   "ownerCNIC": "7777"}
	//  }`
	// GetQueryResult
	QryIterator, err := stub.GetQueryResult(qry)
	if err != nil {
		return shim.Error("Error in executing rich query !!!! " + err.Error())
	}
	// hold the result json
	resultJSON := "["
	counter := 0
	for QryIterator.HasNext() {
		// Hold pointer to the query result
		var resultKV *queryresult.KV

		// Get the next element
		resultKV, _ = QryIterator.Next()

		value := string(resultKV.GetValue())
		if counter > 0 {
			resultJSON += ", "
		}
		resultJSON += value
		counter++
	}
	resultJSON += "]"

	return shim.Success([]byte(resultJSON))
}

// // SetupSampleData creates multiple instances of the ERC20history
// func (history *VehicleChaincode) SetupSampleData(stub shim.ChaincodeStubInterface) {

// 	// This the car data for testing
// 	AddData(stub, "100", "toyota", "corolla", 2011, "J Smith", "2015-12-20")
// 	AddData(stub, "200", "honda", "civic", 2012, "G Roger", "2016-01-15")
// 	AddData(stub, "300", "audi", "a5", 2015, "S Ripple", "2018-07-22")
// 	AddData(stub, "400", "bmw", "x5", 2013, "M Jane", "2019-02-19")
// 	AddData(stub, "500", "toyota", "camry", 2018, "J Hoover", "2019-01-15")

// 	fmt.Println("Initialized with the sample data!!")
// }

// AddData adds a car asset to the chaincode asset database
// Structure is created and initialized then it is marshalled to JSON for storage using PutState
// func AddData(stub shim.ChaincodeStubInterface, vin string, make, model string, year uint, owner, transfer string) {
// 	vehicle := Vehicle{DocType: DocType, VIN: vin, Year: year, Make: make, Model: model, Owner: owner, Transfer: transfer}
// 	jsonVehicle, _ := json.Marshal(vehicle)
// 	// Key = VIN#, Value = Car's JSON representation
// 	stub.PutState(vin, jsonVehicle)
// }

// Chaincode registers with the Shim on startup
func main() {
	fmt.Printf("Started Chaincode. token/history\n")
	err := shim.Start(new(VehicleChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}

/**
 * Getter function
 * function totalSupply() public view returns (uint);
 * Returns the totalSupply for the Pkrtoken
 **/
func totalSupply(stub shim.ChaincodeStubInterface) peer.Response {

	bytes, err := stub.GetState("Pkrtoken")
	if err != nil {
		return errorResponse(err.Error(), 5)
	}

	// Read the JSON and Initialize the struct
	var Pkr PkrToken
	_ = json.Unmarshal(bytes, &Pkr)

	// Create the JSON Response with totalSupply
	return successResponse(strconv.FormatUint(Pkr.TotalSupply, 10))
}

/**
 * Getter function
 * function balanceOf(address tokenOwner) public view returns (uint balance);
 * Returns the balance for the specified owner
 * {"Args":["balanceOf","cnic"]}
 **/
func balanceOf(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// Check if owner id is in the arguments
	if len(args) < 1 {
		return errorResponse("Needs OwnerID!!!", 6)
	}
	OwnerID := args[0]
	//contains cnic
	bytes, err := stub.GetState(OwnerPrefix + OwnerID)
	if err != nil {
		return errorResponse(err.Error(), 7)
	}
	//convert the balance in json and return it
	var wallet Wallet
	_ = json.Unmarshal(bytes, &wallet)
	response := balanceJSON(OwnerID, wallet)

	return successResponse(response)
}

/**
 * Setter function
 * function transfer(address to, uint tokens) public returns (bool success);
 * Transfer tokens
 * {"Args":["from","to","amount"]}
 **/
func transfer(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// Check if owner id is in the arguments
	if len(args) < 4 {
		return errorResponse("Needs to, from & amount!!!", 700)
	}
	// return errorResponse(args[0]+" ,,"+args[1]+" ,,"+args[2]+" ,,"+args[3],701)
	from := string(args[0])
	to := string(args[1])
	amount, err := strconv.ParseUint(args[2], 10, 64)
	if err != nil {
		return errorResponse(err.Error(), uint(amount))
	}
	if amount <= 0 {
		return errorResponse("Amount MUST be > 0!!!", 702)
	}

	// Get the Balance for from
	bytes, _ := stub.GetState(OwnerPrefix + from)

	if len(bytes) == 0 {
		// That means 0 token balance
		return errorResponse("Balance MUST be > 0!!!", 703)
	}
	var wallet Wallet
	_ = json.Unmarshal(bytes, &wallet)
	fromBalance, _ := strconv.ParseUint(wallet.Balance, 10, 64)

	if fromBalance < amount {
		return errorResponse("Insufficient balance to cover transfer!!!", 704)
	}
	// Reduce the tokens in from account
	fromBalance = fromBalance - amount

	// Get the balance in to account
	bytes, _ = stub.GetState(OwnerPrefix + to)
	var toBalance uint64= 0
	var wallet1 Wallet
	//if wallet exist than get his balance
	if len(bytes) > 0 {

		_ = json.Unmarshal(bytes, &wallet1)

		toBalance, _ = strconv.ParseUint(wallet1.Balance, 10, 64)
	}
	toBalance += amount

	// Update the balance
	// bytes = []byte(strconv.FormatInt(int64(fromBalance), 10))
	wallet.Balance = strconv.FormatInt(int64(fromBalance), 10)
	bytes1, _ := json.Marshal(wallet)
	err = stub.PutState(OwnerPrefix+from, []byte(bytes1))

	// bytes = []byte(strconv.FormatInt(int64(toBalance), 10))
	wallet1.Balance = strconv.FormatInt(int64(toBalance), 10)
	wallet1.Name = args[3] //setting name of to owner
	bytes2, _ := json.Marshal(wallet1)

	err = stub.PutState(OwnerPrefix+to, []byte(bytes2))

	// Emit Transfer Event
	display := "{\"from\":\"" + from + "\", \"to\":\"" + to + "\",\"amount\":" + strconv.FormatInt(int64(amount), 10) + "}"
	// stub.SetEvent("transfer", []byte(eventPayload))
	return successResponse(display)
}

// balanceJSON creates a JSON for representing the balance
func balanceJSON(OwnerID string, wallet Wallet) string {
	return "{\"owner-cnic\":\"" + OwnerID + "\", \"balance\":" + wallet.Balance + "\", \"name\":" + wallet.Name + "}"
}

func errorResponse(err string, code uint) peer.Response {
	codeStr := strconv.FormatUint(uint64(code), 10)
	// errorString := "{\"error\": \"" + err +"\", \"code\":"+codeStr+" \" }"
	errorString := "{\"error\":" + err + ", \"code\":" + codeStr + " \" }"
	return shim.Error(errorString)
}


func successResponse(dat string) peer.Response {
	success := "{\"response\": " + dat + ", \"code\": 0 }"
	return shim.Success([]byte(success))
}
