package method

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/machinebox/graphql"
	"smart-contact/conf"
	"strings"
)

type graphQlClient struct {
	client *graphql.Client
}

var GraphQLClient graphQlClient

func InitGraphClient() {
	GraphQLClient = graphQlClient{
		graphql.NewClient(conf.NodeURL + "/graphql"),
	}
}

type ArgumentEthCall struct {
	Block    *int64
	To, Sign string
	Args     []interface{}
}

// expect result
//query {
//	q1: block {
//		call(data: {to: $to, data: $data}) {
//			data
//		}
//	}
//	q2: block {
//		call(data: {to: $to, data: $data}) {
//			data
//		}
//	}
//}
func genQuery(listQuery map[string]ArgumentEthCall) string {
	var query string
	for k, v := range listQuery {
		query += fmt.Sprintf("%s:%s\n", k, constructEthCallQuery(v))
	}

	return fmt.Sprintf("{%s}", query)
}

func constructEthCallQuery(args ArgumentEthCall) string {
	data := genData(args.Sign, args.Args)
	if args.Block == nil {
		return fmt.Sprintf(`block {
		call(data: { to: "%s", data: "%s" }) {
			data
		}
	}`, args.To, data)
	}

	return fmt.Sprintf(`block(number: %x) {
		call(data: { to: "%s", data: "%s" }) {
			data
		}
	}`, args.Block, args.To, data)
}

func genData(sign string, args []interface{}) string {
	commonHash := "00000000000000000000000000000000000000000000000000000000000000000"
	signHash := crypto.Keccak256Hash([]byte(sign))
	hashString := signHash.String()
	data := hashString[:10] + "000000000000000000000000"

	for i, arg := range args {
		switch arg.(type) {
		case int64:
			h := fmt.Sprintf("%x", arg)
			length := len(h)
			if length < 64 {
				h += commonHash[:64-length]
			}

			data += h
		case float64:
			argFloat := arg.(float64)
			argInt := int(argFloat)
			h := fmt.Sprintf("%x", argInt)
			length := len(h)
			if i < len(args)-1 && length < 64 {
				h += commonHash[:64-length]
			}

			data += h
		case string:
			argStr := arg.(string)
			data += argStr[2:] + "000000000000000000000000"
		}
	}

	return strings.ToLower(data)
}

func parseResult(data string) []string {
	data = data[2:]
	length := len(data)
	num := length / 64
	rs := make([]string, num)
	for i := 0; i < num; i++ {
		rs[i] = data[i*64 : (i+1)*64]
	}

	return rs
}

func (c *graphQlClient) Call(args map[string]ArgumentEthCall) {
	query := genQuery(args)

	graphqlRequest := graphql.NewRequest(query)
	var graphqlResponse map[string]map[string]map[string]string
	err := c.client.Run(context.Background(), graphqlRequest, &graphqlResponse)
	if err != nil {
		fmt.Println(err)
	}

	rs := make(map[string][]string)
	for k, v := range graphqlResponse {
		rs[k] = parseResult(v["call"]["data"])
	}

	fmt.Print("Graphql result: ")
	fmt.Println(rs)
}
