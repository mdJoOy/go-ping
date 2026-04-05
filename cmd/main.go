package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	ICMPProtocol   = 1
	ICMPv6Protocol = 58
)

func main() {
	config := &Config{}

	flag.IntVar(&config.count, "c", 0, "stop after <count> replies")
	flag.BoolVar(&config.ipv6, "6", false, "use Ipv6")
	flag.BoolVar(&config.flood, "f", false, "flood ping")
	//ttl, size, deadline, interval
	flag.IntVar(&config.ttl, "t", 64, "define time to live")
	flag.IntVar(&config.size, "s", 56, "use <size> as number of data bytes to send")
	flag.DurationVar(&config.interval, "i", time.Second, "seconds between sending each packet")

	// flag.IntVar(&config.interval, "i", timetime.Second, "seconds between sending each packet")
	var usages string = `Usage
  goping [options] <destination>

Options:
  <destination>      DNS name or IP address
  -c <count>         stop after <count> replies
  -d                 use SO_DEBUG socket option
  -f                 flood ping
  -i <interval>      seconds between sending each packet
  -s <size>          use <size> as number of data bytes to be sent
  -t <ttl>           define time to live
  -w <deadline>      reply wait <deadline> in seconds
  -W <timeout>       time to wait for response

  -6                 use IPv6
`
	flag.Usage = func() { fmt.Print(usages) }
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Print(usages)
		os.Exit(1)
	}

	config.destination = flag.Arg(0)

	ip, err := resolveHostIP(config.destination, config.ipv6)
	if err != nil {
		log.Fatal("couldnot resolve ip address")
	}

	ping(ip, config)
}
