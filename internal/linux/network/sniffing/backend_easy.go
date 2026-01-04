//go:build easy
// +build easy

package sniffing

import (
	"context"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"time"
)

type Sniffing struct {
}

func NewSniffing() *Sniffing {
	return &Sniffing{}
}

func (s *Sniffing) Sniff(ctx context.Context, iface, filter string, snaplen int32, promisc bool, timeout time.Duration) {
	if iface == "" {
		// Liste les interfaces dispo
		devs, err := pcap.FindAllDevs()
		if err != nil {
			log.Fatalf("Impossible de lister les interfaces: %v", err)
		}
		fmt.Println("Interfaces disponibles :")
		for _, d := range devs {
			fmt.Printf("- %s (%s)\n", d.Name, d.Description)
		}
		fmt.Println("\nRelance avec: -i <interface> [-f <filtre>]")
		return
	}

	handle, err := pcap.OpenLive(iface, int32(snaplen), promisc, timeout)
	if err != nil {
		log.Fatalf("OpenLive échoué: %v", err)
	}
	defer handle.Close()

	if filter != "" {
		if err := handle.SetBPFFilter(filter); err != nil {
			log.Fatalf("Filtre BPF invalide: %v", err)
		}
		log.Printf("Filtre appliqué: %q", filter)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	//packetSource.NoCopy = true

	log.Printf("Capture sur %s (Ctrl+C pour arrêter)...", iface)
	for {
		select {
		case <-ctx.Done():
			log.Println("You requested to stop the capture...")
			return
		case pkt, ok := <-packetSource.Packets():
			if !ok {
				return
			}
			printPacket(pkt)
		}
	}
}

func printPacket(packet gopacket.Packet) {
	// Timestamp & longueur
	ci := packet.Metadata().CaptureInfo
	fmt.Printf("\n[%s] len=%d caplen=%d\n", ci.Timestamp.Format(time.RFC3339Nano), ci.Length, ci.CaptureLength)

	// IPv4/IPv6
	var srcIP, dstIP string
	if ip4 := packet.Layer(layers.LayerTypeIPv4); ip4 != nil {
		l := ip4.(*layers.IPv4)
		srcIP, dstIP = l.SrcIP.String(), l.DstIP.String()
		fmt.Printf("IPv4  %s -> %s  proto=%s\n", srcIP, dstIP, l.Protocol)
	} else if ip6 := packet.Layer(layers.LayerTypeIPv6); ip6 != nil {
		l := ip6.(*layers.IPv6)
		srcIP, dstIP = l.SrcIP.String(), l.DstIP.String()
		fmt.Printf("IPv6  %s -> %s  next=%s\n", srcIP, dstIP, l.NextHeader)
	} else {
		// Pas d'IP (ARP, etc.)
		fmt.Println("Non-IP (ex: ARP/LLDP/etc.)")
	}

	// TCP / UDP
	if tcpL := packet.Layer(layers.LayerTypeTCP); tcpL != nil {
		t := tcpL.(*layers.TCP)
		fmt.Printf("TCP   %s:%d -> %s:%d  flags=[S:%v A:%v F:%v R:%v P:%v U:%v]\n",
			srcIP, t.SrcPort, dstIP, t.DstPort,
			t.SYN, t.ACK, t.FIN, t.RST, t.PSH, t.URG,
		)
		return
	}
	if udpL := packet.Layer(layers.LayerTypeUDP); udpL != nil {
		u := udpL.(*layers.UDP)
		fmt.Printf("UDP   %s:%d -> %s:%d\n", srcIP, u.SrcPort, dstIP, u.DstPort)
		return
	}

	// Autres couches utiles (ex: DNS)
	if dnsL := packet.Layer(layers.LayerTypeDNS); dnsL != nil {
		d := dnsL.(*layers.DNS)
		fmt.Printf("DNS   q=%d a=%d rcode=%v\n", len(d.Questions), len(d.Answers), d.ResponseCode)
	}
}
