package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "opg-sirius-deputy-hub", log.LstdFlags)

	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":1234", nil)
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		logger.Fatalln(err)
	}
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world!")
}
