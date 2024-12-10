package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/iwanhae/uuid"
)

type ExampleStruct struct {
	ID uuid.V7 `json:"id"`
}

func main() {
	id := uuid.NewV7()

	// parse base58 formatted uuid v7
	{
		if err := id.Parse("CJ3mHjsWxTmqDzNe1TawS"); err != nil {
			panic(err)
		}

		// validates uuid version while parsing
		v4string := uuid.NewV4().String()
		if err := id.Parse(v4string); err != nil {
			// output:
			// expect uuid v7, but get uuid v4
			fmt.Println(err.Error())
		}
	}

	// could handle metadata from uuid
	{
		// output:
		// 2024-12-10 18:53:43.46 +0900 KST
		fmt.Println(id.Timestamp())
	}

	// print base58 formatted uuid v7
	{
		// output:
		// CJ3mHjsWxTmqDzNe1TawS
		fmt.Println(id) // returns base58 encoded uuid v7

		// output:
		// CJ3mHjsWxTmqDzNe1TawS
		fmt.Printf("%v\n", id)
	}

	// support "encoding/json" package
	{
		// output:
		// {"id":"CJ3mHjsWxTmqDzNe1TawS"}
		if err := json.NewEncoder(os.Stdout).Encode(ExampleStruct{
			ID: id,
		}); err != nil {
			panic(err)
		}

		v := ExampleStruct{}
		if err := json.Unmarshal(
			[]byte(`{"id":"CJ3mHjsWxTmqDzNe1TawS"}`), &v,
		); err != nil {
			panic(err)
		}

		// output:
		// {CJ3mHjsWxTmqDzNe1TawS}
		fmt.Println(v)
	}

}
