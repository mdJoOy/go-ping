package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	//	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

func ping(ip net.IP, cf *Config) {
	stat := &Stats{}
	var (
		network     string
		listener    string
		icmpMsgType icmp.Type
		protocol    int
	)

	//deciding using ipv4 or ipv6
	if ip.To4() != nil && !cf.ipv6 {
		network = "ip4:icmp"
		listener = "0.0.0.0"
		icmpMsgType = ipv4.ICMPTypeEcho
		protocol = ICMPProtocol
	} else if ip.To4() == nil && cf.ipv6 {
		network = "ip6:ipv6-icmp"
		listener = "::"
		icmpMsgType = ipv6.ICMPTypeEchoRequest
		protocol = ICMPv6Protocol
	}
	c, err := icmp.ListenPacket(network, listener)
	if err != nil {
		log.Fatalf("cannot open listener connection %w\n", err)
	}
	defer c.Close()

	//setting up ttl
	if !cf.ipv6 {
		if err := c.IPv4PacketConn().SetTTL(cf.ttl); err != nil {
			fmt.Errorf("couldn't set ttl: %w", err)
		}

	} else {
		if err := c.IPv6PacketConn().SetHopLimit(cf.ttl); err != nil {
			fmt.Errorf("couldn't set ttl: %w", err)
		}

	}
	//body for icmp echo request, which contains id, seq number or data in raw bytes
	id := os.Getpid() & 0xffff
	payload := makePayload(cf.size)

	dst := &net.IPAddr{IP: ip}

	fmt.Printf("PING %s (%s): %d bytes of data\n", cf.destination, ip, cf.size)

	//handling os termination like CTRL + C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	done := make(chan struct{})
	go func() {
		select {
		case <-sigCh:
		case <-done:
			return
		}
		printStats(*stat, *cf)
		if cf.histogram {
			stat.histogram()
		}
		os.Exit(0)

	}()

	for seq := 0; cf.count == 0 || seq < cf.count; seq++ {
		echoRequestBody := &icmp.Echo{ID: id, Seq: seq, Data: payload}
		icmpEchoMsg := icmp.Message{Type: icmpMsgType, Code: 0, Body: echoRequestBody}

		//marshaling the body
		icmpEchoMsgInBytes, err := icmpEchoMsg.Marshal(nil)
		if err != nil {
			log.Fatalf("could not marshal the icmp echo req, %w\n", err)
		}
		//setting deadline

		//	c.SetWriteDeadline(time.Now().Add(3 * time.Second))
		//send the echo msg
		timeNow := time.Now()
		if _, err := c.WriteTo(icmpEchoMsgInBytes, dst); err != nil {
			fmt.Println(err)
			fmt.Fprintf(os.Stderr, "write failed:%w\n", err)
			time.Sleep(cf.interval)
			continue
		}
		//reading the request
		rb := make([]byte, 1500)
		n, peer, err := c.ReadFrom(rb)
		if err != nil {
			log.Fatalf("couldnot read the echo reply message: %w\n", err)
			time.Sleep(cf.interval)
			continue
		}

		rm, err := icmp.ParseMessage(protocol, rb[:n])
		if err != nil {
			fmt.Println("error while parsing the icmp request msg")
			log.Fatal(err)
		}

		rtt := time.Since(timeNow).Seconds() * 1000
		time.Sleep(cf.interval)
		stat.add(rtt)
		rmBody := rm.Body.(*icmp.Echo)
		if rmBody.ID != id {
			continue
		}
		switch rm.Type {
		case ipv4.ICMPTypeEchoReply, ipv6.ICMPTypeEchoReply:
			if !cf.flood && !cf.quiet {
				fmt.Printf("%d bytes from %s: icmp_seq=%d ttl=%d time=%v ms\n", cf.size, peer, seq, cf.ttl, rtt)
			}
		case ipv4.ICMPTypeDestinationUnreachable, ipv6.ICMPTypeDestinationUnreachable:
			fmt.Printf("Destination: %s UNREACHABLE\n", peer)
		case ipv4.ICMPTypeTimeExceeded, ipv6.ICMPTypeTimeExceeded:
			fmt.Printf("Reply time out\n")
		}

		if cf.count != 0 && seq+1 >= cf.count {
			close(done)
			break
		}
	}
	stat.loss = float64((stat.sent - stat.received) / stat.sent * 100)
	printStats(*stat, *cf)
	if cf.histogram {
		stat.histogram()
	}

}
