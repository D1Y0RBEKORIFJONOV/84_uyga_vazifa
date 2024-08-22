package main

import (
	"expvar"
	"net/http"
)

var (
	requestCount = expvar.NewInt("request_count")
)

func handler(w http.ResponseWriter, r *http.Request) {
	requestCount.Add(1)
	w.Write([]byte("Hello, World!"))
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
