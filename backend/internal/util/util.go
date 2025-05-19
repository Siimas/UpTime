package util

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
