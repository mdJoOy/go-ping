package main

import "time"

type Config struct {
	destination string
	ipv4        bool
	ipv6        bool
	count       int
	ttl         int
	size        int
	flood       bool
	deadLine    time.Duration
	interval    time.Duration
}
