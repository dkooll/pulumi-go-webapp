package main

import (
	"fmt"
	"net/http"
	"os"
)

func hostname(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	fmt.Fprintf(w, "You are on host: %s", hostname)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi!")
}

func main() {
	http.HandleFunc("/hostname", hostname)
	http.HandleFunc("/", home)
	http.ListenAndServe(":80", nil)
}
