package network

import (
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type Nodes struct {
	Nodes []Node `json:"nodes"`
}

type Node struct {
	Ip string `json:"ip"`
	Peers string `json:"peers"`
}

func ReadNetworkConfig(nodesPtr *Nodes, directory string) {
	nodesJson, err := os.Open(directory)
	if err != nil {
		fmt.Println("Cannot read network configuration file")
		return
	}
	byteValue, _ := ioutil.ReadAll(nodesJson)
	nodesJson.Close()

	json.Unmarshal(byteValue, &nodesPtr)
}