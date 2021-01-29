package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/hfgo/datafile"
)

type CallbackResp struct {
	Msg   string
	Reqid int
}

func handleCallBackResp(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 30)
	var resp CallbackResp
	resp.Msg = "wait"
	resp.Reqid = 100
	for {
		respLen, err := conn.Read(buf)
		fmt.Println("Received as callback", string(buf))
		if err != nil {
			log.Fatal(err)
		} else if err == io.EOF {
			fmt.Println("Server has closed the connection")
			break
		} else {
			err = json.Unmarshal(buf[:respLen], &resp)
			if err != nil {
				fmt.Println("Error during unMarshall")
				log.Fatal(err)
			} else {
				//fmt.Println(respLen, ":", resp)
				resp.Msg = ""
				resp.Reqid = 0
			}
		}
	}

}
func handleCallback() {
	listner, err := net.Listen("tcp", ":8012")
	if err != nil {
		log.Panic(err)
	}
	defer listner.Close()

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Println(err)
		}
		go handleCallBackResp(conn)
	}

}
func connectServer() (net.Conn, error) {
	var conn net.Conn
	var err error
	conn, err = net.Dial("tcp", "127.0.0.1:8011")
	if err != nil {
		log.Fatalln(err)
	}
	return conn, err
}
func main() {
	fileName := "input.txt"
	clientReader, err := datafile.GetString(fileName)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := connectServer()
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	go handleCallback()
	serverReader := bufio.NewReader(conn)
	startTime := time.Now()
	log.Printf("Begin to start transmitting now \n")
	for _, value := range clientReader {
		log.Printf("Sending data\n")
		fmt.Fprintf(conn, value+"\n")
		log.Printf("Waiting for ACK \n")               // Write data to server
		servResp, err := serverReader.ReadString('\n') // Read response from the server
		log.Printf("Recvd ACK\n")
		if err == io.EOF {
			fmt.Println("EOF.Server closed the connection")

		} else if err == nil {
			fmt.Println(strings.TrimSpace(servResp))
		} else {
			log.Printf("Server error : %v\n", err)
			time.Sleep(2 * time.Second)
			conn, err = connectServer()
			continue
			//time.Sleep(4 * time.Second)

		}

	}
	time.Sleep(time.Duration(2) * time.Second)
	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Print("Total time taken for 500 req is ", diff.Seconds(), " seconds")
	log.Printf("Ending now\n")

}
