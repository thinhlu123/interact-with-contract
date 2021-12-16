package method

import (
	"context"
	"fmt"
	ether "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"os"
)

type CustomClient struct {
	ethClient *ethclient.Client
	abi       abi.ABI
	contract  ethereum.Address
}

func NewCustomClient(ethClient *ethclient.Client, contract, abiPath string) *CustomClient {
	//_, fileLocation, _, _ := runtime.Caller(1)
	//abiPath = filepath.Join(fileLocation, abiPath)
	file, err := os.Open(abiPath)
	if err != nil {
		panic(err)
	}
	parsed, err := abi.JSON(file)
	if err != nil {
		panic(err)
	}

	return &CustomClient{
		ethClient: ethClient,
		abi:       parsed,
		contract:  ethereum.HexToAddress(contract),
	}
}

func (bl *CustomClient) Call(signature string, blockNumber *big.Int, args ...interface{}) {
	input, err := bl.abi.Pack(signature, args...)
	if err != nil {
		log.Println(err)
	}
	value := big.NewInt(0)
	msg := ether.CallMsg{To: &bl.contract, Value: value, Data: input}
	result, err := bl.ethClient.CallContract(context.Background(), msg, blockNumber)
	if err != nil {
		fmt.Println(err)
		return
	}

	var trueResult interface{}
	trueResult, err = bl.abi.Unpack(signature, result)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("EthClient_ABI result: ", trueResult)
}