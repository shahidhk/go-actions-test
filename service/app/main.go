package main

import (
	"log"
	"net/http"
	"os"
)

var port string

func init() {
	port = os.Getenv("PORT")
	if len(port) == 0 {
		port = "4000"
	}
}

func main() {
	var handler = http.NewServeMux()

	// add routers here
	handler.HandleFunc("/foo", fooHandler)
	handler.HandleFunc("/createAddress", createAddressHandler)

	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// copy routes inside main

// handler.HandleFunc("/actionName", actionNameHandler)
