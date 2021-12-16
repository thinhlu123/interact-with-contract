package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"smart-contact/conf"
	"smart-contact/method"
)

func main() {
	// TODO: Change NodeURL in conf folder to your node
	// make sure your node has enable graphql mode
	eth, err := ethclient.Dial(conf.NodeURL)
	if err != nil {
		panic(err)
	}

	// method 1: Use abigen
	// STEP: Find the abi on bscscan => save to your code and run abigen in generate.go
	erc20, err := method.NewERC20Token(common.HexToAddress("0x7130d2a12b9bcbfae4f2634d864a1ee1ce3ead9c"), eth)
	if err != nil {
		panic(err)
	}
	decimal, err := erc20.Decimals(nil)
	fmt.Print("Abigen result: ")
	fmt.Println(decimal)

	// method 2: Use ethclient and abi
	ethAbi := method.NewCustomClient(eth, "0x7130d2a12b9bcbfae4f2634d864a1ee1ce3ead9c", "abi/erc20.json")
	ethAbi.Call("decimals", nil)

	// method 3: Use graphql
	method.InitGraphClient()
	method.GraphQLClient.Call(map[string]method.ArgumentEthCall{
		"query1": {
			Block: nil,
			To:    "0x7130d2a12b9bcbfae4f2634d864a1ee1ce3ead9c",
			Sign:  "decimals()",
			Args:  nil,
		},
	})

	// method 4: Use RPC
	method.CallContractByRPC(nil, "", "0x7130d2a12b9bcbfae4f2634d864a1ee1ce3ead9c", "decimals()")
}
