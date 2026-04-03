package main

import (
	"flag"
	"fmt"
	"log"
)

const (
	ICMPProtocol   = 1
	ICMPv6Protocol = 58
)

func main() {
	config := &Config{}

	flag.IntVar(&config.count, "count", 0, "stop after <count> replies")
	flag.Parse()
	config.destination = flag.Arg(0)
	fmt.Println(config.destination)
	ip, err := resolveHostIP(config.destination)
	fmt.Println(ip)
	if err != nil {
		log.Fatal("couldnot resolve ip address")
	}
	ping(ip)
}
