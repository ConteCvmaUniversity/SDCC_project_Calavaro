package main

import (
	"bufio"
	"container/list"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/IlConteCvma/SDCC_Project/utility"
)

/*
	This package implements the sequencer node, and all the function used for communication centralized between peers

*/
const (
	bufferSize int = 100
)

var (
	storeBuff = make(chan utility.Message, bufferSize)
	mutex     sync.Mutex
	seqNum    uint64
	verbose   bool
)

func main() {
	flag.BoolVar(&verbose, "v", false, "use this flag to get verbose info on messages")
	flag.Parse()

	if verbose {
		fmt.Println("VERBOSE FLAG ON")
	}

	//
	peers := list.New()
	utility.Registration(peers, utility.Sequencer, utility.Server_port)

	//debug
	for e := peers.Front(); e != nil; e = e.Next() {
		item := utility.Info(e.Value.(utility.Info))
		log.Printf("Type: %d, Address: %s:%s", item.Type, item.Address, item.Port)
	}

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(utility.Server_port))
	if err != nil {
		log.Fatal("net.Lister fail")
	}
	defer listener.Close()

	//start responce routine
	go response(peers)
	//handle peer connection
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept fail")
		}
		go handleConnection(connection)
	}

	//defer close(storeBuff)
}

func handleConnection(c net.Conn) {
	if verbose {
		log.Printf("Serving %s\n", c.RemoteAddr().String())
	}

	scanner := bufio.NewScanner(c)
	for {
		ok := scanner.Scan()

		if !ok {
			break
		}

		/*
			netData, err := bufio.NewReader(c).ReadString('\n')
			if err != nil {
				log.Fatal("Error ReadString")
			}
			fmt.Printf("%d , %s",tmp,netData)
			//stop connection
			temp := strings.TrimSpace(string(netData))
			if temp == utility.End_string {
				break
			}

		*/
		var msg utility.Message
		//critic section
		mutex.Lock()
		seqNum++
		msg.SeqNum = append(msg.SeqNum, seqNum)
		//msg.SeqNum[0] = seqNum
		mutex.Unlock()

		msg.Date = time.Now().Format("2006/01/02 15:04:05")
		msg.Text = scanner.Text()
		msg.Type = utility.SeqMsg
		msg.SendID = 0
		//send
		storeBuff <- msg

	}
	c.Close()
}

func response(peers *list.List) {
	//while channel storedBuff is open
	for msg := range storeBuff {
		//send responce to peer
		go sendResponse(peers, msg)
	}
}

func sendResponse(peers *list.List, msg utility.Message) {
	if verbose {
		fmt.Printf("[%d] Date: %s Text: %s\n", msg.SeqNum, msg.Date, msg.Text)
	}

	for element := peers.Front(); element != nil; element = element.Next() {
		item := utility.Info(element.Value.(utility.Info))

		//send to peer only
		if item.Type == utility.Peer {
			//open send channel
			peer_conn := item.Address + ":" + item.Port
			conn, err := net.Dial("tcp", peer_conn)
			if err != nil {
				log.Println("Send response error on Dial")
			}
			//send msg
			enc := gob.NewEncoder(conn)
			enc.Encode(msg)
		}
	}
}
