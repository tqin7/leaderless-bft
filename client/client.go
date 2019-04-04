package main

import (
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
	conn, err := grpc.Dial(gossiperIp, grpc.WithInsecure())
	if err != nil {
		log.WithFields(log.Fields{
			"node": gossiperIp,
			"error": err,
		}).Error("Cannot dial node")
	}
	defer conn.Close()

	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Message to send: ")
		msg, _ := reader.ReadBytes('\n')

		// send message to node
		client := pb.NewGossipClient(conn)
		//exists, _ := client.Poke(context.Background(), &pb.ReqId{Hash: msg})
		//fmt.Println(exists.Status)
		client.Push(context.Background(), &pb.ReqBody{Body: msg})


		//// listen for reply
		//message, _ := bufio.NewReader(conn).ReadString('\n')
		//fmt.Print("Message from server: " + message)
	}
}
