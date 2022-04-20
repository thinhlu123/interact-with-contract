package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"os"
)

func GenByAbi(signature string, args ...interface{}) string {

	file, err := os.Open("")
	if err != nil {
		panic(err)
	}
	parsed, err := abi.JSON(file)
	if err != nil {
		panic(err)
	}

	input, err := parsed.Pack(signature, args...)
	if err != nil {
		log.Println(err)
	}

	return string(input)
}

// Data in param
// first 64 bytes is keccak256 of signature ex: decimals() -> 0x313ce567000000000000000000000000
// then concat hex 64bytes of each args
func genParam(funcName string, args []interface{}) string {
	commonHash := "00000000000000000000000000000000000000000000000000000000000000000"
	signHash := crypto.Keccak256Hash([]byte(funcName))
	hashString := signHash.String()
	data := hashString[:10]

	for i, arg := range args {
		switch arg.(type) {
		case int64:
			h := fmt.Sprintf("%x", arg)
			length := len(h)
			if length < 64 {
				h = commonHash[:64-length] + h
			}

			data += h
		case float64:
			argFloat := arg.(float64)
			argInt := int(argFloat)
			h := fmt.Sprintf("%x", argInt)
			length := len(h)
			if i < len(args)-1 && length < 64 {
				h = commonHash[:64-length] + h
			}

			data += h
		case string:
			argStr := arg.(string)
			data += "000000000000000000000000" + argStr[2:]
		}
	}

	return data
}
