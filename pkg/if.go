package ethtool

const IFNAMESIZE = 16

type if_nameindex struct {
	if_index uint  /* 1, 2, ... */
	if_name  *byte /* null terminated name: "eth0", ... */
}

const (
	IFF_UP          = 0x1
	IFF_BROADCAST   = 0x2
	IFF_DEBUG       = 0x4
	IFF_LOOPBACK    = 0x8
	IFF_POINTOPOINT = 0x10
	IFF_NOTRAILERS  = 0x20
	IFF_RUNNING     = 0x40
	IFF_NOARP       = 0x80
	IFF_PROMISC     = 0x100
	IFF_ALLMULTI    = 0x200
	IFF_MASTER      = 0x400
	IFF_SLAVE       = 0x800
	IFF_MULTICAST   = 0x1000
	IFF_PORTSEL     = 0x2000
	IFF_AUTOMEDIA   = 0x4000
	IFF_DYNAMIC     = 0x8000
)

type sockaddr struct {
	sa_family uint16
	sa_data   [14]byte
}

type ifaddr struct {
	ifa_addr sockaddr
	ifu_addr sockaddr
	ifa_ifp  uintptr
	ifa_next *ifaddr
}

type ifmap struct {
	mem_start uint64
	mem_end   uint64
	base_addr uint16
	irq       uint8
	dma       uint8
	port      uint8
}

type ifreq struct {
	ifr_name [IFNAMESIZE]byte
	ifr_data uintptr
}

type ifconf struct {
	ifc_len  int
	ifcu_buf uintptr
}
