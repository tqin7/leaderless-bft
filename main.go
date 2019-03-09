package main

import (
	"flag"
	. "github.com/spockqin/leaderless-bft/proto"
	"strings"
)

func main() {
	var ipAddr string
	var peersStr string

	flag.StringVar(&ipAddr,"ip", "127.0.0.2", "ip addr of node")
	flag.StringVar(&peersStr, "peers", "", "comma separated list of peers' ip addresses")

	flag.Parse()

	node := CreateNode(ipAddr)
	peers := strings.Split(peersStr, ",")

	for _, peer := range peers {
		node.AddPeer(peer)
	}

	node.NodeUp()
}
