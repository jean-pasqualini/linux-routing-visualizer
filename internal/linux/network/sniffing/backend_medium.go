//go:build medium
// +build medium

package sniffing

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

type Sniffing struct{}

func NewSniffing() *Sniffing { return &Sniffing{} }

func htons(i uint16) uint16 { return (i<<8)&0xff00 | i>>8 }

const (
	// From linux/if_packet.h
	PACKET_RX_RING = 5
	PACKET_VERSION = 10

	TPACKET_V3 = 2

	TP_STATUS_KERNEL = 0
	TP_STATUS_USER   = 1
)

type TpacketReq3 struct {
	BlockSize      uint32
	BlockNr        uint32
	FrameSize      uint32
	FrameNr        uint32
	RetireBlkTov   uint32
	SizeofPriv     uint32
	FeatureReqWord uint32
}

type TpacketBlockDesc struct {
	Version uint32
	Offset  uint32
}

type TpacketHdrV1 struct {
	BlockStatus      uint32
	NumPkts          uint32
	OffsetToFirstPkt uint32
	BlkLen           uint32
	SeqNum           uint64
	TsFirstPktSec    uint64
	TsFirstPktNSec   uint64
}

// Minimal tpacket3_hdr (enough for NextOffset/Mac/SnapLen)
type Tpacket3Hdr struct {
	NextOffset uint32
	Sec        uint32
	NSec       uint32
	SnapLen    uint32
	Len        uint32
	Status     uint32
	Mac        uint16
	Net        uint16
	Pad        uint16
}

// Sniff prints Ethernet frames as they arrive.
func (s *Sniffing) Sniff(ctx context.Context, iface string, filter string, snaplen int32, promisc bool, timeout time.Duration) error {
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	// Bind to interface
	ifi, err := net.InterfaceByName(iface)
	if err != nil {
		return err
	}
	if err := syscall.Bind(fd, &syscall.SockaddrLinklayer{
		Protocol: htons(syscall.ETH_P_ALL),
		Ifindex:  ifi.Index,
	}); err != nil {
		return err
	}

	// Optional promiscuous mode
	if promisc {
	}

	/**

	FrameSize: 0, // obligatoire
	FrameNr:   0, // obligatoire

	const TP_FT_REQ_FILL_RXHASH = 1 << 0

	req.FeatureReqWord = TP_FT_REQ_FILL_RXHASH
	FeatureReqWord: 0, // <-- OBLIGATOIRE avec filtre
	*/
	//// Filter
	prog := unix.SockFprog{
		Len:    uint16(len(filters)),
		Filter: &filters[0],
	}

	if err := unix.SetsockoptSockFprog(fd, unix.SOL_SOCKET, unix.SO_ATTACH_FILTER, &prog); err != nil {
		fmt.Println("noooo")
		return err
	}

	////

	// Configure TPACKET_V3 ring
	req := TpacketReq3{
		BlockSize:    1 << 20, // 1MB
		BlockNr:      64,
		FrameSize:    2048,                  // not used by V3 in the same way, but must be non-zero
		FrameNr:      (1 << 20) * 64 / 2048, // ✅ parentheses fixed
		RetireBlkTov: 60,                    // ms
	}

	if err := syscall.SetsockoptInt(fd, syscall.SOL_PACKET, PACKET_VERSION, TPACKET_V3); err != nil {
		return err
	}

	_, _, errno := syscall.Syscall6(
		syscall.SYS_SETSOCKOPT,
		uintptr(fd),
		uintptr(syscall.SOL_PACKET),
		uintptr(PACKET_RX_RING),
		uintptr(unsafe.Pointer(&req)),
		unsafe.Sizeof(req),
		0,
	)
	if errno != 0 {
		return errno
	}

	mmapSize := int(req.BlockSize * req.BlockNr)
	data, err := syscall.Mmap(fd, 0, mmapSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return err
	}
	defer syscall.Munmap(data)

	// poll(2) via x/sys/unix
	pfds := []unix.PollFd{{Fd: int32(fd), Events: unix.POLLIN}}
	pollTimeoutMs := int(timeout.Milliseconds())
	if pollTimeoutMs <= 0 {
		pollTimeoutMs = 1000
	}

	for i := 0; ; i = (i + 1) % int(req.BlockNr) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		block := data[i*int(req.BlockSize):]
		//desc := (*TpacketBlockDesc)(unsafe.Pointer(&block[0]))
		//hdr := (*TpacketHdrV1)(unsafe.Pointer(&block[8]))
		desc := (*TpacketBlockDesc)(unsafe.Pointer(&block[0]))
		// tpacket_hdr_v1 est juste après TpacketBlockDesc (8 bytes)
		hdr := (*TpacketHdrV1)(unsafe.Pointer(&block[unsafe.Sizeof(*desc)]))

		if hdr.BlockStatus&TP_STATUS_USER == 0 {
			_, err = unix.Poll(pfds, pollTimeoutMs)
			if err != nil {
				fmt.Println(err.Error())
			}
			continue
		}

		offset := int(hdr.OffsetToFirstPkt)
		for p := 0; p < int(hdr.NumPkts); p++ {
			pktHdr := (*Tpacket3Hdr)(unsafe.Pointer(&block[offset]))

			start := offset + int(pktHdr.Mac)
			end := start + int(pktHdr.SnapLen)

			if start < 0 || end > len(block) || end < start {
				break
			}

			frame := block[start:end]
			printEthernet(frame)
			decodeIPv4(frame)

			if pktHdr.NextOffset == 0 {
				break
			}
			offset += int(pktHdr.NextOffset)
		}

		// Release block back to kernel
		hdr.BlockStatus = TP_STATUS_KERNEL
	}
}

func decodeIPv4(frame []byte) {
	if len(frame) < 34 {
		return
	}

	ethType := binary.BigEndian.Uint16(frame[12:14])
	if ethType != 0x0800 {
		return // pas IPv4
	}

	ip := frame[14:]
	ihl := int(ip[0]&0x0F) * 4
	if len(ip) < ihl {
		return
	}

	src := net.IP(ip[12:16])
	dst := net.IP(ip[16:20])
	proto := ip[9]

	fmt.Printf("IPv4 %s → %s proto=%d\n", src, dst, proto)
}

func printEthernet(b []byte) {
	if len(b) < 14 {
		return
	}
	dst := net.HardwareAddr(b[:6])
	src := net.HardwareAddr(b[6:12])
	eth := binary.BigEndian.Uint16(b[12:14])
	fmt.Printf("[%s → %s] EtherType=0x%04x Len=%d\n", src, dst, eth, len(b))
}
