package main

import (
	"net"
	"bufio"
	"os"
	"fmt"
	"google.golang.org/grpc"
	"context"
	log "github.com/sirupsen/logrus"
	pb "github.com/spockqin/leaderless-bft/proto"
)

func main() {

	gossiperIp := "127.0.0.2:7777"
	conn, err := grpc.Dial(gossiperIp)
	if err != nil {
		log.WithFields(log.Fields{
			"node": gossiperIp,
		}).Error("Cannot dial node")
	}
	defer conn.Close()

	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Message to send: ")
		msg, _ := reader.ReadString('\n')

		// send message to node
		client := pb.NewGossipClient(conn)
		client.Push(context.Background(), &pb.ReqBody{Body: msg})
		//// listen for reply
		//message, _ := bufio.NewReader(conn).ReadString('\n')
		//fmt.Print("Message from server: " + message)
	}
}
