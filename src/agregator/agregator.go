package agregator

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

// Agregator struct
type Agregator struct {
	client  client.Client
	context context.Context
}

// New - constructor for Agregator
func New() *Agregator {
	ag := Agregator{}
	ag.setClient()

	ag.context = context.TODO()

	return &ag
}

func (ag *Agregator) setClient() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err, "")
		os.Exit(1)
	}
	// clientset, err := apiextensionsclient.NewForConfig(cfg)
	cli, err := client.New(cfg, client.Options{})
	if err != nil {
		log.Fatal(err, "")
		os.Exit(1)
	}
	ag.client = cli
}

// GetPodList func
func (ag *Agregator) GetPodList() *corev1.PodList {
	podList := &corev1.PodList{}
	err := ag.client.List(ag.context, podList, client.MatchingLabels{"tungstenfabric": "status"})
	if err != nil {
		fmt.Printf("ERR client output %v\n", err)
		os.Exit(1)
	}

	return podList
}

// GetPodJSONStatus func
func (ag *Agregator) GetPodJSONStatus(podIP string, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("INFO request %s\n", podIP)

	nodePort := os.Getenv("NODE_PORT")
	if len(nodePort) == 0 {
		nodePort = "80"
	}

	resp, err := http.Get("http://" + podIP + ":" + nodePort + "/json")
	if err != nil {
		fmt.Printf("ERR request %s, %v\n", podIP, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ERR get content %s, %v\n", podIP, err)
	}
	fmt.Printf("INFO json from %s has been recieved\n", podIP)
	ch <- string(body)
}

// GetJSONStatusForAll func
func (ag *Agregator) GetJSONStatusForAll() []string {
	podList := ag.GetPodList()
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

	return jsonOuts
}

// GetPlainStatusByNode func
func (ag *Agregator) GetPlainStatusByNode() string {
	nodes := ag.GetJSONStatusForAll()

	var results []TFStatus
	for _, node := range nodes {
		var result TFStatus
		err := json.Unmarshal([]byte(node), &result)
		if err != nil {
			fmt.Printf("ERR Unmarshal: %v", err)
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

	fmt.Printf("%s", out)
	return out
}

func (ag *Agregator) getPlainStatusOfGroup(group *TFGroup) string {
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
func (ag *Agregator) GetPlainStatusByGroup() string {
	nodes := ag.GetJSONStatusForAll()

	var results []TFStatus
	for _, node := range nodes {
		var result TFStatus
		err := json.Unmarshal([]byte(node), &result)
		if err != nil {
			fmt.Printf("ERR Unmarshal: %v", err)
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

	fmt.Printf("%s", out)
	return out
}

func generateArrayForGroups(statuses []TFStatus) map[string]map[string][][]string {

	data := make(map[string]map[string][][]string)
	for _, node := range statuses {
		nameN := node.PodName
		fmt.Printf("Node: %s\n", nameN)

		for _, grp := range node.Groups {
			nameG := grp.Name
			fmt.Printf("Group: %s\n", nameG)
			if _, ok := data[nameG]; !ok {
				data[nameG] = make(map[string][][]string)
			}
			for _, svc := range grp.Services {
				nameS := svc.Name
				fmt.Printf("Service: %s\n", nameS)
				data[nameG][nameS] = append(data[nameG][nameS], []string{nameN, svc.Status})
			}
		}
	}
	return data
}

func (ag *Agregator) getPlainStatusOfServices(services [][]string) string {
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
