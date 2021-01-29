package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8011")
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Server started")
		io.WriteString(conn, "\nHello from the server\n")
		fmt.Fprintln(conn, "How is your day")
		fmt.Fprintf(conn, "%v", "well i hope")
		conn.Close()
		//message, _ := bufio.NewReader(conn).ReadString('\n')
		//fmt.Println("Received :", string(message))
	}
}
