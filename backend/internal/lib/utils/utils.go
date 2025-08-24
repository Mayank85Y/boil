package utils

import (
	"encoding/json" 
	"fmt"
)

func PrintJSON(v interface{}) {
	json, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		fmt.Println("error marshalling to json:", err)
		return
	}
	fmt.Println("JSON", string(json))
}