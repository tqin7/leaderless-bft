package main

import (
	"github.com/spockqin/leaderless-bft/node"
	"strconv"
	"strings"
	. "github.com/spockqin/leaderless-bft/tests/network"
)

func main() {
	var nodes Nodes
	ReadNetworkConfig(&nodes, "./tests/network/config.json")

	nodeTable := make(map[string]string)
	for i, node := range nodes.Nodes {
		nodeTable[strconv.Itoa(i)] = node.Ip
	}

	for i, node := range nodes.Nodes {
		pbfter := proto.CreatePbfter(strconv.Itoa(i), node.Ip, nodeTable)
		peers := strings.Split(node.Peers, ",")
		for _, peer := range peers {
			pbfter.AddPeer(peer)
		}

		if i == len(nodes.Nodes) - 1 {
			pbfter.SnowerUp()
		} else {
			go pbfter.SnowerUp()
		}
	}
}
