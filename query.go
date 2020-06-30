package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {

	req, err := http.NewRequest("GET", "http://18.162.47.230/queue/api/json", nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth("admin", "lihongzhi")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	found := strings.Contains(string(body), `There are no nodes with the label ‘work-node’`)
	fmt.Println(found)
}
