package aggregator

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/olekukonko/tablewriter"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// Aggregator struct
type Aggregator struct {
	client  client.Client
	context context.Context
}

// New - constructor for Aggregator
func New() *Aggregator {
	ag := Aggregator{}
	ag.setClient()

	ag.context = context.TODO()

	return &ag
}

func (ag *Aggregator) setClient() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err, "Cannot get K8s config")
		os.Exit(1)
	}
	cli, err := client.New(cfg, client.Options{})
	if err != nil {
		log.Fatal(err, "Cannot create new client")
		os.Exit(1)
	}
	ag.client = cli
}

// GetPodList func
func (ag *Aggregator) GetPodList() (*corev1.PodList, error) {
	podList := &corev1.PodList{}
	err := ag.client.List(ag.context, podList, client.MatchingLabels{"tungstenfabric": "status"})
	if err != nil {
		log.Fatalf("Cannot get list of tf-status nodes %v", err)
		return nil, err
	}
	return podList, err
}

// GetPodJSONStatus func
func (ag *Aggregator) GetPodJSONStatus(podIP string, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Requesting %s for status\n", podIP)

	nodePort := os.Getenv("NODE_PORT")
	if len(nodePort) == 0 {
		nodePort = "80"
	}

	resp, err := http.Get("http://" + podIP + ":" + nodePort + "/json")
	if err != nil {
		log.Fatalf("Cannot get status from %s, error:  %v\n", podIP, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Cannot read output from %s, error: %v\n", podIP, err)
	}
	log.Printf("Status from %s has been received successfully\n", podIP)
	ch <- string(body)
}

// GetJSONStatusForAll func
func (ag *Aggregator) GetJSONStatusForAll() ([]string, error) {
	podList, err := ag.GetPodList()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	channel := make(chan string, podList.Size())
	jsonOuts := []string{}

	for _, pod := range podList.Items {
		wg.Add(1)
		go ag.GetPodJSONStatus(pod.Status.PodIP, channel, &wg)
	}

	go func() {
		defer close(channel)
		wg.Wait()
	}()

	for ch := range channel {
		jsonOuts = append(jsonOuts, ch)
	}

	return jsonOuts, nil
}

// GetPlainStatusByNode func
func (ag *Aggregator) GetPlainStatusByNode() (string, error) {
	nodes, err := ag.GetJSONStatusForAll()
	if err != nil {
		return "", err
	}

	var results []TFStatus
	for _, node := range nodes {
		var result TFStatus
		err := json.Unmarshal([]byte(node), &result)
		if err != nil {
			log.Fatalf("Cannot unmarshal input data: %v", err)
			return "", err
		}
		results = append(results, result)
	}

	out := ""
	for _, node := range results {
		title := getTitle(node.PodName)
		out += fmt.Sprintf("%s\n", title)
		for _, grp := range node.Groups {
			out += fmt.Sprintf("\n== %s ==\n", grp.Name)
			out += ag.getPlainStatusOfGroup(grp)
		}
		out += fmt.Sprintf("\n")
	}

	return out, nil
}

func (ag *Aggregator) getPlainStatusOfGroup(group *TFGroup) string {
	data := [][]string{}
	for idS, svc := range group.Services {
		if idS > 0 {
		}
		data = append(data, []string{svc.Name, svc.Status})
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data)
	table.Render()

	return tableString.String()
}

// GetPlainStatusByGroup func
func (ag *Aggregator) GetPlainStatusByGroup() (string, error) {
	nodes, err := ag.GetJSONStatusForAll()
	if err != nil {
		return "", err
	}
	var results []TFStatus
	for _, node := range nodes {
		var result TFStatus
		err := json.Unmarshal([]byte(node), &result)
		if err != nil {
			log.Fatalf("Cannot unmarshal input data: %v", err)
		}
		results = append(results, result)
	}
	// Group - Service - node:status
	data := generateArrayForGroups(results)

	out := ""
	for grpName, svcs := range data {
		out += fmt.Sprintf("\n%s\n", getTitle(grpName))
		for svcName, nodes := range svcs {
			out += fmt.Sprintf("\n== %s ==\n", svcName)
			tbl := ag.getPlainStatusOfServices(nodes)
			out += fmt.Sprintf("%s\n", tbl)
		}
	}

	return out, nil
}

func generateArrayForGroups(statuses []TFStatus) map[string]map[string][][]string {

	data := make(map[string]map[string][][]string)
	for _, node := range statuses {
		nameN := node.PodName
		for _, grp := range node.Groups {
			nameG := grp.Name
			if _, ok := data[nameG]; !ok {
				data[nameG] = make(map[string][][]string)
			}
			for _, svc := range grp.Services {
				nameS := svc.Name
				data[nameG][nameS] = append(data[nameG][nameS], []string{nameN, svc.Status})
			}
		}
	}
	return data
}

func (ag *Aggregator) getPlainStatusOfServices(services [][]string) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(services)
	table.Render()

	return tableString.String()
}

func getTitle(title string) string {
	size := len(title)
	line := strings.Repeat("-", size)
	return fmt.Sprintf("%s\n%s\n%s", line, title, line)
}
