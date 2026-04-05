package main

import "time"

type Config struct {
	destination string
	ipv6        bool
	count       int
	ttl         int
	size        int
	flood       bool
	quit        bool
	interval    time.Duration
	timeout     time.Duration
	deadLine    time.Duration
}
