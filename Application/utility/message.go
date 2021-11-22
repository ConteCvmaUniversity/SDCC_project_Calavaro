package utility

import (
	"container/list"
)

type MsgType int

const (
	SeqMsg         MsgType = 0
	ScalarClockMsg         = 1
	VectorClockMsg         = 2
	ScalarACK              = 3
)

type Message struct {
	Type   MsgType
	SendID int
	Date   string
	Text   string
	SeqNum []uint64 //used by sequencer or scalar clock value
}

func InsertInOrder(l *list.List, msg Message) *list.Element {
	//scan list element for the right position
	tmp := msg.SeqNum[0]
	//fmt.Println("MSG whit seq: "+ strconv.FormatUint(tmp,10))
	for e := l.Front(); e != nil; e = e.Next() {
		item := Message(e.Value.(Message))
		//fmt.Println("ITEM whit seq: "+ strconv.FormatUint(item.SeqNum,10))
		if tmp < item.SeqNum[0] {
			//found the next item
			//fmt.Println("IF CONDITION OK")
			return l.InsertBefore(msg, e)
		}
	}
	//fmt.Println("PUSHBACK")
	return l.PushBack(msg)
}
