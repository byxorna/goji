package main

import (
	"log"
	"time"
)

// coalesce events on a channel within a window, and then call fn()
// http://blog.gopheracademy.com/advent-2013/day-24-channel-buffering-patterns/
//TODO this doesnt quite work right... It fires fn() immediately on first event in ch
func coalesceEvents(ch <-chan string, window time.Duration, fn func()) {
	ticker := time.NewTimer(0)
	var timerCh <-chan time.Time
	i := 0
	for {
		select {
		case e := <-ch:
			// count how many events we coalesce, for fun
			i = i + 1
			log.Printf("Coalescing event %s. (%d events so far)\n", e, i)
			log.Printf("%s\n", timerCh)
			if timerCh == nil {
				ticker.Reset(window)
				timerCh = ticker.C
			}
		case <-timerCh:
			log.Printf("Coalesced %d events\n", i)
			fn()
			i = 0
			timerCh = nil
		}
	}
}
