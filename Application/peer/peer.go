package main

import (
	"container/list"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/IlConteCvma/SDCC_Project/utility"
)

/*
	Main function for register node
*/
var (
	peers   *list.List
	myId    int
	allId   []int
	verbose bool
)

func main() {
	flag.BoolVar(&verbose, "v", false, "use this flag to get verbose info on messages")
	flag.Parse()

	if verbose {
		fmt.Println("VERBOSE FLAG ON")
	}
	//phase registration
	peers = list.New()
	utility.Registration(peers, utility.Peer, utility.Client_port)
	//debugging output remove ?

	for e := peers.Front(); e != nil; e = e.Next() {
		item := utility.Info(e.Value.(utility.Info))
		log.Printf("Type: %d, Address: %s:%s", item.Type, item.Address, item.Port)
	}

	//get myId
	setMyID()
	//start clock
	startClocks()

	//open listen channel for messages
	//service on port 2345
	go message_handler()

	if utility.Launch_Test {
		startTests()
		os.Exit(2) //test complete
	}
	//open menu
	open_menu()
}

func setMyID() {

	for e := peers.Front(); e != nil; e = e.Next() {
		item := utility.Info(e.Value.(utility.Info))
		if item.Address == utility.GetLocalIP() {
			myId = item.ID
			allId = append(allId, item.ID)
		} else {
			allId = append(allId, item.ID)
		}
	}
	allId = allId[1:len(allId)] //remove sequencer
}
