package method

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"io/ioutil"
	"net/http"
	"smart-contact/conf"
	"strconv"
	"strings"
)

type NodeResp struct {
	JsonRpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  string `json:"result"`
}

// BodyInput for rpc call
type BodyInput struct {
	JsonRpc string        `json:"jsonrpc"`
	ID      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type CommonParams struct {
	To   string `json:"to"`
	From string `json:"from"`
	Data string `json:"data"`
}

func CallContractByRPC(blockNumber *int64, from, to, sign string, args ...interface{}) {
	body := &BodyInput{
		JsonRpc: "2.0",
		ID:      "1",
		Method:  "eth_call",
		Params:  genParam(blockNumber, from, to, sign, args),
	}

	jsonStr, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", conf.NodeURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("RPC result: ", string(data))
}

// Data in param
// first 64 bytes is keccak256 of signature ex: decimals() -> 0x313ce567000000000000000000000000
// then concat hex 64bytes of each args
func genParam(blockNumber *int64, from, to, sign string, args []interface{}) []interface{} {
	commonHash := "00000000000000000000000000000000000000000000000000000000000000000"
	signHash := crypto.Keccak256Hash([]byte(sign))
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

	if len(from) == 0 || from == "" {
		from = "0x0000000000000000000000000000000000000000"
	}

	common := &CommonParams{
		To:   to,
		From: from,
		Data: strings.ToLower(data),
	}

	blockHex := "latest"
	if blockNumber != nil {
		blockHex = "0x" + strconv.FormatInt(*blockNumber, 16)
	}

	return []interface{}{common, blockHex}
}
