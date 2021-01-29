package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type CallbackResp struct {
	Msg   string
	Reqid int
}

func main() {
	clid := 0
	listner, err := net.Listen("tcp", ":8011")
	if err != nil {
		log.Panic(err)
	}
	defer listner.Close()

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Println(err)
		}
		clid++
		go handle(conn, clid)
	}
}

func handle(conn net.Conn, clid int) error {
	defer conn.Close()
	callbackConn, err := connectCallback(clid)
	if err != nil {
		log.Fatal(err)
		fmt.Println("Failed in connect callback")
		return err
	}
	defer callbackConn.Close()
	var ack string
	log.Printf("Begin to process client %d now", clid)
	clientReader := bufio.NewReader(conn)
	for {
		clientReq, err := clientReader.ReadString('\n')
		switch err {
		case nil:
			ack = fmt.Sprintf("ACK Recvd for %s ", clientReq)
			io.WriteString(conn, ack)
			fmt.Println("ACK sent.Processing requests for Client: ", clid)
			clientReq = strings.TrimSpace(clientReq)
			time.Sleep(time.Duration(350) * time.Millisecond)
			id, _ := strconv.Atoi(clientReq)
			err1 := sendCallback(callbackConn, id)
			if err1 != nil {
				callbackConn, err2 := connectCallback(clid)
				if err2 != nil {
					log.Fatal(err2)
					fmt.Println("Failed in connect callback")
					return err
				}
				defer callbackConn.Close()
				err1 = sendCallback(callbackConn, id)
			}

		case io.EOF:
			log.Printf("In EOF. Received close from %d\n", clid)
			return err
		default:
			log.Printf("In Default. Received %v from client\n", err)
			return err
		}

	}

}
func connectCallback(clid int) (net.Conn, error) {
	var callbackConn net.Conn
	var err error
	if clid == 1 {
		callbackConn, err = net.Dial("tcp", "127.0.0.1:8012")
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		callbackConn, err = net.Dial("tcp", "127.0.0.1:8013")
		if err != nil {
			log.Fatalln(err)
		}
	}
	return callbackConn, err
}
func sendCallback(conn net.Conn, id int) error {

	var resp CallbackResp
	resp.Msg = "Received"
	resp.Reqid = id
	buf, _ := json.Marshal(resp)
	fmt.Println(string(buf))
	sent, err := conn.Write(buf)
	if err != nil {
		log.Fatal(err)
		fmt.Println("error during callback sent")
		return err
	}
	fmt.Println("callback ", id, " sent back :", sent, "bytes")
	return nil
	//conn.Close()

}
