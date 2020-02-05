package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	ag "mirantis.com/tungsten-operator/tf-status-agregator/src/agregator"
)

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Requested: /")
	out := "TungsteFabric status agregator\n"
	out += "\nRoutes:\n"
	out += "\n\t/pod-list\n\tReturns information about all detected \"tf-status\" pods\n\tTFStatus pod have to have following label: \"tungstenfabric\": \"status\"\n"
	out += "\n\t/status/json\n\tReturns agregated json from all \"tf-status\" pods.\n"
	out += "\n\t/status or /status/node\n\tReturns formated output from all detected \"tf-status\" pods\n\tat standart format for each node\n"
	out += "\n\t/status/group\n\tReturns formated output from all detected \"tf-status\" pod\n\tagregated by service for all nodes which handles the service"
	fmt.Fprintf(w, out)
}

func main() {
	ag := ag.New()
	h := Handler{}
	h.SetAgregator(ag)

	fmt.Println("Starting server...")

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
	log.Fatal(http.ListenAndServe(":"+serverPort, nil))
}
