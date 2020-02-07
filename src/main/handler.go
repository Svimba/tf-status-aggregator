package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"mirantis.com/tungsten-operator/tf-status-aggregator/src/aggregator"
)

// Handler struct
type Handler struct {
	ag *aggregator.Aggregator
}

// SetAggregator func
func (h *Handler) SetAggregator(ag *aggregator.Aggregator) {
	h.ag = ag
}

func (h *Handler) handleListOfPod(w http.ResponseWriter, r *http.Request) {

	log.Println("Requested: " + r.URL.Path)
	podList, err := h.ag.GetPodList()
	if err != nil {
		log.Fatal("URL: "+r.URL.Path, err)
		http.Error(w, err.Error(), 500)
		return
	}

	var out string
	for _, pod := range podList.Items {
		out += fmt.Sprintf("---------------------------------------------------\n")
		out += fmt.Sprintf("Name: %v\n", pod.Name)
		out += fmt.Sprintf("IP: %v\n", pod.Status.PodIP)
		out += fmt.Sprintf("UUID: %v\n", pod.UID)
		out += fmt.Sprintf("HostIP: %v\n", pod.Status.HostIP)
		out += fmt.Sprintf("Host: %v\n\n", pod.Spec.NodeName)
	}

	w.Header().Add("Content-Type", "text/plain")
	_, err = fmt.Fprintf(w, out)
	if err != nil {
		log.Fatal("URL: "+r.URL.Path, err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h *Handler) handleStatusJSON(w http.ResponseWriter, r *http.Request) {
	log.Println("Requested: " + r.URL.Path)

	jsons, err := h.ag.GetJSONStatusForAll()
	if err != nil {
		log.Fatal("URL: "+r.URL.Path, err)
		http.Error(w, err.Error(), 500)
		return
	}
	out := strings.Join(jsons, ",")

	w.Header().Add("Content-Type", "application/json")
	_, err = fmt.Fprintf(w, out)
	if err != nil {
		log.Fatal("URL: "+r.URL.Path, err)
		http.Error(w, err.Error(), 500)
		return
	}

}

func (h *Handler) handleStatusPlainNode(w http.ResponseWriter, r *http.Request) {
	log.Println("Requested: " + r.URL.Path)

	out, err := h.ag.GetPlainStatusByNode()
	if err != nil {
		log.Fatal("URL: "+r.URL.Path, err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Add("Content-Type", "plan/text")
	_, err = fmt.Fprintf(w, out)
	if err != nil {
		log.Fatal("URL: "+r.URL.Path, err)
		http.Error(w, err.Error(), 500)
		return
	}

}

func (h *Handler) handleStatusPlainGroup(w http.ResponseWriter, r *http.Request) {
	log.Println("Requested: " + r.URL.Path)

	out, err := h.ag.GetPlainStatusByGroup()
	if err != nil {
		log.Fatal("URL: "+r.URL.Path, err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Add("Content-Type", "plan/text")
	_, err = fmt.Fprintf(w, out)
	if err != nil {
		log.Fatal("URL: "+r.URL.Path, err)
		http.Error(w, err.Error(), 500)
		return
	}
}
