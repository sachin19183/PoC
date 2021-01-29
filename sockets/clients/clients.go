package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/hfgo/datafile"
)

func main() {
	fileName := "input.txt"
	clientReader, err := datafile.GetString(fileName)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.Dial("tcp", "127.0.0.1:8011")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	serverReader := bufio.NewReader(conn)
	log.Printf("Begin to start transmitting now \n")
	for _, value := range clientReader {
		fmt.Fprintf(conn, value+"\n") // Write data to server
		//_,err = conn.Write([]byte(value+'\n'))   //  if a byte format data is to be sent towards server
		servResp, err := serverReader.ReadString('\n') // Read response from the server
		if err == io.EOF {
			fmt.Println("Server closed the connection")
		} else if err == nil {
			fmt.Println(strings.TrimSpace(servResp))
		} else {
			log.Printf("Server error : %v\n", err)
		}

	}
	log.Printf("Ending now\n")

}
