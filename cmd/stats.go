package main

import "time"

type Stats struct {
	sent     int
	received int
	rtt      int
	maxRtt   int
	minRtt   int
	loss     int
	execTime time.Time
}

func (s *Stats) add() {
}
