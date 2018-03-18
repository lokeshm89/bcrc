package main

import (
	"fmt"
	"strconv"
	"time"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type MilkProduct struct{
	BatchId		string	`json:"batchId"`
	Farm		string	`json:"farm"`
	Quantity	int	`json:"quantity"`
	Fat		string	`json:"fat"`
	Protein		string	`json:"protein"`
	Juicines	string	`json:"juicines"`
	Lead		string	`json:"lead"`
	BatchStatus	string	`json:"batchStatus"`
	Owner		string	`json:"owner"`
	Temprature	string	`json:"temprature"`
	Amount		string	`json:"amount"`
	Transactions []Transaction	`json:"transactions"`

}

type Transaction struct {

	Farm		string	`json:"farm"`
	Fat		string	`json:"fat"`
	Protein		string	`json:"protein"`
	Juicines	string	`json:"juicines"`
	Lead		string	`json:"lead"`
	Pasturization	string	`json:"pasturization"`
	Adulteration	string	`json:"adulteration"`
	NoOfBottles	int	`json:"noOfBottles"`
	TransName	string	`json:"transName"`
	Status		string	`json:"status"`
	TransTime	string	`json:"transTime"`
	TPstatus	string	`json:"tpStatus"`
	RTstatus	string	`json:"rtStatus"`
	Owner		string	`json:"owner"`
	Temprature	string	`json:"temprature"`

}

type AllBatches struct{
	AllBatches []string
}

type AllBatchesDetails struct{
	Batches []MilkProduct
}

type TimeTracker struct{
	DispachedTime	string
	ReachedTime	string
	TPTemprature	string
}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("chaincode_custom Init")
	var err error
	var batches AllBatches
	
	batchesAsBytes,_ :=json.Marshal(batches)
	err = stub.PutState("AllBatches", batchesAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("chaincode_custom Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "addProduct"{
		return t.addProduct(stub, args)
	} else if function == "getAllManuFacturerBatches"{
		return t.getAllManuFacturerBatches(stub, args)
	} else if function == "transferProduct"{
		return t.transferProduct(stub,args)
	}else if function == "addProccesingInput"{
		return t.addProccesingInput(stub,args)	
	} else if function == "getAllProcessingBatches"{
		return t.getAllProcessingBatches(stub,args)
	} else if function == "getAllTransporterBatches"{
		return t.getAllTransporterBatches(stub,args)
	} else if function == "getAllRetailerBatches"{
		return t.getAllRetailerBatches(stub,args)
	}else if function == "addIOTResults"{
		return t.addIOTResults(stub,args)
	} else if function == "retailerQualityCheck"{
		return t.retailerQualityCheck(stub,args)
	} else if function == "getTimeTracker"{
		return t.getTimeTracker(stub,args)	
	}else if function == "getAllConsumerBatches"{
		return t.getAllConsumerBatches(stub,args)	
	}
	
	return shim.Error("Invalid invoke function name.")
}

//func to add milk product
func (t *SimpleChaincode) addProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Add product...")
	var product MilkProduct
	var err error

	isExistAsBytes,_ := stub.GetState(args[0])

	if(isExistAsBytes != nil){
		return shim.Error(err.Error())
	}

	product.BatchId=args[0]
	product.Farm=args[1]
	product.Quantity, err=strconv.Atoi(args[2])
	product.Fat=args[3]
	product.Protein=args[4]
	product.Juicines=args[5]
	product.Lead=args[6]
	product.Owner=args[7]
	product.Amount=args[8]
	product.BatchStatus="MF"
	product.Temprature=""

	var tx Transaction

	tx.Farm=product.Farm
	tx.Fat=product.Fat
	tx.Protein=product.Protein
	tx.Juicines=product.Juicines
	tx.Lead=product.Lead
	tx.Pasturization=""
	tx.Adulteration=""
	tx.NoOfBottles=0
	tx.TransName=""
	tx.Status="MF"
	tx.TransTime=time.Now().UTC().String()
	tx.TPstatus=""
	tx.RTstatus=""
	tx.Owner=product.Owner
	tx.Temprature=""
	product.Transactions=append(product.Transactions,tx)

	productAsBytes, _ := json.Marshal(product)
	err = stub.PutState(product.BatchId, productAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//update all batches
	AllBatchesAsBytes, err := stub.GetState("AllBatches")

	var allBatches AllBatches

	err= json.Unmarshal(AllBatchesAsBytes, &allBatches)
	allBatches.AllBatches=append(allBatches.AllBatches,product.BatchId)

	allbatchesAsBytes,_ :=json.Marshal(allBatches)
	err = stub.PutState("AllBatches", allbatchesAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//fill time tracker
	var timeTracker TimeTracker
		timeTracker.DispachedTime=time.Now().UTC().String()
		timeTracker.ReachedTime=""
		timeTracker.TPTemprature=""
	
	timeTrackerAsBytes,_ :=json.Marshal(timeTracker)
	err = stub.PutState("time"+product.BatchId, timeTrackerAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}	


	return shim.Success(nil)
}

// ============================================================================================================================
// Get All Batches Details for processing input
// ============================================================================================================================
func (t *SimpleChaincode) getAllProcessingBatches(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	
	//get the AllBatches index
	allBAsBytes,_ := stub.GetState("AllBatches")
	
	var res AllBatches
	json.Unmarshal(allBAsBytes, &res)
	
	var rab AllBatchesDetails

	for i := range res.AllBatches{

		sbAsBytes,_ := stub.GetState(res.AllBatches[i])
		
		var sb MilkProduct
		json.Unmarshal(sbAsBytes, &sb)

	if(sb.BatchStatus == "In Process") {
		rab.Batches = append(rab.Batches,sb); 
	}
	}

	rabAsBytes, _ := json.Marshal(rab)

	return shim.Success(rabAsBytes)
	
}

// ============================================================================================================================
// Get All Batches Details for Transporter
// ============================================================================================================================
func (t *SimpleChaincode) getAllTransporterBatches(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	
	//get the AllBatches index
	allBAsBytes,_ := stub.GetState("AllBatches")
	
	var res AllBatches
	json.Unmarshal(allBAsBytes, &res)
	
	var rab AllBatchesDetails

	for i := range res.AllBatches{

		sbAsBytes,_ := stub.GetState(res.AllBatches[i])
		
		var sb MilkProduct
		json.Unmarshal(sbAsBytes, &sb)

	if(sb.BatchStatus == "TP" || sb.BatchStatus == "RT" || sb.BatchStatus == "Paid") {
		rab.Batches = append(rab.Batches,sb); 
	}

	}

	rabAsBytes, _ := json.Marshal(rab)

	return shim.Success(rabAsBytes)
	
}

// ============================================================================================================================
// Get All Batches Details for Transporter
// ============================================================================================================================
func (t *SimpleChaincode) getAllManuFacturerBatches(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	
	//get the AllBatches index
	allBAsBytes,_ := stub.GetState("AllBatches")
	
	var res AllBatches
	json.Unmarshal(allBAsBytes, &res)
	
	var rab AllBatchesDetails

	for i := range res.AllBatches{

		sbAsBytes,_ := stub.GetState(res.AllBatches[i])
		
		var sb MilkProduct
		json.Unmarshal(sbAsBytes, &sb)

	//if(sb.BatchStatus == "MF") {
		rab.Batches = append(rab.Batches,sb); 
	//}

	}

	rabAsBytes, _ := json.Marshal(rab)

	return shim.Success(rabAsBytes)
}

// ============================================================================================================================
// Get All Batches Details for Transporter
// ============================================================================================================================
func (t *SimpleChaincode) getAllRetailerBatches(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	
	//get the AllBatches index
	allBAsBytes,_ := stub.GetState("AllBatches")
	
	var res AllBatches
	json.Unmarshal(allBAsBytes, &res)
	
	var rab AllBatchesDetails

	for i := range res.AllBatches{

		sbAsBytes,_ := stub.GetState(res.AllBatches[i])
		
		var sb MilkProduct
		json.Unmarshal(sbAsBytes, &sb)

	if(sb.BatchStatus == "RT" || sb.BatchStatus == "Paid") {
		rab.Batches = append(rab.Batches,sb); 
	}

	}

	rabAsBytes, _ := json.Marshal(rab)

	return shim.Success(rabAsBytes)
	
}


// ============================================================================================================================
// Transfer a batch - can only be done by shipping company
// ============================================================================================================================
func (t *SimpleChaincode) transferProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

	//Update Batch data
	bAsBytes, err := stub.GetState(args[0])
	
	var bch MilkProduct
	err = json.Unmarshal(bAsBytes, &bch)
	if err != nil {
		return shim.Error(err.Error())
	}
	if(bch.BatchStatus != "In Process"){
		return shim.Error(err.Error())
	}

	bch.BatchStatus="TP"
	var tx Transaction
	tx.Farm=bch.Farm
	tx.Fat=bch.Fat
	tx.Protein=bch.Protein
	tx.Juicines=bch.Juicines
	tx.Lead=bch.Lead
	tx.NoOfBottles, err=strconv.Atoi(args[1])
	tx.TransName=args[2]
	tx.Status="TP"
	tx.TransTime=time.Now().UTC().String()
	tx.Owner=tx.TransName
	tx.TPstatus="in-transist"
	bch.Owner=tx.TransName

	bch.Transactions = append(bch.Transactions, tx)

	//Commit updates batch to ledger
	btAsBytes, _ := json.Marshal(bch)
	err = stub.PutState(bch.BatchId, btAsBytes)	
	if err != nil {
		return shim.Error(err.Error())
	}	

	return shim.Success(nil)
}


func (t *SimpleChaincode) addProccesingInput(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	
	//Update Batch data
	bAsBytes, err := stub.GetState(args[0])
	
	var bch MilkProduct
	err = json.Unmarshal(bAsBytes, &bch)
	if err != nil {
		return shim.Error(err.Error())
	}
	if(bch.BatchStatus != "MF"){
		return shim.Error(err.Error())
	}

	bch.BatchStatus="In Process"
	var tx Transaction
	
	tx.Farm=bch.Farm
	tx.Fat=bch.Fat
	tx.Protein=bch.Protein
	tx.Juicines=bch.Juicines
	tx.Lead=bch.Lead
	tx.Pasturization=args[1]
	tx.Adulteration=args[2]
	tx.Status="In Process"
	tx.TransTime=time.Now().UTC().String()
	tx.Temprature=args[3]
	bch.Transactions = append(bch.Transactions, tx)

	//Commit updates batch to ledger
	btAsBytes, _ := json.Marshal(bch)
	err = stub.PutState(bch.BatchId, btAsBytes)	
	if err != nil {
		return shim.Error(err.Error())
	}	

	return shim.Success(nil)
}


//Add iot device results
func (t *SimpleChaincode) addIOTResults(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

	//Update Batch data
	bAsBytes, err := stub.GetState(args[0])
	
	var bch MilkProduct
	var tx Transaction
	err = json.Unmarshal(bAsBytes, &bch)
	if err != nil {
		return shim.Error(err.Error())
	}

	if(bch.Transactions[len(bch.Transactions)-1].TPstatus == "Delivered" || bch.Transactions[len(bch.Transactions)-1].TPstatus == "Rejected"){

		return shim.Error(err.Error())
	}

	var temp int
		temp, err=strconv.Atoi(args[2])
	if temp < 15{
		bch.BatchStatus="RT"
		tx.Farm=bch.Farm
		tx.Fat=bch.Fat
		tx.Protein=bch.Protein
		tx.Juicines=bch.Juicines
		tx.Lead=bch.Lead
		tx.Status="RT"
		tx.TransTime=time.Now().UTC().String()
		tx.TransName=bch.Transactions[len(bch.Transactions)-1].TransName
		tx.NoOfBottles, err=strconv.Atoi(args[1])
		tx.TPstatus="Delivered"
		tx.Temprature=args[2]
	} else {
		bch.BatchStatus="TP"
		tx.Farm=bch.Farm
		tx.Fat=bch.Fat
		tx.Protein=bch.Protein
		tx.Juicines=bch.Juicines
		tx.Lead=bch.Lead
		tx.Status="TP"
		tx.TransTime=time.Now().UTC().String()
		tx.TransName=bch.Transactions[len(bch.Transactions)-1].TransName
		tx.NoOfBottles, err=strconv.Atoi(args[1])
		tx.TPstatus="Rejected"
		tx.Temprature=args[2]
	}

	bch.Transactions = append(bch.Transactions, tx)

	//Commit updates batch to ledger
	btAsBytes, _ := json.Marshal(bch)
	err = stub.PutState(bch.BatchId, btAsBytes)	
	if err != nil {
		return shim.Error(err.Error())
	}

	//fill time tracker
	var timeTracker TimeTracker
		timeTracker.ReachedTime=time.Now().UTC().String()
		timeTracker.TPTemprature=args[2]
	
	timeTrackerAsBytes,_ :=json.Marshal(timeTracker)
	err = stub.PutState("time"+bch.BatchId, timeTrackerAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}	

	return shim.Success(nil)
}

//Retailer quality check
func (t *SimpleChaincode) retailerQualityCheck(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	
	//Update Batch data
	bAsBytes, err := stub.GetState(args[0])
	
	var bch MilkProduct
	var tx Transaction
	err = json.Unmarshal(bAsBytes, &bch)
	if err != nil {
		return shim.Error(err.Error())
	}

	var check string
		check=args[1]
	if check == "accept"{

		bch.BatchStatus="Paid"
		tx.Status="Paid"
		tx.TransTime=time.Now().UTC().String()
		tx.RTstatus="Accepted"
		tx.Farm=bch.Farm
		tx.Fat=bch.Fat
		tx.Protein=bch.Protein
		tx.Juicines=bch.Juicines
		tx.Lead=bch.Lead
		tx.TransName=bch.Transactions[len(bch.Transactions)-1].TransName
		tx.NoOfBottles=bch.Transactions[len(bch.Transactions)-1].NoOfBottles
		tx.TPstatus=bch.Transactions[len(bch.Transactions)-1].TPstatus
		tx.Temprature=bch.Transactions[len(bch.Transactions)-1].Temprature

	} else {

		bch.BatchStatus="RT"
		tx.Farm=bch.Farm
		tx.Fat=bch.Fat
		tx.Protein=bch.Protein
		tx.Juicines=bch.Juicines
		tx.Lead=bch.Lead
		tx.Status="RT"
		tx.TransTime=time.Now().UTC().String()
		tx.RTstatus="Rejected"
		tx.TransName=bch.Transactions[len(bch.Transactions)-1].TransName
		tx.NoOfBottles=bch.Transactions[len(bch.Transactions)-1].NoOfBottles
		tx.TPstatus=bch.Transactions[len(bch.Transactions)-1].TPstatus
		tx.Temprature=bch.Transactions[len(bch.Transactions)-1].Temprature

	}
	bch.Transactions = append(bch.Transactions, tx)

	//Commit updates batch to ledger
	btAsBytes, _ := json.Marshal(bch)
	err = stub.PutState(bch.BatchId, btAsBytes)	
	if err != nil {
		return shim.Error(err.Error())
	}	
	return shim.Success(nil)
}

// ============================================================================================================================
// Get time of the Batch
// ============================================================================================================================
func (t *SimpleChaincode) getTimeTracker(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	
	//get the time tracker
	timeTrackerAsBytes,_ := stub.GetState("time"+args[0])

	return shim.Success(timeTrackerAsBytes)
}

//Get all batches that have status PAID in android App
func (t *SimpleChaincode) getAllConsumerBatches(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	
	//get the AllBatches index
	allBAsBytes,_ := stub.GetState("AllBatches")
	
	var res AllBatches
	json.Unmarshal(allBAsBytes, &res)
	
	var rab AllBatchesDetails

	for i := range res.AllBatches{

		sbAsBytes,_ := stub.GetState(res.AllBatches[i])
		
		var sb MilkProduct
		json.Unmarshal(sbAsBytes, &sb)

	if(sb.BatchStatus == "Paid") {
		rab.Batches = append(rab.Batches,sb); 
	}

	}

	rabAsBytes, _ := json.Marshal(rab)

	return shim.Success(rabAsBytes)
	
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
