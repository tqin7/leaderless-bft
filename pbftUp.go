package main

import (
	proto "github.com/spockqin/leaderless-bft/node"
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
		pbfter := proto.CreatePbfter(node.Ip, 1, node.Ip, allIps)
		peers := strings.Split(node.Peers, ",")
		for _, peer := range peers {
			pbfter.AddPeer(peer)
		}
		go pbfter.GetMsgFromGossip()
		if i == len(nodes.Nodes) - 1 {
			pbfter.PbfterUp()
		} else {
			go pbfter.PbfterUp()
		}
	}

	//reader := bufio.NewReader(os.Stdin)
	//text, _ := reader.ReadString('\n')
	//text = text[:(len(text)-1)]
	//index, _ := strconv.Atoi(string(text))
}