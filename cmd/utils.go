package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

// echo payload creator functin
func makePayload(size int) []byte {
	b := make([]byte, size)

	//trying to include the startTime in the paylaod
	startTime := time.Now().UnixNano()

	binary.BigEndian.PutUint64(b[:8], uint64(startTime))

	for i := 8; i < size; i++ {
		b[i] = byte(i & 0xff)
	}
	return b
}

// host resolver func
func resolveHostIP(host string, v6 bool) net.IP {
	addrs, err := net.LookupHost(host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldnot resolve host: %v\n", err)
	}
	for _, addr := range addrs {
		ipAddr := net.ParseIP(addr)
		if ipAddr.To4() != nil && !v6 {
			return ipAddr
		} else if ipAddr.To4() == nil && v6 {
			return ipAddr
		}

	}
	return net.ParseIP(addrs[0])
}

// print stats func
// --- google.com ping statistics ---
// 20 packets transmitted, 20 received, 0% packet loss, time 19030ms
// rtt min/avg/max/mdev = 4.120/10.022/31.023/7.191 ms
func printStats(s Stats, c Config) {
	fmt.Printf("\n--- %s ping statistics ---\n", c.destination)
	fmt.Printf("%d packets transmitted, %d received, %.2f%% packet loss, time ms\n", s.sent, s.received, s.loss())
	fmt.Printf("rtt min/avg/max/mdev = %.3f/%.3f/%.3f/%.3f ms\n", s.minRtt, s.maxRtt, s.avgRtt, s.stdDev())
}

// repeat func as helper func for the histrogram
func repeat(ch byte, rc int) string {
	b := make([]byte, rc)
	for i := range b {
		b[i] = ch
	}
	return string(b)
}
