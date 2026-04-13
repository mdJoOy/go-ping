package main

import (
	"net"
)

type Reply struct {
	from net.IP
	id   int
	seq  int
	rtt  float64
}

type Target struct {
	host    string
	ip      net.IP
	id      int
	stats   *Stats
	replies chan Reply
}

func pingMultiple(cfg Config, hosts []string) error {
	return nil
}
