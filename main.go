package main

import (
	"fmt"
	"net/http"
)

// TODO:
// реализовать три сервиса event service/psotgresql
// реализовать booking service/postgresql, redis, kafka -> noti worker
// paymant service -> postgresql

func main() {
	mux := http.NewServeMux()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "Healthy") })

	http.ListenAndServe(":8080", mux)
}
