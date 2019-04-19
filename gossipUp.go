package main

import (
	"github.com/spockqin/leaderless-bft/node"
	"strings"
	. "github.com/spockqin/leaderless-bft/tests/network"
)

func main() {
	var nodes Nodes
	ReadNetworkConfig(&nodes, "./tests/network/config.json")

	for i, node := range nodes.Nodes {
		gossiper := proto.CreateGossiper(node.Ip)
		//gossiper := proto.CreateGossiper(node.Ip, allIps)
		peers := strings.Split(node.Peers, ",")
		for _, peer := range peers {
			gossiper.AddPeer(peer)
		}

		if i == len(nodes.Nodes) - 1 {
			gossiper.GossiperUp()
		} else {
			go gossiper.GossiperUp()
		}
	}

	// COMMAND LINE FLAG
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
