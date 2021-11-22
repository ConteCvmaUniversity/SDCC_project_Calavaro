package utility

import (
	"math/rand"
	"sync"
	"time"
)

const maxDelay = 5

type Clock interface {
	Start()           // Start not concurrent operation
	Increment(id int) //id: my node identifier, do not care if ScalarClock
	Update(timestamp []uint64)
	GetValue() []uint64
	//AtomicIncAndGet ??
}

type ScalarClock struct {
	counter uint64
	mutex   sync.Mutex
}

func (clock *ScalarClock) Start() {
	clock.counter = 0
}

func (clock *ScalarClock) Increment(_id int) {
	clock.mutex.Lock()
	defer clock.mutex.Unlock()
	clock.counter++
}
func (clock *ScalarClock) Update(timestamp []uint64) {
	clock.mutex.Lock()
	defer clock.mutex.Unlock()
	clock.counter = MaxOf(clock.counter, timestamp[0])
}
func (clock *ScalarClock) GetValue() []uint64 {
	ret := make([]uint64, 1)
	clock.mutex.Lock()
	defer clock.mutex.Unlock()
	ret[0] = clock.counter
	return ret
}

type VectorClock struct {
	counter []uint64
	mutex   sync.Mutex
}

func (clock *VectorClock) Start() {
	clock.counter = make([]uint64, MAXPEERS, MAXPEERS)
}
func (clock *VectorClock) Increment(id int) {
	clock.mutex.Lock()
	defer clock.mutex.Unlock()
	clock.counter[id] = clock.counter[id] + 1

}

func (clock *VectorClock) Update(timestamp []uint64) {
	clock.mutex.Lock()
	defer clock.mutex.Unlock()
	for k := 0; k < len(clock.counter); k++ {
		clock.counter[k] = MaxOf(clock.counter[k], timestamp[k])
	}
}
func (clock *VectorClock) GetValue() []uint64 {
	clock.mutex.Lock()
	defer clock.mutex.Unlock()
	return clock.counter
}

func Delay() {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(maxDelay)
	time.Sleep(time.Duration(n) * time.Second)
}

func Delay_ms(maxTime int) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(maxTime)
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func Delay_sec(exactTime int) {
	time.Sleep(time.Duration(exactTime) * time.Second)
}

func Timer(timeout int, outChan chan bool)  {
	Delay_sec(timeout)
	outChan <- true
}
