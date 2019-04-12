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

func createGrpcConn(ip string) *grpc.ClientConn {
	conn, err := grpc.Dial(ip, grpc.WithInsecure())
	if err != nil {
		log.WithFields(log.Fields{
			"node": ip,
			"error": err,
		}).Error("Cannot dial node")
	}
	return conn
}

func main() {

	numOfNodes, start := 100, 2

	var gossipers []string
	for i := start; i < start + numOfNodes; i++ {
		gossipers = append(gossipers, fmt.Sprintf("127.0.0.%d:7777", i))
	}

	// create connection with main point of contact
	mainIp := gossipers[0]
	mainConn, err := grpc.Dial(mainIp, grpc.WithInsecure())
	if err != nil {
		log.WithFields(log.Fields{
			"node": mainIp,
			"error": err,
		}).Error("Cannot dial node")
	}
	defer mainConn.Close()

	mainClient := pb.NewGossipClient(mainConn)

	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Message to send: ")
		msg, _ := reader.ReadBytes('\n')
		msg = msg[:(len(msg)-1)] // get rid of \n at the end

		switch {
		case string(msg) == "get all":
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
