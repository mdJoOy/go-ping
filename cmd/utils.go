package main

import (
	"fmt"
	"net"
)

// echo body creator functin
func makeBody(size int) []byte {
	b := make([]byte, size)
	for i := range b {
		b[i] = byte(i & 0xff)
	}
	return b
}

func resolveHostIP(host string, v6 bool) (net.IP, error) {
	addrs, err := net.LookupHost(host)
	if err != nil {
		fmt.Errorf("couldnot resolve host: %w\n", err)
	}
	for _, addr := range addrs {
		ipAddr := net.ParseIP(addr)
		if ipAddr.To4() != nil && !v6 {
			return ipAddr, nil
		} else if ipAddr.To4() == nil && v6 {
			return ipAddr, nil
		}

	}
	return net.ParseIP(addrs[0]), nil
}
