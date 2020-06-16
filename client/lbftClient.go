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
	"regexp"
	"github.com/spockqin/leaderless-bft/util"
)

func main() {

	var nodes network.Nodes
	network.ReadNetworkConfig(&nodes, "../tests/network/config.json")

	var lbfters []string
	for _, node := range nodes.Nodes {
		lbfters = append(lbfters, node.Ip)
	}

	// create connection with main point of contact
	mainIp := lbfters[0]
	mainConn, err := grpc.Dial(mainIp, grpc.WithInsecure())
	if err != nil {
		log.WithFields(log.Fields{
			"node": mainIp,
			"error": err,
		}).Error("Cannot dial pbfter")
		panic(err)
	} else {
		log.WithFields(log.Fields{
			"node": mainIp,
		}).Info("Connected to node")
	}
	defer mainConn.Close()

	mainClient := pb.NewLbftClient(mainConn)
	numPattern, _ := regexp.Compile("get num (.+)")
	// throughputPattern, _ := regexp.Compile("throughput (.+)")

	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Message to send: ")
		msg, _ := reader.ReadString('\n')
		msg = msg[:(len(msg)-1)] // get rid of \n at the end
		msgStr := string(msg)
		switch {
		case msgStr == "get reqs":
			for _, ip := range lbfters {
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
			for _, ip := range lbfters {
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
		case msgStr == "latency same":
			testLatencySameConn(lbfters)
		case msgStr == "throughput":
			testThroughPutSameConn(lbfters)
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

			_, err = mainClient.LSendReq(context.Background(), &pb.ReqBody{Body:reqBytes})
			if err != nil {
				log.Error("Client sendReq error!")
				panic(err)
			}
		}
	}
}

func constructLbftReqBytes(msgStr string) []byte {
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
	return reqBytes
}

func testLatencySameConn(lbfters []string) {
	// fmt.Println("Timestamp right before first dialing: ",
	// 	time.Now().Format("2006-01-01 15:04:05 .000"))

	mainIp := lbfters[0]
	mainConn, err := grpc.Dial(mainIp, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Cannot establish TCP connection with lbfter")
		return
	}
	defer mainConn.Close()

	mainClient := pb.NewLbftClient(mainConn)

	for i := 0; i < 3; i++ {
		msgStr := fmt.Sprintf("%d %d %d", i, i, i)
		reqBytes := constructLbftReqBytes(msgStr)
		_, err := mainClient.LSendReq(context.Background(), &pb.ReqBody{Body: reqBytes})
		if err != nil {
			log.Error("Client sendReq error!")
			panic(err)
		} else {
			fmt.Printf("req %d done ", i)
		}
	}
}

func testThroughPutSameConn(lbfters []string) {
	mainIp := lbfters[0]
	mainConn, err := grpc.Dial(mainIp, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Cannot establish TCP connection with lbfter")
		return
	}
	defer mainConn.Close()

	mainClient := pb.NewLbftClient(mainConn)

	i := 0

	for {
		msgStr := fmt.Sprintf("%d %d %d", i, i, i)
		reqBytes := constructLbftReqBytes(msgStr)
		_, err := mainClient.LSendReq(context.Background(), &pb.ReqBody{Body: reqBytes})
		if err != nil {
			log.Error("Client sendReq error!")
			panic(err)
		} else {
			fmt.Printf("req %d done ", i)
		}
		i += 1
	}
}

