package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func main() {
	deregisterAgent("18.163.184.77:80", "jenkins-cli.jar", "i-0743e48f00bf94046")
}

func deregisterAgent(url, jenkinsCli, agent string) {
	cmd := exec.Command("java", "-jar", jenkinsCli, "-s", "http://"+url, "delete-node", agent)
	err := cmd.Run()
	log.Printf("Command finished with error: %v", err)
}

func queryAgent(url, agent string) bool {
	// http://18.163.184.77/computer/i-0912b2c691a1f56e5/api/json
	requrl := fmt.Sprintf("http://%v/computer/%v/api/json", url, agent)
	req, err := http.NewRequest("GET", requrl, nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(body))
	idle := strings.Contains(string(body), `"idle":true`)
	return idle
}

func queryQueue(url string) bool {
	requrl := fmt.Sprintf("http://%v/queue/api/json", url)
	req, err := http.NewRequest("GET", requrl, nil)
	if err != nil {
		panic(err)
	}
	// req.SetBasicAuth("admin", "lihongzhi")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(body))
	found := strings.Contains(string(body), `There are no nodes with the label ‘work-node’`)
	return found
}
