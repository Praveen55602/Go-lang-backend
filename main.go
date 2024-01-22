package main

import (
	"fmt"
	"net/http"
)

func HelloWorldPage(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "hello world")
}

func main() {
	http.HandleFunc("/", HelloWorldPage)
	http.ListenAndServe("", nil)
}
