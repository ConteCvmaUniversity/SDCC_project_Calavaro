package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/IlConteCvma/SDCC_Project/utility"
	"io"
	"log"
	"os"
	"strconv"
)

var (
	testMsgChan = make(chan *utility.Message, bufferSize) //channel for receive test msg
	results     = make(map[int]bool)
)

func setCommunicationState(state CommunicationType) {
	communicationState = state
}

func sendMsg_whitDelay(msg string, delay int) {
	if !(delay == 0) {
		utility.Delay_sec(utility.GetRandInt(delay))
	}
	err := sendMessages(msg)
	if err != nil {
		return
	}

}

/*
	numMsg are the number of msg that test must wait
*/
func waitMsg(numMsg int, file *os.File, respChan chan bool) {
	i := 0
	w := bufio.NewWriter(file)
	for msg := range testMsgChan {
		//save msg on file
		_, err := fmt.Fprintf(w, "Date: %s msg: %s seq: %d\n", msg.Date, msg.Text, msg.SeqNum)
		err = w.Flush()
		if err != nil {
			respChan <- false
			break
		} else {
			i++
			if i == numMsg {
				respChan <- true
				break
			}
		}
	}

}


func waitMsgWhitTimeout(timeout int, file *os.File,respChan chan bool) error{
	w := bufio.NewWriter(file)
	tout := make(chan bool)
	go utility.Timer(timeout, tout)
	select {
	case <-tout:

		respChan <- true

	case msg := <- testMsgChan:

		_, err := fmt.Fprintf(w, "Date: %s msg: %s seq: %d\n", msg.Date, msg.Text, msg.SeqNum)
		err = w.Flush()
		if err != nil {
			log.Printf("Error on waitMsgWhitTimeout")
			return errors.New("FILE write error")
		}
		respChan <- false
	}

	return nil
}

func executeTest(testId int, testFunc func(testId int) bool) {
	log.Printf("Starting test number %d\n", testId)
	res := testFunc(testId)
	results[testId] = res

	if verbose {
		if res {
			log.Printf("Test number %d PASS\n", testId)
		} else {
			log.Printf("Test number %d FAILED\n", testId)
		}
	}
}

func printResults() {
	var res string
	for k, v := range results {
		if v {
			res = "PASS"
		} else {
			res = "FAILED"
		}
		fmt.Printf("Test number %d %s\n", k, res)
	}
}

func compareTestFiles(testNum int) bool {

	for i := 0; i < len(allId)-1; i++ {
		targetFile := utility.Test_Dir + "test" + strconv.Itoa(testNum) + "_peer" + strconv.Itoa(allId[i])
		file1, err := os.OpenFile(targetFile, os.O_RDONLY, 0755)
		if err != nil {
			log.Fatal("Impossible to open file")
		}
		targetFile = utility.Test_Dir + "test" + strconv.Itoa(testNum) + "_peer" + strconv.Itoa(allId[i+1])
		file2, err := os.OpenFile(targetFile, os.O_RDONLY, 0755)
		if err != nil {
			log.Fatal("Impossible to open file")
		}

		ret := compareFile(file1, file2)
		file1.Close()
		file2.Close()

		if !ret {
			return false
		}
	}
	return true
}

func compareAllFileEmpty(testNum int) bool{
	for i := 0; i < len(allId); i++ {
		targetFile := utility.Test_Dir + "test" + strconv.Itoa(testNum) + "_peer" + strconv.Itoa(allId[i])
		file, err := os.OpenFile(targetFile, os.O_RDONLY, 0755)
		if err != nil {
			log.Fatal("Impossible to open file")
		}
		lines,_ := lineCounter(file)

		file.Close()
		if lines != 0{
			return false
		}

	}
	return true
}

func lineCounter(f *os.File) (int, error) {
	r := bufio.NewReader(f)
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func compareFile(file1, file2 *os.File) bool {
	wr1 := bufio.NewScanner(file1)
	wr2 := bufio.NewScanner(file2)

	ret := true
	for wr1.Scan() {
		wr2.Scan()
		if !(wr1.Text() == wr2.Text()) {
			ret = false
			break
		}
	}

	return ret
}
