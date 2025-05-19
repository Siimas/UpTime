package util

import (
	"fmt"
	"encoding/json"
)

func PrettyPrint(v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}