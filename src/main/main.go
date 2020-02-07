package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	ag "mirantis.com/tungsten-operator/tf-status-aggregator/src/aggregator"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	out := "TungsteFabric status aggregator\n"

	if r.URL.Path != "/" {
		out += "Unknown URL please use one of the following paths\n"
	}

	out += "\nRoutes:\n"
	out += "\n\t/pod-list\n\tReturns information about all detected \"tf-status\" pods\n\tTFStatus pod have to have following label: \"tungstenfabric\": \"status\"\n"
	out += "\n\t/status/json\n\tReturns agregated json from all \"tf-status\" pods.\n"
	out += "\n\t/status or /status/node\n\tReturns formated output from all detected \"tf-status\" pods\n\tat standart format for each node\n"
	out += "\n\t/status/group\n\tReturns formated output from all detected \"tf-status\" pod\n\tagregated by service for all nodes which handles the service"

	_, err := fmt.Fprintf(w, out)
	if err != nil {
		log.Fatal("URL: / ", err)
		http.Error(w, err.Error(), 500)
	}
}

func main() {
	ag := ag.New()
	h := Handler{}
	h.SetAggregator(ag)

	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler).Methods("GET")
	r.HandleFunc("/pod-list", h.handleListOfPod).Methods("GET")
	r.HandleFunc("/status/json", h.handleStatusJSON).Methods("GET")
	r.HandleFunc("/status/node", h.handleStatusPlainNode).Methods("GET")
	r.HandleFunc("/status", h.handleStatusPlainNode).Methods("GET")
	r.HandleFunc("/status/group", h.handleStatusPlainGroup).Methods("GET")
	http.Handle("/", r)

	serverPort := os.Getenv("SERVER_PORT")
	if len(serverPort) == 0 {
		serverPort = "80"
	}
	log.Printf("Starting server on port %s", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, nil))
}
