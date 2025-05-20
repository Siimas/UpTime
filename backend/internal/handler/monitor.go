package handler

import (
	"fmt"
	"net/http"
)

func Test(w http.ResponseWriter, req *http.Request)  {
	fmt.Fprintf(w, "yoo")
}