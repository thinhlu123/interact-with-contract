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

	// erc20 contract
	var contract = "0x7130d2a12b9bcbfae4f2634d864a1ee1ce3ead9c"

	// method 1: Use abigen
	// STEP: Find the abi on bscscan => save to your code and run abigen in generate.go
	erc20, err := method.NewERC20Token(common.HexToAddress(contract), eth)
	if err != nil {
		panic(err)
	}
	decimal, err := erc20.Decimals(nil)
	fmt.Print("Abigen result: ")
	fmt.Println(decimal)

	// method 2: Use ethclient and package abi
	ethAbi := method.NewCustomClient(eth, contract, "abi/erc20.json")
	ethAbi.Call("balanceOf", nil, common.HexToAddress("0xe1bfa3d9994a88f81909fd5d8cef2642159c4e78"))
	ethAbi.Call("symbol", nil)

	// method 3: Use graphql
	method.InitGraphClient()
	method.GraphQLClient.Call(map[string]method.ArgumentEthCall{
		"query1": {
			Block: nil,
			To:    contract,
			Sign:  "decimals()",
			Args:  nil,
		},
	})

	// method 4: Use RPC
	method.CallContractByRPC(nil, "", contract, "decimals()")

	// method 5: use abi package with input and output like [ "address", "uint256",  "uint256", "uint256" ]
	ethAbi = method.NewCustomClientWithString(eth, contract, method.GenAbi([]string{"address"}, []string{"uint256"}))
	ethAbi.Call("balanceOf", nil, common.HexToAddress("0xe1bfa3d9994a88f81909fd5d8cef2642159c4e78"))
	ethAbi.Call("symbol", nil)
}
