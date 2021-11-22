package main

import (
	"bufio"
	"container/list"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/IlConteCvma/SDCC_Project/utility"
)

const (
	bufferSize int = 100
)

var (
	msgSeqFile     *os.File
	msgScaFile     *os.File
	msgVecFile     *os.File
	filesMutex     [3]sync.Mutex
	mutex          sync.Mutex
	scalarMsgQueue *list.List
	ackCounter     map[string]int //key : msg.id-msg.seqNum
)
var (
	seqNum      uint64 = 1 //sequence number start from 1
	seqMsgChan         = make(chan *utility.Message, bufferSize)
	ackChan            = make(chan string, bufferSize)
	vectMsgChan        = make(chan *utility.Message, bufferSize)
)

func message_handler() {

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(utility.Client_port))
	if err != nil {
		log.Fatal("net.Lister fail")
	}
	defer listener.Close()

	//open file for save msg
	open_files()
	defer close_files()

	//function for
	go seqMsg_reodering()
	go ack_menager()
	go vectMsg_reordering()
	//scalar variable init
	scalarMsgQueue = list.New()
	ackCounter = make(map[string]int)

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept fail")
		}
		go handleConnection(connection)
	}
}

//Save message
func handleConnection(conn net.Conn) {
	// read msg and save on file
	defer conn.Close()
	msg := new(utility.Message)

	dec := gob.NewDecoder(conn)
	dec.Decode(msg)

	//msg management depends on msg.Type
	switch msg.Type {
	case utility.SeqMsg:
		//to msg reordering
		seqMsgChan <- msg

	case utility.ScalarClockMsg:
		//update clock
		tmp := msg.SeqNum
		updateClock(&scalarClock, tmp)
		incrementClock(&scalarClock, myId)
		//add in queue and send ack
		//e := scalarMsgQueue.PushBack(*msg)
		e := utility.InsertInOrder(scalarMsgQueue, *msg)
		tmpId := strconv.Itoa(msg.SendID) + "-" + strconv.FormatUint(msg.SeqNum[0], 10)

		go scalarMsgDemon(msg, e)
		go send_scalar_ack(tmpId)

	case utility.VectorClockMsg:
		if msg.SendID == myId {
			save_msg(msg)
		} else {
			//receive
			vectMsgChan <- msg
		}

	case utility.ScalarACK:
		text := msg.Text
		//fmt.Println("ACK FOR: " + text)
		ackChan <- text

	}

}

//-------------------------------------------------------------------------------------------------------------------------
/*
	FUNZIONI PER GESTIONE DI FILES
*/

func open_files() {
	var err error
	msgSeqFile, err = os.OpenFile(utility.Peer_msg_seq_file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("Impossible to open file")
	}
	msgScaFile, err = os.OpenFile(utility.Peer_msg_sca_file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("Impossible to open file")
	}
	msgVecFile, err = os.OpenFile(utility.Peer_msg_vec_file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("Impossible to open file")
	}
}

func close_files() {
	err := msgSeqFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = msgScaFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = msgVecFile.Close()
	if err != nil {
		log.Fatal(err)
	}

}

func save_msg(msg *utility.Message) {
	if utility.Launch_Test {
		testMsgChan <- msg
	}
	if verbose {
		fmt.Printf("Recv MSG--->Date: %s msg: %s seq: %d\n", msg.Date, msg.Text, msg.SeqNum)
	}

	//save on right file
	var (
		err error
		w   *bufio.Writer
	)
	switch msg.Type {
	case utility.SeqMsg:
		//save on msgSeqFile
		w = bufio.NewWriter(msgSeqFile)
		filesMutex[0].Lock()
		_, err = fmt.Fprintf(w, "Date: %s msg: %s seq: %d\n", msg.Date, msg.Text, msg.SeqNum)
		w.Flush()
		if err != nil {
			//TODO
			log.Fatal(err)
		}
		filesMutex[0].Unlock()

	case utility.ScalarClockMsg:
		//save on msgScaFile
		w = bufio.NewWriter(msgScaFile)
		filesMutex[1].Lock()
		_, err = fmt.Fprintf(w, "Date: %s msg: %s seq: %d\n", msg.Date, msg.Text, msg.SeqNum)
		w.Flush()
		if err != nil {
			//TODO
			log.Fatal(err)
		}
		filesMutex[1].Unlock()

	case utility.VectorClockMsg:
		//save on msgVecFile
		w = bufio.NewWriter(msgVecFile)
		filesMutex[2].Lock()
		_, err = fmt.Fprintf(w, "Date: %s msg: %s seq: %d\n", msg.Date, msg.Text, msg.SeqNum)
		w.Flush()
		if err != nil {
			//TODO
			log.Fatal(err)
		}
		filesMutex[2].Unlock()

	}
	if verbose {
		fmt.Println("Message saved")
	}

}

//-------------------------------------------------------------------------------------------------------------------------
/*
	FUNZIONI PER GESTIONE DI MESSAGGI DA SEQ
*/

func seqMsg_reodering() {
	for msg := range seqMsgChan {

		if msg.SeqNum[0] == seqNum {
			mutex.Lock()
			seqNum++
			mutex.Unlock()
			save_msg(msg)
		} else {
			//out of order msg go back
			seqMsgChan <- msg
		}

	}
}

//-------------------------------------------------------------------------------------------------------------------------
/*
	FUNZIONI PER GESTIONE DI MESSAGGI VECT
*/

func vectMsg_reordering() {
	for msg := range vectMsgChan {

		if condOne(*msg) && condTwo(*msg) {
			//can commit msg
			updateClock(&vectorClock, msg.SeqNum)
			save_msg(msg)
		} else {
			//Soluzione non efficiente ?? (con tanti peer?)
			vectMsgChan <- msg
		}
	}
}

//check if is the next expected msg
func condOne(msg utility.Message) bool {
	sendId := msg.SendID - 1 //because peer id start from 1

	if msg.SeqNum[sendId] == getValueClock(&vectorClock)[sendId]+1 {
		return true
	}
	return false
}

//check if
func condTwo(msg utility.Message) bool {
	msgId := msg.SendID

	for i := 0; i < len(allId); i++ {
		if !(msgId == allId[i]) {
			index := allId[i] - 1
			if msg.SeqNum[index] <= getValueClock(&vectorClock)[index] {
				break
			} else {
				return false
			}
		}
	}

	return true
}

//-------------------------------------------------------------------------------------------------------------------------

/*

	FUNZIONI PER IL RECEIVE DI SCALAR MSG
	p j consegna msg i all’applicazione se:
	1. msg i è in testa a queue j (e tutti gli ack relativi a msg i sono
	stati ricevuti da p j )
	2. per ogni processo p k c’è un messaggio msg k in queue j con
	timestamp maggiore di quello di msg i
*/

func scalarMsgDemon(msg *utility.Message, element *list.Element) {
	for !(checkConditionOne(*msg) && checkConditionTwo(*msg)) {
		//busy wait
		utility.Delay_ms(100)
	}
	//can commit msg
	save_msg(msg)
	//remove msg from queue and ack reference
	msgID := strconv.Itoa(msg.SendID) + "-" + strconv.FormatUint(msg.SeqNum[0], 10)
	delete(ackCounter, msgID)
	scalarMsgQueue.Remove(element)
}

func checkConditionOne(msg utility.Message) bool {
	//get head on queue
	tmp := utility.Message(scalarMsgQueue.Front().Value.(utility.Message))
	tmpId := strconv.Itoa(tmp.SendID) + "-" + strconv.FormatUint(tmp.SeqNum[0], 10)
	msgID := strconv.Itoa(msg.SendID) + "-" + strconv.FormatUint(msg.SeqNum[0], 10)
	//fmt.Println("tmpid:",tmpId,"msgid:",msgID)
	//fmt.Println(ackCounter[tmpId])
	mutex.Lock()
	ack := ackCounter[tmpId]
	mutex.Unlock()
	if tmpId == msgID && ack == utility.MAXPEERS {

		return true
	} else {
		return false
	}

}

func checkConditionTwo(msg utility.Message) bool {

	msgiQ := scalarMsgQueue.Front()
	msgi := utility.Message(msgiQ.Value.(utility.Message))

	//if msg not head return false
	if !checkEqualMSG(msg, msgi) {
		return false
	}

	for i := 0; i < len(allId); i++ {
		internal := false
		for e := msgiQ.Next(); e != nil; e = e.Next() {
			item := utility.Message(e.Value.(utility.Message))
			if item.SendID == allId[i] && item.SeqNum[0] > msg.SeqNum[0] {
				internal = true
				break
			}
		}
		if !internal {
			return false
		}
	}

	return true

}

func checkEqualMSG(msg1, msg2 utility.Message) bool {
	if msg1.SendID != msg2.SendID {
		return false
	}
	if msg1.SeqNum[0] != msg2.SeqNum[0] {
		return false
	}
	return true
}

func ack_menager() {
	for text := range ackChan {
		//fmt.Printf("Prima [%s]: %d\n",text,ackCounter[text])
		mutex.Lock()
		ackCounter[text] = ackCounter[text] + 1
		mutex.Unlock()
		//fmt.Printf("Dopo [%s]: %d\n " ,text,ackCounter[text])
	}

}

//----------------------------------------------------------
//Test function

func showMsg(_ ...string) error{
	var f []byte
	var err error
	switch communicationState {
	case Sequencer:
		f,err =os.ReadFile(utility.Peer_msg_seq_file)
	case Scalar:
		f,err =os.ReadFile(utility.Peer_msg_sca_file)
	case Vector:
		f,err =os.ReadFile(utility.Peer_msg_vec_file)

	}
	if err != nil {
		log.Println("Impossible to open file")
		return err
	}

	fmt.Print(string(f))


	return nil
}


func test(args ...string) error {

	fmt.Println("My ID: " + strconv.Itoa(myId))
	fmt.Printf("OTHERS : %d\tcap: %d\n", allId, cap(allId))
	//print actual state
	fmt.Println("ackCounter SIZE: " + strconv.Itoa(len(ackCounter)))
	//fmt.Println(ackCounter[args[0]])
	for e := scalarMsgQueue.Front(); e != nil; e = e.Next() {
		item := utility.Message(e.Value.(utility.Message))
		fmt.Println("SEND ID:" + strconv.Itoa(item.SendID))
		fmt.Println("ITEM SEQ:" + strconv.FormatUint(item.SeqNum[0], 10))
		fmt.Println("TEXT: " + item.Text)
		fmt.Println("----------------------------------------------------------")

	}

	return nil
}

func cleanScalarMSGQueue(){
	scalarMsgQueue = list.New()
}

//-------------------------------------------------------------------------------------------------------------------------
