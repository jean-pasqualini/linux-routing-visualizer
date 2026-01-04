package sniffing

import "golang.org/x/sys/unix"

var afilters = []unix.SockFilter{
	{Code: 0x28, Jt: 0, Jf: 0, K: 0x0000000c},
	{Code: 0x15, Jt: 0, Jf: 8, K: 0x00000800},
	{Code: 0x30, Jt: 0, Jf: 0, K: 0x00000017},
	{Code: 0x15, Jt: 0, Jf: 6, K: 0x00000006},
	{Code: 0x28, Jt: 0, Jf: 0, K: 0x00000014},
	{Code: 0x45, Jt: 4, Jf: 0, K: 0x00001fff},
	{Code: 0x28, Jt: 0, Jf: 0, K: 0x00000016},
	{Code: 0x15, Jt: 0, Jf: 1, K: 0x00000050}, // port 80
	{Code: 0x06, Jt: 0, Jf: 0, K: 0x00040000},
	{Code: 0x06, Jt: 0, Jf: 0, K: 0x00000000},
}

var filters = []unix.SockFilter{
	// EtherType == IPv4
	{Code: 0x28, Jt: 0, Jf: 9, K: 0x0000000c}, // ldh [12]
	{Code: 0x15, Jt: 0, Jf: 8, K: 0x00000800}, // jeq #0x0800

	// IP proto == TCP
	{Code: 0x30, Jt: 0, Jf: 6, K: 0x00000017}, // ldb [23]
	{Code: 0x15, Jt: 0, Jf: 5, K: 0x00000006}, // jeq #6

	// TCP dst port == 80
	{Code: 0x28, Jt: 0, Jf: 3, K: 0x00000024}, // ldh [36]
	{Code: 0x15, Jt: 0, Jf: 2, K: 0x00000050}, // jeq #80

	// Accept
	{Code: 0x06, Jt: 0, Jf: 0, K: 0x00040000},
	// Reject
	{Code: 0x06, Jt: 0, Jf: 0, K: 0x00000000},
}
