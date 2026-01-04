package sniffing

// sudo setcap cap_net_raw,cap_net_admin+ep ./yourbin
// getcap $(which tcpdump)

/**
struct tpacket_block_desc {
  __u32 version;
  __u32 offset_to_priv;   // ⚠️ PAS un offset vers le header
  struct tpacket_hdr_v1 hdr;   // ✅ le header est juste après
};

*/

/**
var filters = []unix.SockFilter{
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

*/
