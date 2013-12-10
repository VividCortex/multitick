// Package multitick broadcasts a time.Ticker to multiple receivers,
// all of which receive the same values, aligned to an offset.
package multitick

import (
	"sync"
	"time"
)

// Ticker is a broadcaster for time.Time tick events.
type Ticker struct {
	mux    sync.Mutex
	chans  []chan time.Time
	ticker *time.Ticker
}

// NewTicker creates and starts a new Ticker, whose ticks are sent to
// subscribers at the specified interval, delayed until the specified offset
// after the interval (there will be some imprecision). The first parameter is
// equivalent to "d" in time.NewTicker(d). If the offset is negative, ticking
// is not aligned to an offset, and begins immediately.
func NewTicker(interval, offset time.Duration) *Ticker {
	t := &Ticker{}
	if offset >= 0 { // Sleep till the specified offset past the interval.
		var (
			now   = time.Now().UnixNano()
			intNs = interval.Nanoseconds()
			offNs = offset.Nanoseconds()
			sleep time.Duration
		)
		elapsed := now % intNs // since the last interval boundary
		if elapsed <= offNs {  // we haven't passed the offset, sleep till it
			sleep = time.Duration(offNs - elapsed)
		} else { // we passed the offset, sleep till next interval+offset
			sleep = time.Duration((intNs - elapsed) + offNs)
		}
		time.Sleep(sleep)
	}
	t.ticker = time.NewTicker(interval)
	go t.tick()
	return t
}

// Subscribe returns a channel to which ticks will be delivered. Ticks that
// can't be delivered to the channel, because it is not ready to receive, are
// discarded.
func (t *Ticker) Subscribe() <-chan time.Time {
	t.mux.Lock()
	defer t.mux.Unlock()
	c := make(chan time.Time)
	t.chans = append(t.chans, c)
	return c
}

// Stop stops the ticker. As in time.Ticker, it does not close channels.
func (t *Ticker) Stop() {
	t.ticker.Stop()
}

// This could be inlined as an anonymous function, but I think it's easier to
// read stacktraces with real function names in them.
func (t *Ticker) tick() {
	for tick := range t.ticker.C {
		t.mux.Lock()
		for i := range t.chans {
			select {
			case t.chans[i] <- tick:
			default:
			}
		}
		t.mux.Unlock()
	}
}
