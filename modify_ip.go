package main

import (
	"golang.org/x/net/icmp"
	"golang.org/x/net/internal/iana"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"os"
	"flag"
	"fmt"
)

const targetIP4address = "8.8.8.8"

func main() {
	dstipPtr := flag.String("dst_ip", "8.8.8.8", "Destination IP Address")
	srcipPtr := flag.String("src_ip", "100.64.0.10", "Src IP Address")
	
	flag.Parse()
	fmt.Printf(" Dst IP: %s Src IP: %s \n ", *dstipPtr, *srcipPtr)
	
	//cm := newControlMessage()
	conn, err := icmp.ListenPacket("ip4:icmp","0.0.0.0")
	if err != nil {
		log.Fatalf("listen error, %s", err)
	}
	defer conn.Close()

	icmpPacket := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
		ID: os.Getpid() & 0xffff, Seq: 1,
		Data: []byte("NITPING"),
		},
	}
	icmpBytes, err := icmpPacket.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := conn.WriteTo(icmpBytes, &net.IPAddr{IP: net.ParseIP(*dstipPtr)}); err !=nil {
		log.Fatalf("WriteTo error, %s", err)
	}
	
	//cm.Src = netAddrToIP4(*srcipPtr)
	readBytes := make([]byte, 1500)
	n, peer, err := conn.ReadFrom(readBytes)
	if err != nil {
		log.Fatal(err)
	}

	remove, err := icmp.ParseMessage(iana.ProtocolICMP, readBytes[:n])
	if err != nil {
		log.Fatal(err)
	}

	switch remove.Type {
	case ipv4.ICMPTypeEchoReply:
		log.Printf("Got NITPONG from %v", peer)
	default:
		log.Printf("Got %+v; Wanted an NITPONG(echo reply)", remove)
	}
}

