package main

import (
	"fmt"
	"net/http"
)

func main() {
	multiplexer := http.NewServeMux()

	server := http.Server{
		Handler: multiplexer,
		Addr:    ":8080",
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
