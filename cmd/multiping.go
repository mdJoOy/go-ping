package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
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
	//one raw socket for every hosts
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}

	//making a slice of targets
	targets := make([]*Target, len(hosts))

	for i, host := range hosts {
		targets[i] = &Target{
			host:    host,
			ip:      resolveHostIP(host, cfg.ipv6),
			id:      i + 1,
			stats:   &Stats{},
			replies: make(chan Reply, 10),
		}
	}
	//creating sender go rutine for each target
	for _, target := range targets {
		go sender(*conn, target, cfg)
	}

	rawReplies := make(chan Reply, 50)
	go reader(*conn, rawReplies)
	go display(rawReplies, targets)

	return nil
}

func sender(c icmp.PacketConn, t *Target, cfg Config) {
	ticker := time.NewTicker(cfg.interval)
	defer ticker.Stop()
	seq := 0
	for range ticker.C {

		echoMessage := icmp.Message{Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{ID: t.id, Seq: seq, Data: makePayload(cfg.size)},
		}
		rawMessage, err := echoMessage.Marshal(nil)
		if err != nil {
			log.Fatal(err)
		}
		dst := &net.IPAddr{IP: t.ip}
		if _, err := c.WriteTo(rawMessage, dst); err != nil {
			log.Fatal(err)
		}
		if seq+1 == cfg.count {
			return
		}
		seq++
	}
}

func reader(c icmp.PacketConn, out chan<- Reply) {
	buf := make([]byte, 1500)
	for {
		n, peer, err := c.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}
		parseReply, err := icmp.ParseMessage(ICMPProtocol, buf[:n])
		if err != nil {
			log.Fatal(err)
		}
		if parseReply.Type == ipv4.ICMPTypeEchoReply {
			icmpEchoReplyBody := parseReply.Body.(*icmp.Echo)
			//
			returnedPayload := icmpEchoReplyBody.Data
			sentTimeNano := int64(binary.BigEndian.Uint64(returnedPayload[:8]))
			rtt := time.Since(time.Unix(0, sentTimeNano))
			out <- Reply{from: net.IP(peer.String()), id: icmpEchoReplyBody.ID, seq: icmpEchoReplyBody.Seq, rtt: rtt.Seconds() * 1000}
		}

	}
}

func display(replies chan Reply, targets []*Target) {
	for reply := range replies {
		for _, t := range targets {
			if t.id == reply.id {
				fmt.Printf("received reply from %s, seq=%d, id=%d, rtt=%vms\n", t.host, reply.seq, reply.id, reply.rtt)
			}
		}

	}
}
