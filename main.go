package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", greeter)
	err := http.ListenAndServe(":5656", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func greeter(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello Yopta")
}
