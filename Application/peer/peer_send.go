package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/IlConteCvma/SDCC_Project/utility"
)

type CommunicationType int

const (
	Sequencer CommunicationType = 0
	Scalar                      = 1
	Vector                      = 2
)

var communicationState CommunicationType

func sendMessages(args ...string) error {
	var function func(msgs []string)
	if len(args) < 1 {
		fmt.Println("Wrong usage of send expected: send arg1 arg2 ... ")
		return nil
	}

	switch communicationState {
	case Sequencer:
		function = send_to_seq

	case Scalar:
		function = send_to_scalar

	case Vector:
		function = send_to_vector
	}

	//chiamata sincrona se aggiungo go la rendo asincrona
	function(args)

	return nil
}

//msg get from stdin
func send_to_seq(msgs []string) {
	//connect to sequencer
	seq_conn := utility.Seq_addr + ":" + strconv.Itoa(utility.Server_port)
	conn, err := net.Dial("tcp", seq_conn)
	defer conn.Close()
	if err != nil {
		log.Println("Send to sequencer error on Dial")
	}
	for _, element := range msgs {
		//fmt.Println(element)
		//fmt.Fprintf(conn, element+"\n")
		if verbose {
			fmt.Printf("SENDING %s to Sequencer\n", element)
		}
		_, err := conn.Write([]byte(element + "\n"))
		if err != nil {
			log.Println("Error writing to stream.")
			break
		}
	}
	//fmt.Fprintf(conn, utility.End_string+"\n")

}

func send_to_scalar(msgs []string) {

	for _, text := range msgs {
		//increment local clock
		incrementClock(&scalarClock, myId)

		//prepare msg to send
		var msg utility.Message
		msg.Type = utility.ScalarClockMsg
		msg.SeqNum = append(msg.SeqNum, getValueClock(&scalarClock)[0])
		msg.Date = time.Now().Format("2006/01/02 15:04:05")
		msg.Text = text
		msg.SendID = myId

		send_to_peer(msg)
	}

}

func send_to_vector(msgs []string) {

	for _, text := range msgs {
		incrementClock(&vectorClock, myId)
		//prepare msg to send
		var msg utility.Message
		msg.Type = utility.VectorClockMsg
		msg.SeqNum = append(msg.SeqNum, getValueClock(&vectorClock)...)
		msg.Date = time.Now().Format("2006/01/02 15:04:05")
		msg.Text = text
		msg.SendID = myId
		send_to_peer(msg)

	}
}

func send_scalar_ack(text string) {
	//prepare msg to send
	var msg utility.Message
	msg.Type = utility.ScalarACK
	msg.Text = text
	msg.SendID = myId

	send_to_peer(msg)

}

func send_to_peer(msg utility.Message) {
	if verbose {
		if msg.Type == utility.ScalarClockMsg{
			fmt.Printf("SEND MSG--->Date: %s msg: %s seq: %d\n", msg.Date, msg.Text, msg.SeqNum)
		}else {
			fmt.Printf("ACK: %s",msg.Text)
		}

	}

	for e := peers.Front(); e != nil; e = e.Next() {
		dest := utility.Info(e.Value.(utility.Info))
		//only peer are destination of msgs
		if dest.Type == utility.Peer {

			//open connection whit peer
			peer_conn := dest.Address + ":" + dest.Port
			conn, err := net.Dial("tcp", peer_conn)
			defer conn.Close()
			if err != nil {
				log.Println("Send response error on Dial")
			}
			enc := gob.NewEncoder(conn)
			enc.Encode(msg)
		}
	}
}
