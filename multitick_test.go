package multitick

// Copyright (c) 2013 VividCortex, Inc. All rights reserved.
// Please see the LICENSE file for applicable license terms.

import (
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	tick := NewTicker(time.Second, time.Millisecond*250)
	c := tick.Subscribe()
	chans := make([]<-chan time.Time, 0)
	i := 0
	for now := range c {
		chans = append(chans, tick.Subscribe())
		if i > 5 {
			tick.Stop()
			break
		}
		i++
		t.Log(now, now.Nanosecond())
		if now.Nanosecond() < 247000000 || now.Nanosecond() > 253000000 {
			t.Errorf("%v isn't within 3ms of 250ms offset", now)
		}
	}
}

func TestRandomTicker(t *testing.T) {
  var seed int64 = 0
	tick := NewRandomTicker(time.Second,seed)
	c := tick.Subscribe()
	chans := make([]<-chan time.Time, 0)
	i := 0
  oldNow := time.Now()
	for now := range c {
		chans = append(chans, tick.Subscribe())
		if i > 5 {
			tick.Stop()
			break
		}
		i++
    delTime := time.Since(oldNow).Nanoseconds()
    if delTime < 999700000 || delTime > 1003000000 {
      t.Error("Time difference should always be 1 second for this test. Was",delTime)
    } else {
      oldNow = now
    }
		t.Log(now, now.Nanosecond())
    //TODO 744ms is a magic number. It is the offset that gets specified when
    //the seed is 0
		if now.Nanosecond() < 741000000 || now.Nanosecond() > 747000000 {
			t.Error(now,"isn't within 3ms of 747ms offset")
		}
	}
}
