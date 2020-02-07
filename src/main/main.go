package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	ag "mirantis.com/tungsten-operator/tf-status-aggregator/src/aggregator"
)

func handle(w http.ResponseWriter, r *http.Request) {
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

	http.HandleFunc("/", handle)
	http.HandleFunc("/pod-list", h.handleListOfPod)
	http.HandleFunc("/status/json", h.handleStatusJSON)
	http.HandleFunc("/status/node", h.handleStatusPlainNode)
	http.HandleFunc("/status", h.handleStatusPlainNode)
	http.HandleFunc("/status/group", h.handleStatusPlainGroup)

	serverPort := os.Getenv("SERVER_PORT")
	if len(serverPort) == 0 {
		serverPort = "80"
	}
	log.Printf("Starting server on port %s", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, nil))
}
