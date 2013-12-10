package multitick

import (
	"testing"
	"time"
)

func Test1(t *testing.T) {
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
