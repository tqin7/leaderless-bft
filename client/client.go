package main

import (
	"net"
	"bufio"
	"os"
	"fmt"
)

func main() {

	nodeIp := "127.0.0.2:7777"
	conn, _ := net.Dial("tcp", nodeIp)
	defer conn.Close()

	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Message to send: ")
		msg, _ := reader.ReadString('\n')

		//m := proto.Message{proto.StoreKey(msg)}
		//b, err := json.Marshal(m)

		// send message to node
		fmt.Fprintf(conn, msg + "\n")
		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}
}
