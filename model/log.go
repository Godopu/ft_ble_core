package model

import (
	"fmt"
	"time"
)

type TimeLog struct {
	ExpType         string
	start           time.Time
	endComputing    time.Time
	end             time.Time
	NetworkingDelay uint64 `json:"networking_delay"`
	ComputingDelay  uint64 `json:"computing_delay"`
	TotalDelay      uint64 `json:"total_delay"`
}

func (tl *TimeLog) Start() {
	tl.start = time.Now()
}

func (t1 *TimeLog) SetComputingTime(t uint64) {
	t1.ComputingDelay = t
}

func (tl *TimeLog) EndComputing() {
	tl.endComputing = time.Now()
}

func (tl *TimeLog) End() {
	tl.end = time.Now()
	tl.NetworkingDelay =
		uint64(tl.endComputing.Sub(tl.start).Milliseconds()) - tl.ComputingDelay + uint64(tl.end.Sub(tl.endComputing).Milliseconds())

	tl.TotalDelay = tl.NetworkingDelay + tl.ComputingDelay
}

func (tl *TimeLog) String() string {
	// return "Hello"
	return fmt.Sprintf("%s] - (Networking, Compugting, Total) = (%d ms, %d ms, %d ms)",
		tl.ExpType, tl.NetworkingDelay, tl.ComputingDelay, tl.TotalDelay)
}

var LogStack []*TimeLog

func init() {
	LogStack = make([]*TimeLog, 0, 100)
}

func AddTimeLog(tl *TimeLog) {
	LogStack = append(LogStack, tl)
}

func ClearLogStack() {
	LogStack = LogStack[:0]
}
