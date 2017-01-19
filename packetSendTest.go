package main

import (
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/examples/util"
	"code.google.com/p/gopacket/layers"
	"golang.org/x/net/ipv4"
	"log"
	"net"
)

func main() {
	defer util.Run()()

	var srcMac, dstMac net.HardwareAddr
	var srcMacstr string = "01:01:de:ad:be:ef"
	var dstMacstr string = "01:02:be:ef:de:ad"
	
	var srcIP, dstIP net.IP
	var srcIPstr string = "127.0.0.1"
	var dstIPstr string = "127.0.0.1"

	// source ip
	srcIP = net.ParseIP(srcIPstr)
	if srcIP == nil {
		log.Printf("non-ip target: %q\n", srcIPstr)
	}
	srcIP = srcIP.To4()
	if srcIP == nil {
		log.Printf("non-ipv4 target: %q\n", srcIPstr)
	}

	// destination ip
	dstIP = net.ParseIP(dstIPstr)
	if dstIP == nil {
		log.Printf("non-ip target: %q\n", dstIPstr)
	}
	dstIP = dstIP.To4()
	if dstIP == nil {
		log.Printf("non-ipv4 target: %q\n", dstIPstr)
	}

	srcMac, err1 := net.ParseMAC(srcMacstr)
	if err1 != nil {
		log.Printf("Src Mac wrong: %q\n", srcMacstr)
	}
	
	dstMac, err2 := net.ParseMAC(dstMacstr)
	if err2 != nil {
		log.Printf("Dst Mac wrong: %q\n", dstMacstr)
	}
	
	ethernet := layers.Ethernet{
			SrcMAC: srcMac,
			DstMAC: dstMac,
			EthernetType: layers.EthernetTypeIPv4,
			Length: 0,
	}

	// build tcp/ip packet
	ip := layers.IPv4{
		SrcIP:    srcIP,
		DstIP:    dstIP,
		Version:  4,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
	}

	srcport := layers.TCPPort(666)
	dstport := layers.TCPPort(1022)
	tcp := layers.TCP{
		SrcPort: srcport,
		DstPort: dstport,
		Window:  1505,
		Urgent:  0,
		Seq:     11050,
		Ack:     0,
		ACK:     false,
		SYN:     false,
		FIN:     false,
		RST:     false,
		URG:     false,
		ECE:     false,
		CWR:     false,
		NS:      false,
		PSH:     false,
	}

	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	tcp.SetNetworkLayerForChecksum(&ip)

	ipHeaderBuf := gopacket.NewSerializeBuffer()
	err := ip.SerializeTo(ipHeaderBuf, opts)
	if err != nil {
		panic(err)
	}
	ipHeader, err := ipv4.ParseHeader(ipHeaderBuf.Bytes())
	if err != nil {
		panic(err)
	}

	tcpPayloadBuf := gopacket.NewSerializeBuffer()
	payload := gopacket.Payload([]byte("meowmeowmeow"))
	err = gopacket.SerializeLayers(tcpPayloadBuf, opts, &ethernet, &ip, &tcp, payload)
	if err != nil {
		panic(err)
	}

	var packetConn net.PacketConn
	var rawConn *ipv4.RawConn
	packetConn, err = net.ListenPacket("ip4:tcp", "127.0.0.1")
	if err != nil {
		panic(err)
	}
	rawConn, err = ipv4.NewRawConn(packetConn)
	if err != nil {
		panic(err)
	}

	err = rawConn.WriteTo(ipHeader, tcpPayloadBuf.Bytes(), nil)
	log.Printf("Packet of length %d sent!\n", (len(tcpPayloadBuf.Bytes()) + len(ipHeaderBuf.Bytes())))
}
