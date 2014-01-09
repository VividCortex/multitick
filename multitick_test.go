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

func samplingTestor(samplingInterval time.Duration, seed int64, timesShouldBe []int, t *testing.T) {
	tick := NewTicker(time.Second, time.Millisecond*250)
	c := tick.Subscribe()
	tick.Sample(samplingInterval, seed)
	i := 0
	start := time.Now()
	for now := range c {
		if i > 5 {
			tick.Sample(time.Duration(0), 0)
		}
		if i > 8 {
			tick.Stop()
			break
		}
		tenthsOfSeconds := int(now.Sub(start).Nanoseconds() / 100000000)
		t.Log("time (tenths or seconds)", tenthsOfSeconds)
		if tenthsOfSeconds != timesShouldBe[i] {
			t.Error("Expected", timesShouldBe[i], "but got", tenthsOfSeconds, "for tick.Sample(", samplingInterval, ",", seed, ")")
		}
		i++
	}
}

func TestSampling1(t *testing.T) {
	samplingTestor(time.Second*2, 0, []int{20, 30, 60, 70, 90, 120, 140, 150, 160, 170}, t)
}

func TestSampling2(t *testing.T) {
	//NOTE: this the the same as the last test b/c of integer rounding.
	//This is expected
	samplingTestor(time.Second*2+time.Millisecond*999, 0, []int{20, 30, 60, 70, 90, 120, 140, 150, 160, 170}, t)
}

func TestSampling3(t *testing.T) {
	samplingTestor(time.Second*3, 0, []int{30, 40, 80, 120, 130, 160, 190, 200, 210, 220}, t)
}

func TestSampling4(t *testing.T) {
	samplingTestor(time.Second*4, 1, []int{30, 80, 100, 160, 180, 210, 270, 280, 290, 300}, t)
}
