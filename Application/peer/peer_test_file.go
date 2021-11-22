package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/IlConteCvma/SDCC_Project/utility"
)

/*
	Tests are developed into main package because the topology of peers need a setup, this is the fastest way to do test
*/

const (
	maxDelay = 5
)

func startTests() {

	utility.Delay_sec(3)
	//Run tests
	executeTest(1, testSequencer)
	executeTest(2, testSeqMultiMSG)
	
	/*First run comment this test on second one*/
	executeTest(3, testScalar)
	executeTest(5, testVector)
	
	/*Second run comment this test on first one*/
	//executeTest(4, testScaMultiMSG)
	//executeTest(6, testVecMultiMSG)
	
	//executeTest(7, testSlideEX)


	if myId == 1 {
		utility.Delay_sec(5)
		printResults()
		if utility.Clean_Test_Dir {
			err := os.RemoveAll(utility.Test_Dir)

			if err != nil {
				log.Println(err)
				//return
			}
		}

	}


	//infinite loop only if Clean_test_dir false to see results
	for !utility.Clean_Test_Dir {

	}

}

/*
	Testing if sequencer
*/
func testSequencer(testId int) bool {
	const numMsg = 10
	msgs := [numMsg]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	//only peer 1 do send
	if myId == 1 {
		setCommunicationState(Sequencer)
		for _, s := range msgs {
			sendMsg_whitDelay(s, 2)
		}
	}

	var respChan = make(chan bool)
	targetFile := utility.Test_Dir + "test" + strconv.Itoa(testId) + "_peer" + strconv.Itoa(myId)
	file, err := os.OpenFile(targetFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("Impossible to open file")
	}

	go waitMsg(numMsg, file, respChan)
	resp := <-respChan
	err = file.Close()
	if err != nil || !resp {
		log.Fatal(err)
	}

	log.Printf("Waiting test messages\n")
	utility.Delay_sec(2) //busy time to complete all file arrival
	//start compare test file
	log.Printf("Compare results\n")
	resp = compareTestFiles(testId)

	return resp
}

/*
	Similar to testSequencer but all peer send msg to sequencer
	The order to sequencer cannot be predicted, but the order back to peers must be the same
*/

func testSeqMultiMSG(testId int) bool {
	const numMsg = 3 * utility.MAXPEERS //3 msg per peer
	msgs := [3]string{"1", "2", "3"}

	setCommunicationState(Sequencer)
	for _, s := range msgs {
		sendMsg_whitDelay(s+"peer"+strconv.Itoa(myId), 2)
	}

	var respChan = make(chan bool)
	targetFile := utility.Test_Dir + "test" + strconv.Itoa(testId) + "_peer" + strconv.Itoa(myId)
	file, err := os.OpenFile(targetFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("Impossible to open file")
	}

	go waitMsg(numMsg, file, respChan)
	resp := <-respChan
	err = file.Close()
	if err != nil || !resp {
		log.Fatal(err)
	}

	log.Printf("Waiting test messages\n")
	utility.Delay_sec(2) //busy time to complete all file arrival
	//start compare test file
	log.Printf("Compare results\n")
	resp = compareTestFiles(testId)

	return resp
}


func testScalar(testId int) bool{
	const numMsg = 10
	msgs := [numMsg]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	//only peer 1 do send
	if myId == 1 {
		setCommunicationState(Scalar)
		for _, s := range msgs {
			sendMsg_whitDelay(s, 2)
		}
	}

	var respChan = make(chan bool)
	targetFile := utility.Test_Dir + "test" + strconv.Itoa(testId) + "_peer" + strconv.Itoa(myId)
	file, err := os.OpenFile(targetFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("Impossible to open file")
	}
	go func() {
		err = waitMsgWhitTimeout(numMsg, file, respChan)
		if err != nil {
			log.Fatal(err)
		}
	}()

	resp := <-respChan
	err = file.Close()
	if err != nil || !resp {
		log.Println(err)
		return false
	}
	log.Printf("Waiting test messages\n")
	utility.Delay_sec(2) //busy time to complete all file arrival
	log.Printf("Compare results\n")
	resp = compareAllFileEmpty(testId)
	return resp
}


/*
	Testing scalar send by all peer
	3 message send by peer but expected 2 back
*/
func testScaMultiMSG(testId int) bool{
	const numMsg = 3 * utility.MAXPEERS //3 msg per peer
	msgs := [3]string{"1", "2", "3"}

	setCommunicationState(Scalar)
	for _, s := range msgs {
		sendMsg_whitDelay(s+"peer"+strconv.Itoa(myId), 2)
	}

	var respChan = make(chan bool)
	targetFile := utility.Test_Dir + "test" + strconv.Itoa(testId) + "_peer" + strconv.Itoa(myId)
	file, err := os.OpenFile(targetFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("Impossible to open file")
	}

	go waitMsg(2, file, respChan)

	resp := <-respChan
	err = file.Close()

	if err != nil || !resp {
		log.Println(err)
		return false
	}
	if myId == 1{
		//show some queue info
		test("")
	}

	log.Printf("Waiting test messages\n")
	utility.Delay_sec(2) //busy time to complete all file arrival
	log.Printf("Compare results\n")
	resp = compareTestFiles(testId)
	return resp
}

func testVector(testId int) bool {
	const numMsg = 10
	msgs := [numMsg]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	//only peer 1 do send
	if myId == 1 {
		setCommunicationState(Vector)
		for _, s := range msgs {
			sendMsg_whitDelay(s, 2)
		}
	}

	var respChan = make(chan bool)
	targetFile := utility.Test_Dir + "test" + strconv.Itoa(testId) + "_peer" + strconv.Itoa(myId)
	file, err := os.OpenFile(targetFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("Impossible to open file")
	}

	go waitMsg(numMsg, file, respChan)
	resp := <-respChan
	err = file.Close()
	if err != nil || !resp {
		log.Fatal(err)
	}

	log.Printf("Waiting test messages\n")
	utility.Delay_sec(2) //busy time to complete all file arrival
	//start compare test file
	log.Printf("Compare results\n")
	resp = compareTestFiles(testId)

	return resp
}

func testVecMultiMSG(testId int) bool {
	const numMsg = 3 * utility.MAXPEERS //3 msg per peer
	msgs := [3]string{"1", "2", "3"}

	setCommunicationState(Vector)
	for _, s := range msgs {
		sendMsg_whitDelay(s+"peer"+strconv.Itoa(myId), 2)
	}

	var respChan = make(chan bool)
	targetFile := utility.Test_Dir + "test" + strconv.Itoa(testId) + "_peer" + strconv.Itoa(myId)
	file, err := os.OpenFile(targetFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("Impossible to open file")
	}

	go waitMsg(numMsg, file, respChan)
	resp := <-respChan
	err = file.Close()
	if err != nil || !resp {
		log.Fatal(err)
	}

	log.Printf("Waiting test messages\n")
	utility.Delay_sec(2) //busy time to complete all file arrival
	//start compare test file
	log.Printf("Compare results\n")
	resp = compareTestFiles(testId)

	return resp
}

func testSlideEX(testId int) bool{
	var res bool
	setCommunicationState(Vector)

	if myId == 1 {
		sendMsg_whitDelay("Ma",2)
		sendMsg_whitDelay("Mb",2) //do not care result
	}

	targetFile := utility.Test_Dir + "test" + strconv.Itoa(testId) + "_peer" + strconv.Itoa(myId)
	file, err := os.OpenFile(targetFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("Impossible to open file")
	}
	w := bufio.NewWriter(file)
	i := 0
	for msg := range testMsgChan {
		_, err := fmt.Fprintf(w, "Date: %s msg: %s seq: %d\n", msg.Date, msg.Text, msg.SeqNum)
		err = w.Flush()
		if err != nil {
			res = false
			break
		}else {
			if msg.Text == "Ma" && myId == 2 {
				sendMsg_whitDelay("Mc",2)
				sendMsg_whitDelay("Md",2)
			}
			i++
			if i == 4 {
				res = true
				break
			}
		}
	}

	if !res {
		log.Fatal("Some error")
	}
	file.Close()
	//start compare test file
	log.Printf("Compare results\n")
	utility.Delay_sec(2) //busy time to complete all file arrival
	res = compareTestFiles(testId)

	return res

}