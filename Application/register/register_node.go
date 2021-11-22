package main

import (
	"log"
	"net"
	"net/rpc"
	"strconv"

	"github.com/IlConteCvma/SDCC_Project/utility"
)

/*
	Main function for register node
*/

func main() {
	var connect_num int

	utils := new(utility.Utility)

	server := rpc.NewServer()
	//register method
	err := server.RegisterName("Utility", utils)
	if err != nil {
		log.Fatal("Format of service Utility is not correct: ", err)
	}

	port := utility.Server_port
	// listen for a connection
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("Error in listening:", err)
	}
	// Close the listener whenever we stop
	defer listener.Close()

	log.Printf("RPC server on port %d", port)

	go server.Accept(listener)

	//Wait connection
	for connect_num < utility.MAXCONNECTION {
		ch := <-utility.Connection
		if ch == true {
			connect_num++
		}
	}

	log.Printf("Max Number of Connection reached up")

	utility.Wg.Add(-utility.MAXCONNECTION)
	//send client a responce for max number of peer registered

	for {
		//TODO after registration this peer must be off ??
	}
}
