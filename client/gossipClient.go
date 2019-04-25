package main

import (
	"bufio"
	"os"
	"fmt"
	"google.golang.org/grpc"
	log "github.com/sirupsen/logrus"
	pb "github.com/spockqin/leaderless-bft/proto"
	"context"
	"github.com/spockqin/leaderless-bft/tests/network"
)

func main() {

	var nodes network.Nodes
	network.ReadNetworkConfig(&nodes, "../tests/network/config.json")

	var gossipers []string
	for _, node := range nodes.Nodes {
		gossipers = append(gossipers, node.Ip)
	}

	// create connection with main point of contact
	mainIp := gossipers[0]
	mainConn, err := grpc.Dial(mainIp, grpc.WithInsecure())
	if err != nil {
		log.WithFields(log.Fields{
			"node": mainIp,
			"error": err,
		}).Error("Cannot dial gossiper")
	}
	defer mainConn.Close()
	fmt.Println("established main connection with ", mainIp)

	mainClient := pb.NewGossipClient(mainConn)

	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Message to send: ")
		msg, _ := reader.ReadBytes('\n')
		msg = msg[:(len(msg)-1)] // get rid of \n at the end
		msgStr := string(msg)

		switch {
		case msgStr == "get all":
			//requests, _ := client.GetAllRequests(context.Background(), &pb.Void{})
			//fmt.Println(mainIp, requests)
			for _, ip := range gossipers {
				conn, err := grpc.Dial(ip, grpc.WithInsecure())
				if err != nil {
					log.WithFields(log.Fields{
						"node": ip,
						"error": err,
					}).Error("Cannot dial node")
					continue
				}
				c := pb.NewGossipClient(conn)
				requests, _ := c.GetAllRequests(context.Background(), &pb.Void{})
				fmt.Println(ip, "-", requests.Requests)
				conn.Close()
			}
		default:
			mainClient.Push(context.Background(), &pb.ReqBody{Body: msg})
		}
	}
}
