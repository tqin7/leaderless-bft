package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/spockqin/leaderless-bft/node"
	"strings"
)

type Nodes struct {
	Nodes []Node `json:"nodes"`
}

type Node struct {
	Ip string `json:"ip"`
	Peers string `json:"peers"`
}

func main() {
	nodesJson, err := os.Open("./tests/network/config.json")
	if err != nil {
		fmt.Println("Cannot read network configuration file")
		return
	}
	defer nodesJson.Close()

	byteValue, _ := ioutil.ReadAll(nodesJson)

	var nodes Nodes
	json.Unmarshal(byteValue, &nodes)

	for i, node := range nodes.Nodes {
		gossiper := proto.CreateGossiper(node.Ip)
		peers := strings.Split(node.Peers, ",")
		for _, peer := range peers {
			gossiper.AddPeer(peer)
		}

		if i == len(nodes.Nodes) - 1 {
			proto.GossiperUp(gossiper)
		} else {
			go proto.GossiperUp(gossiper)
		}
	}


	//var ipAddr string
	//var peersStr string
	//
	//flag.StringVar(&ipAddr,"ip", "127.0.0.2", "ip addr of node")
	//flag.StringVar(&peersStr, "peers", "", "comma separated list of peers' ip addresses")
	//
	//flag.Parse()
	//
	//gossiper := CreateGossiper(ipAddr)
	//peers := strings.Split(peersStr, ",")
	//
	//for _, peer := range peers {
	//	gossiper.AddPeer(peer)
	//}
	//
	//GossiperUp(gossiper)
}
