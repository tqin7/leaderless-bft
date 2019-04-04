package main

import (
	"flag"
	"strings"
	. "github.com/spockqin/leaderless-bft/node"
)

func main() {
	var ipAddr string
	var peersStr string

	flag.StringVar(&ipAddr,"ip", "127.0.0.2", "ip addr of node")
	flag.StringVar(&peersStr, "peers", "", "comma separated list of peers' ip addresses")

	flag.Parse()

	gossiper := CreateGossiper(ipAddr)
	peers := strings.Split(peersStr, ",")

	for _, peer := range peers {
		gossiper.AddPeer(peer)
	}

	go GossiperUp(gossiper)
	gossiper.StartGossip()
}
