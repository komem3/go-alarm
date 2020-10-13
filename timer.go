package alarm

import (
	"fmt"
	"time"
)

func Timer(d time.Duration) {
	start := time.Now()
	stopwatch := time.NewTicker(time.Second)
	ticker := time.NewTimer(d)
	defer stopwatch.Stop()
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			return
		case t := <-stopwatch.C:
			tsub := d - t.Sub(start)
			fmt.Printf("left %s\n", tsub.Round(time.Second))
		}
	}
}
