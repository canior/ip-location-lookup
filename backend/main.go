package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"ip-lookup-app/api"
	"ip-lookup-app/cmd"
	"log"
	"net/http"
)

func main() {
	go cmd.RunIpsLocationLookup()

	clientChannelMap := make(map[string]chan string)

	r := mux.NewRouter()
	r.HandleFunc("/stream/{clientId}", func(w http.ResponseWriter, r *http.Request) {
		api.SseHandler(w, r, clientChannelMap)
	}).Methods("GET")

	r.HandleFunc("/ips/{clientId}", func(w http.ResponseWriter, r *http.Request) {
		api.IpsLookupHandler(w, r, clientChannelMap)
	})

	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		fmt.Fprint(w, "Hello")
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
