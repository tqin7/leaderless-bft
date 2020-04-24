package main

import (
	"os"
	"log"
	"time"
	"strconv"
	"fmt"
	"context"
	"io"
	pb "github.com/spockqin/leaderless-bft/proto"
	"google.golang.org/grpc"
)

func main() {
	serverFile := os.Args[1]

	MAIN_IP := "127.0.0.1:30000"
	mainConn, err := grpc.Dial(MAIN_IP, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Cannot establish TCP connection")
		return
	}
	defer mainConn.Close()

	startTimeLog := fmt.Sprintf("Testing Timestamp: %d", time.Now().Unix())
	appendToFile("log.txt", startTimeLog)

	if serverFile == "gossipUp" {
		mainClient := pb.NewGossipClient(mainConn)
		for i := 0; i < 100; i++ {
			req := []byte(strconv.Itoa(i))
			mainClient.Push(context.Background(), &pb.ReqBody{Body: req})
			// time.Sleep(0.5 * time.Second)
		}
	} else if serverFile == "snowballUp" {
		mainClient := pb.NewSnowballClient(mainConn)
		for i := 0; i < 100; i++ {
			req := []byte(strconv.Itoa(i))
			mainClient.SendReq(context.Background(), &pb.ReqBody{Body: req})
			// time.Sleep(0.5 * time.Second)
		}
	} else if serverFile == "pbftUp" {
		mainClient := pb.NewPbftClient(mainConn)
		for i := 0; i < 100; i++ {
			req := []byte(strconv.Itoa(i))
			mainClient.GetReq(context.Background(), &pb.ReqBody{Body: req}) //TODO: modify msg to 1 msg1 1
			// time.Sleep(0.5 * time.Second)
		}
	} else if serverFile == "lbftUp" {
		mainClient := pb.NewLbftClient(mainConn)
			for i := 0; i < 100; i++ {
			req := []byte(strconv.Itoa(i))
			mainClient.LSendReq(context.Background(), &pb.ReqBody{Body: req}) //TODO: modify msg to 1 msg1 1
			// time.Sleep(0.5 * time.Second)
		}
	} else {
		log.Fatal("client type not supported")
		return
	}
}

func writeToFile(filename string, data string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = io.WriteString(file, data)
    if err != nil {
        return err
    }
    return file.Sync()
}

func appendToFile(filename string, data string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
	    panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(data); err != nil {
	    panic(err)
	}
}