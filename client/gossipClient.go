.package main

import (
	"bufio"
	"os"
	"fmt"
	"google.golang.org/grpc"
	log "github.com/sirupsen/logrus"
	pb "github.com/spockqin/leaderless-bft/proto"
	"context"
	"github.com/spockqin/leaderless-bft/tests/network"
	"github.com/spockqin/leaderless-bft/util"
	// "time"
	"strconv"
	"math/rand"
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
		return
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
		case msgStr == "":
			continue
		case msgStr == "get reqs":
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
		case msgStr == "latency same":
			testThroughPutSameConn(gossipers)
		case msgStr == "latency random":
			testThroughPutRandomConn(gossipers)
		default:
			hash := util.HashBytes([]byte(msgStr))
			exists, _ := mainClient.Poke(context.Background(), &pb.ReqId{Hash: hash})
			if exists.Status { // duplicate requests
				log.Info("Duplicate request")
				continue
			}
			mainClient.Push(context.Background(), &pb.ReqBody{Body: msg})
		}
	}
}

func testThroughPutSameConn(gossipers []string) {
	// fmt.Println("Timestamp right before first dialing: ",
	// 	time.Now().Format("2006-01-01 15:04:05 .000"))

	mainIp := gossipers[0]
	mainConn, err := grpc.Dial(mainIp, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Cannot establish TCP connection with gossiper")
		return
	}
	defer mainConn.Close()

	mainClient := pb.NewGossipClient(mainConn)

	for i := 0; i < 100; i++ {
		req := []byte(strconv.Itoa(i))
		mainClient.Push(context.Background(), &pb.ReqBody{Body: req})
	}
}

func testThroughPutRandomConn(gossipers []string) {
	// fmt.Println("Timestamp right before first dialing: ",
	// 	time.Now().Format("2006-01-01 15:04:05 .000"))

	i, n := 0, len(gossipers)
	for i < 100 {
		conn, err := grpc.Dial(gossipers[rand.Intn(n)], grpc.WithInsecure())
		if err != nil {
			continue
		}
		client := pb.NewGossipClient(conn)
		req := []byte(strconv.Itoa(i))
		client.Push(context.Background(), &pb.ReqBody{Body: req})
		conn.Close()
		i++
	}
}
