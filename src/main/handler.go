package main

import (
	"fmt"
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
	w.Header().Add("Content-Type", "text/plain")
	podList := h.ag.GetPodList()
	var out string
	for _, pod := range podList.Items {
		out += fmt.Sprintf("---------------------------------------------------\n")
		out += fmt.Sprintf("Name: %v\n", pod.Name)
		out += fmt.Sprintf("IP: %v\n", pod.Status.PodIP)
		out += fmt.Sprintf("UUID: %v\n", pod.UID)
		out += fmt.Sprintf("HostIP: %v\n", pod.Status.HostIP)
		out += fmt.Sprintf("Host: %v\n\n", pod.Spec.NodeName)
	}

	fmt.Println("Requested: /pod-list")
	fmt.Fprintf(w, out)
}

func (h *Handler) handleStatusJSON(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Requested: /status/json")
	w.Header().Add("Content-Type", "application/json")
	jsons := h.ag.GetJSONStatusForAll()
	out := strings.Join(jsons, ",")
	fmt.Fprintf(w, out)

}

func (h *Handler) handleStatusPlainNode(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Requested: /status/node")
	w.Header().Add("Content-Type", "plan/text")
	out := h.ag.GetPlainStatusByNode()

	fmt.Fprintf(w, out)

}

func (h *Handler) handleStatusPlainGroup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Requested: /status/group")
	w.Header().Add("Content-Type", "plan/text")
	out := h.ag.GetPlainStatusByGroup()

	fmt.Fprintf(w, out)

}
