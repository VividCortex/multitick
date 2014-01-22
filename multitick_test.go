package multitick

// Copyright (c) 2013 VividCortex, Inc. All rights reserved.
// Please see the LICENSE file for applicable license terms.

import (
	"math/rand"
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

func samplingTestor(samplingInterval time.Duration, timesShouldBe []int, t *testing.T) {
	// Wait till next second boundary
	<-time.After(time.Second - time.Duration(time.Now().Nanosecond()))

	tick := NewTicker(time.Second, time.Millisecond*250)
	tick.mux.Lock()
	tick.randGenerator = rand.New(rand.NewSource(0))
	tick.mux.Unlock()
	c := tick.Subscribe()
	tick.Sample(samplingInterval)
	i := 0

	// Wait till first offset in the second
	start := <-time.After(time.Duration(250) * time.Millisecond -
		time.Duration(time.Now().Nanosecond()))

	for now := range c {
		if i > 5 {
			tick.Sample(time.Duration(0))
		}
		if i > 8 {
			tick.Stop()
			break
		}
		tenthsOfSeconds := int(now.Sub(start).Nanoseconds() / 100000000)
		t.Log("time (tenths of seconds)", tenthsOfSeconds)
		if tenthsOfSeconds != timesShouldBe[i] {
			t.Error("Expected", timesShouldBe[i], "but got", tenthsOfSeconds, "for tick.Sample(", samplingInterval, ")")
		}
		i++
	}
}

func TestSampling1(t *testing.T) {
	samplingTestor(time.Second*2, []int{20, 30, 60, 70, 90, 120, 140, 150, 160, 170}, t)
}

func TestSampling2(t *testing.T) {
	//NOTE: this the the same as the last test b/c of integer rounding.
	//This is expected
	samplingTestor(time.Second*2+time.Millisecond*999, []int{20, 30, 60, 70, 90, 120, 140, 150, 160, 170}, t)
}

func TestSampling3(t *testing.T) {
	samplingTestor(time.Second*3, []int{30, 40, 80, 120, 130, 160, 190, 200, 210, 220}, t)
}

func TestSampling4(t *testing.T) {
	samplingTestor(time.Second*4, []int{20, 50, 120, 150, 190, 240, 260, 270, 280}, t)
}

func TestSeed(t *testing.T) {
	tick := NewTicker(time.Second, time.Millisecond*250)
	tick.Seed()
}
