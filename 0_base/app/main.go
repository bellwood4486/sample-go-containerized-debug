package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "hello world")
	})
	addr := ":8080"
	fmt.Printf("Launching server at %q ...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
