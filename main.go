package main

import (
	"fmt"
	myHttp "net/http" //here myHttp is given as alias to the package
)

func HelloWorldPage(res myHttp.ResponseWriter, req *myHttp.Request) {
	res.Header().Set("Content-Type", "text/html") //by default it's text/html only
	fmt.Fprintf(res, "<h1>hello world</h1>")
}

func main() {
	myHttp.HandleFunc("/", HelloWorldPage)
	//myHttp.ListenAndServe("", nil) //by default this creates a instance of server and starts listening to the provided port we can create our own server so that we can get more flexibity with the setting of server
	myCustomServer := myHttp.Server{
		Addr:         "",
		Handler:      nil,
		ReadTimeout:  1000,
		WriteTimeout: 1000,
	}
	myCustomServer.ListenAndServe()
}
