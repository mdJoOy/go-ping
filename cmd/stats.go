package main

import (
	"math"
	"time"
)

type Stats struct {
	sent      int
	received  int
	rtts      []float64
	maxRtt    float64
	minRtt    float64
	avgRtt    float64
	loss      float64
	startTime time.Time
}

// avg of rtt
func (s *Stats) avg() {
	var sumRtt float64
	for _, rtt := range s.rtts {
		sumRtt += rtt
	}
	s.avgRtt = sumRtt / float64(s.received)
}

// adding
func (s *Stats) add(rtt float64) {
	s.sent++
	s.received++
	s.rtts = append(s.rtts, rtt)

	if s.minRtt == 0 || rtt < s.minRtt {
		s.minRtt = rtt
	}
	if rtt > s.maxRtt {
		s.maxRtt = rtt
	}

	s.avg()
}

func (s *Stats) stdDev() float64 {
	var variance float64

	for _, rtt := range s.rtts {
		diff := rtt - s.avgRtt
		variance += diff * diff
	}
	return math.Sqrt(variance / float64(s.received))
}
