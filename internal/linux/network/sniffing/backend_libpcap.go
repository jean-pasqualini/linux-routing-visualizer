//go:build SNIFF_LIBPCAP

package sniffing

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/logging"
)

func (s *SniffingBackend) Sniff(ctx context.Context, iface, filter string, snaplen int32, promisc bool, timeout time.Duration) (chan Packet, error) {
	logger := logging.FromContext(ctx)
	out := make(chan Packet, 4096)
	if iface == "" {
		// Liste les interfaces dispo
		devs, err := pcap.FindAllDevs()
		if err != nil {
			//logger.Error("Impossible de lister les interfaces: %v", err)
			close(out)
			return out, errors.New(err.Error())
		}
		//logger.Info("Interfaces disponibles :")
		for _, _ = range devs {
			//logger.Info("- %s (%s)\n", d.Name, d.Description)
		}
		close(out)
		return out, errors.New("")
	}

	go func() {
		defer close(out)
		handle, err := pcap.OpenLive(iface, int32(snaplen), promisc, timeout)
		if err != nil {
			logger.Error(fmt.Sprintf("OpenLive échoué: %v", err))
			return
		}
		defer handle.Close()

		if filter != "" {
			if err := handle.SetBPFFilter(filter); err != nil {
				//logger.Error("Filtre BPF invalide: %v", err)
			}
			//logger.Info("Filtre appliqué: %q", filter)
		}

		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		//packetSource.NoCopy = true

		logger.Info("Capture sur %s (Ctrl+C pour arrêter)...", iface)
		for {
			select {
			case <-ctx.Done():
				logger.Info("You requested to stop the capture...")
				return

			case pkt, ok := <-packetSource.Packets():
				if !ok {
					logger.Info("Packet channel closed")
					return
				}
				if result := convert(logger, pkt); result != nil {
					out <- *result
				}
			}
		}
	}()

	return out, nil
}

func convert(logger *slog.Logger, packet gopacket.Packet) *Packet {
	// Timestamp & longueur
	ci := packet.Metadata().CaptureInfo
	logger.Info(fmt.Sprintf("\n[%s] len=%d caplen=%d\n", ci.Timestamp.Format(time.RFC3339Nano), ci.Length, ci.CaptureLength))

	// IPv4/IPv6
	var srcIP, dstIP string
	if ip4 := packet.Layer(layers.LayerTypeIPv4); ip4 != nil {
		l := ip4.(*layers.IPv4)
		srcIP, dstIP = l.SrcIP.String(), l.DstIP.String()
		logger.Info(fmt.Sprintf("IPv4  %s -> %s  proto=%s\n", srcIP, dstIP, l.Protocol))
		return &Packet{
			Source: srcIP,
			Target: dstIP,
		}
	} else if ip6 := packet.Layer(layers.LayerTypeIPv6); ip6 != nil {
		//	l := ip6.(*layers.IPv6)
		//	srcIP, dstIP = l.SrcIP.String(), l.DstIP.String()
		//logger.Info(fmt.Sprintf("IPv6  %s -> %s  next=%s\n", srcIP, dstIP, l.NextHeader))
	} else {
		// Pas d'IP (ARP, etc.)
		///logger.Info(fmt.Sprintln("Non-IP (ex: ARP/LLDP/etc.)"))
	}

	return nil

	// TCP / UDP
	/**
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

	*/
}
