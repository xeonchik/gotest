package misc

import (
	"log"
	"time"
)

var timerStart int64 = 0

func StartTimer() {
	timerStart = time.Now().UnixNano()
}

func LogTimer(str string) {
	resTime := (time.Now().UnixNano() - timerStart) / 1000
	log.Printf(str+": %d mcs", resTime)
}
