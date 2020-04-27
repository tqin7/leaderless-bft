package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/spockqin/leaderless-bft/proto"
	tp "github.com/spockqin/leaderless-bft/types"
	"github.com/spockqin/leaderless-bft/tests/network"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"os"
	"strconv"
	"strings"
)

func main() {

	var nodes network.Nodes
	network.ReadNetworkConfig(&nodes, "../tests/network/config.json")

	var pbfters []string
	for _, node := range nodes.Nodes {
		pbfters = append(pbfters, node.Ip)
	}

	// create connection with main point of contact
	mainIp := pbfters[0]
	mainConn, err := grpc.Dial(mainIp, grpc.WithInsecure())
	if err != nil {
		log.WithFields(log.Fields{
			"node": mainIp,
			"error": err,
		}).Error("Cannot dial pbfter")
		panic(err)
	}
	defer mainConn.Close()

	fmt.Println("Established connection with ", mainIp)

	mainClient := pb.NewPbftClient(mainConn)

	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Message to send: ")
		msg, _ := reader.ReadString('\n')
		msg = msg[:(len(msg)-1)] // get rid of \n at the end
		msgStr := string(msg)
		switch {
		case msgStr == "get reqs":
			for _, ip := range pbfters {
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
		case msgStr == "throughput same":
			testThroughPutSameConn(pbfters)
		default:
			elements := strings.Split(msgStr,  " ")
			timeStamp, err := strconv.ParseInt(elements[2], 10, 64)
			if err != nil {
				fmt.Printf("err happens when converting timestamp to int64")
				panic(err)
			}
			req := &tp.PbftReq{
				ClientID:  	elements[0],
				Operation:  elements[1],
				Timestamp:  timeStamp,
				SequenceID: 0,
				MsgType:    "PbftReq",
			}

			reqBytes, err := json.Marshal(req)
			if err != nil {
				log.Error("client marshal request error!")
				panic(err)
			}

			_, err = mainClient.GetReq(context.Background(), &pb.ReqBody{Body:reqBytes})
			if err != nil {
				log.Error("Client sendReq error!")
				panic(err)
			}
		}
	}
}

func testThroughPutSameConn(pbfters []string) {
	// fmt.Println("Timestamp right before first dialing: ",
	// 	time.Now().Format("2006-01-01 15:04:05 .000"))

	mainIp := pbfters[0]
	mainConn, err := grpc.Dial(mainIp, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Cannot establish TCP connection with pbfter")
		return
	}
	defer mainConn.Close()

	mainClient := pb.NewPbftClient(mainConn)

	for i := 0; i < 100; i++ {
		req := []byte(fmt.Sprintf("70 msg%d %d", i, i)) // 70 is dummy clientID
		mainClient.GetReq(context.Background(), &pb.ReqBody{Body: req})
	}
}
