package main

import (
	"github.com/spockqin/leaderless-bft/node"
	"strings"
	. "github.com/spockqin/leaderless-bft/tests/network"
)

func main() {
	var nodes Nodes
	ReadNetworkConfig(&nodes, "./tests/network/config.json")

	var allIps []string
	for _, node := range nodes.Nodes {
		allIps = append(allIps, node.Ip)
	}

	for i, node := range nodes.Nodes {
		snower := proto.CreateSnower(node.Ip, allIps)
		peers := strings.Split(node.Peers, ",")
		for _, peer := range peers {
			snower.AddPeer(peer)
		}

		if i == len(nodes.Nodes) - 1 {
			snower.SnowerUp()
		} else {
			go snower.SnowerUp()
		}
	}
}
