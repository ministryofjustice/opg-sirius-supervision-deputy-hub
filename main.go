package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":1234", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world!")
}
