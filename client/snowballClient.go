package main

import (
	"bufio"
	"os"
	"fmt"
	"google.golang.org/grpc"
	log "github.com/sirupsen/logrus"
	pb "github.com/spockqin/leaderless-bft/proto"
	"context"
)

func main() {

	numOfNodes, start := 4, 2

	var snowers []string
	for i := start; i < start + numOfNodes; i++ {
		snowers = append(snowers, fmt.Sprintf("127.0.0.%d:7777", i))
	}

	// create connection with main point of contact
	mainIp := snowers[0]
	mainConn, err := grpc.Dial(mainIp, grpc.WithInsecure())
	if err != nil {
		log.WithFields(log.Fields{
			"node": mainIp,
			"error": err,
		}).Error("Cannot dial snower")
	}
	defer mainConn.Close()

	mainClient := pb.NewSnowballClient(mainConn)

	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Message to send: ")
		msg, _ := reader.ReadBytes('\n')
		msg = msg[:(len(msg)-1)] // get rid of \n at the end
		msgStr := string(msg)

		switch {
		case msgStr == "get reqs":
			for _, ip := range snowers {
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
		case msgStr == "get ordered":
			for _, ip := range snowers {
				conn, err := grpc.Dial(ip, grpc.WithInsecure())
				if err != nil {
					log.WithFields(log.Fields{
						"node": ip,
						"error": err,
					}).Error("Cannot dial node")
					continue
				}
				c := pb.NewSnowballClient(conn)
				requests, _ := c.GetOrderedReqs(context.Background(), &pb.Void{})
				fmt.Println(ip, "-", requests.Requests)
				conn.Close()
			}
		default:
			mainClient.SendReq(context.Background(), &pb.ReqBody{Body: msg})
		}
	}
}
