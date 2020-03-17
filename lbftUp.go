package main

import (
	"github.com/spockqin/leaderless-bft/node"
	. "github.com/spockqin/leaderless-bft/tests/network"
	"strings"
)

func main() {
	var nodes Nodes
	ReadNetworkConfig(&nodes, "./tests/network/config.json")
	var allIps []string
	for _, node := range nodes.Nodes {
		allIps = append(allIps, node.Ip)
	}
	for i, node := range nodes.Nodes {
		lbfter := proto.CreateLbfter(node.Ip, 1, node.Ip, allIps)
		peers := strings.Split(node.Peers, ",")
		for _, peer := range peers {
			lbfter.AddPeer(peer)
		}
		go lbfter.GetMsgFromGossip()
		go lbfter.ResolveMsg()
		if i == len(nodes.Nodes) - 1 {
			lbfter.LbfterUp()
		} else {
			go lbfter.LbfterUp()
		}
	}
}