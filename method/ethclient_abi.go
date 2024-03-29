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
	"strings"
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

func NewCustomClientWithString(ethClient *ethclient.Client, contract, abiStr string) *CustomClient {
	//_, fileLocation, _, _ := runtime.Caller(1)
	//abiPath = filepath.Join(fileLocation, abiPath)
	abiReader := strings.NewReader(abiStr)
	parsed, err := abi.JSON(abiReader)
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

	trueResult := make(map[string]interface{})
	err = bl.abi.UnpackIntoMap(trueResult, signature, result)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("EthClient_ABI result: ", trueResult)
}

func (bl *CustomClient) BalanceOf(wallet string, blockNumber *big.Int) (*big.Int, error) {
	return bl.ethClient.BalanceAt(context.Background(), ethereum.HexToAddress(wallet), blockNumber)
}

func GenAbi(input []string, output []string) string {
	in, out := "", ""
	for _, item := range input {
		in += fmt.Sprintf(`{
        "internalType": "%v",
        "name": "",
        "type": "%v"
      }`, item, item)
	}

	for _, item := range output {
		out += fmt.Sprintf(`{
        "internalType": "%v",
        "name": "",
        "type": "%v"
      }`, item, item)
	}

	return fmt.Sprintf(`[{
        "constant": true,
        "inputs": [
      %s
    ],
        "name": "balanceOf",
        "outputs": [%s],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    }]`, in, out)
}
