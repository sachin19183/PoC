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
		go reqDistributor(conn, clid)
	}
}

func reqDistributor(conn net.Conn, clid int) error {
	defer conn.Close()
	log.Printf("Begin to process client %d now", clid)

	var thSwitch bool
	thSwitch = true
	callbackConn, err := connectCallback(clid)
	if err != nil {
		log.Fatal(err)
		fmt.Println("Failed in connect callback")
		return err
	}
	defer callbackConn.Close()
	aChanel := make(chan string)
	bChanel := make(chan string)
	go handleA(aChanel, clid, callbackConn)
	go handleB(bChanel, clid, callbackConn)
	var ack string
	clientReader := bufio.NewReader(conn)
	for {
		clientReq, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			clientReq = strings.TrimSpace(clientReq)
			ack = fmt.Sprintf("ACK Received for %s .waiting for callback", clientReq)
			n, err1 := fmt.Fprintf(conn, ack+"\n")
			//n, err1 := io.WriteString(conn, ack)
			if err1 != nil {
				fmt.Println("Error while sending ACK ", err)
			} else {
				fmt.Printf("ACK sent with %d bytes.Processing requests for Client: %d", n, clid)
				if thSwitch == true {
					aChanel <- clientReq
					thSwitch = false
				} else {
					bChanel <- clientReq
					thSwitch = true
				}

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

func handleA(aChan chan string, clid int, callbackConn net.Conn) error {
	var resp CallbackResp
	fmt.Println("Processing requests for Client: ODD")

	for {
		clientReq := <-aChan
		fmt.Printf("ODD Client %d has sent %s\n", clid, clientReq)
		time.Sleep(time.Duration(350) * time.Millisecond)
		id, _ := strconv.Atoi(clientReq)

		resp.Msg = "Received"
		resp.Reqid = id
		buf, _ := json.Marshal(resp)
		fmt.Println(string(buf))
		sent, _ := callbackConn.Write(buf)
		fmt.Println("callback  ", id, "with", sent, " bytes sent back towards ODD")
		//sendCallback(callbackConn, id)
	}
}

func handleB(bChan chan string, clid int, callbackConn net.Conn) error {
	var resp CallbackResp
	fmt.Println("Processing requests for Client EVEN: ")

	for {
		clientReq := <-bChan
		fmt.Printf("EVEN client %d has sent %s\n", clid, clientReq)
		time.Sleep(time.Duration(350) * time.Millisecond)
		id, _ := strconv.Atoi(clientReq)
		resp.Msg = "Received"
		resp.Reqid = id
		buf, _ := json.Marshal(resp)
		fmt.Println(string(buf))
		sent, _ := callbackConn.Write(buf)
		fmt.Println("callback  ", id, " with ", sent, " bytes sent back towards EVEN")
		//sendCallback(callbackConn, id)
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
func sendCallback(conn net.Conn, id int) {

	var resp CallbackResp
	resp.Msg = "Received"
	resp.Reqid = id
	buf, _ := json.Marshal(resp)
	fmt.Println(string(buf))
	sent, _ := conn.Write(buf)
	fmt.Println("sent ", sent, "bytes towards calback")

}
