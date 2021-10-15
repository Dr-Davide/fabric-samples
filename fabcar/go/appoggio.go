fmt.Println("Creating new Transactions:")
result, err = contract.SubmitTransaction("createGBtr", "GBTR22", "KW", "66", "FMV", "10", "sells", "Market")
if err != nil {
	fmt.Printf("Failed to submit transaction: %s\n", err)
	os.Exit(1)
}
//why is this not printing anything?
fmt.Println(string(result))

fmt.Println("Evaluating Transactions:")
result, err = contract.EvaluateTransaction("queryGBtr", "GBTR22")
if err != nil {
	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	os.Exit(1)
}
fmt.Println(string(result))

fmt.Println("All Transactions on the ledger:")
result, err = contract.EvaluateTransaction("queryAllGBtrs")
if err != nil {
	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	os.Exit(1)
}
fmt.Println(string(result))

// ------ backup functions

/*
	_, err = contract.SubmitTransaction("changeCarOwner", "CAR10", "Archie")
	if err != nil {
		fmt.Printf("Failed to submit transaction: %s\n", err)
		os.Exit(1)
	}

	result, err = contract.EvaluateTransaction("queryCar", "CAR10")
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))
*/