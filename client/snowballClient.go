package main

import (
	"bufio"
	"os"
	"fmt"
	"regexp"
	"google.golang.org/grpc"
	log "github.com/sirupsen/logrus"
	pb "github.com/spockqin/leaderless-bft/proto"
	"github.com/spockqin/leaderless-bft/util"
	"context"
	"github.com/spockqin/leaderless-bft/tests/network"
	"strconv"
)

func main() {

	var nodes network.Nodes
	network.ReadNetworkConfig(&nodes, "../tests/network/config.json")

	var snowers []string
	for _, node := range nodes.Nodes {
		snowers = append(snowers, node.Ip)
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
	numPattern, _ := regexp.Compile("get num (.+)")

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
		case numPattern.MatchString(msgStr):
			reqBytes := []byte(numPattern.FindStringSubmatch(msgStr)[1])
			seqNumRequest := pb.SeqNumMsg{SeqNum: 0, ReqHash: util.HashBytes(reqBytes)}
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
				seqNumRes, _ := c.GetVote(context.Background(), &seqNumRequest)
				fmt.Println(ip, "-", seqNumRes.SeqNum)
				conn.Close()
			}
		case msgStr == "throughput same":
			testThroughPutSameConn(snowers)
		default:
			mainClient.SendReq(context.Background(), &pb.ReqBody{Body: msg})
		}
	}
}

func testThroughPutSameConn(snowers []string) {
	// fmt.Println("Timestamp right before first dialing: ",
	// 	time.Now().Format("2006-01-01 15:04:05 .000"))

	mainIp := snowers[0]
	mainConn, err := grpc.Dial(mainIp, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Cannot establish TCP connection with snower")
		return
	}
	defer mainConn.Close()

	mainClient := pb.NewSnowballClient(mainConn)

	for i := 0; i < 100; i++ {
		req := []byte(strconv.Itoa(i))
		mainClient.SendReq(context.Background(), &pb.ReqBody{Body: req})
		fmt.Println("request ", i, " complete")
	}
}
