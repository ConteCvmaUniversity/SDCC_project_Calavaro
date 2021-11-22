package main

import (
	"github.com/IlConteCvma/SDCC_Project/utility"
)

var (
	scalarClock utility.ScalarClock
	vectorClock utility.VectorClock
)

func startClocks() {
	scalarClock.Start()
	vectorClock.Start()
}

func incrementClock(clock utility.Clock, id int) {
	clock.Increment(id - 1)
}

func updateClock(clock utility.Clock, timestamp []uint64) {
	clock.Update(timestamp)
}

func getValueClock(clock utility.Clock) []uint64 {
	return clock.GetValue()
}
