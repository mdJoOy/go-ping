package main

import (
	"fmt"
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
func (s *Stats) loss() float64 {
	return float64(s.sent-s.received) / float64(s.sent) * 100
}

func (s *Stats) histogram() {
	if len(s.rtts) == 0 {
		fmt.Printf("Not enough rtts to draw a histrogram\n")
		return
	}
	min := s.minRtt
	max := s.maxRtt
	if min == max {
		fmt.Printf("ALL rtts in the same range so, no histrogram necessary\n")
		return
	}

	bucket := 10
	bucketWidth := (max - min) / float64(bucket)

	counts := make([]int, bucket)
	for _, v := range s.rtts {
		idx := int((v - min) / bucketWidth)

		if idx >= bucket {
			idx = bucket - 1
		}
		counts[idx]++
	}
	maxCount := 0
	for _, v := range counts {
		if v > maxCount {
			maxCount = v
		}
	}
	fmt.Printf("\n--- Latency Histogram ---\n")
	for i, v := range counts {
		low := min + float64(i)*bucketWidth
		high := low + bucketWidth
		bar := float64(v) / float64(maxCount) * 40
		fmt.Printf("%6.2f - %6.2f ms | %-*s %d\n", low, high, 40, repeat('#', int(bar)), v)
	}
}
