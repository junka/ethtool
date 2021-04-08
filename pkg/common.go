package ethtool

import "fmt"

const OFF_FLAG_DEF_SIZE = 12

const (
	DEBUG_PARSE       = iota
	DEBUG_NL_MSGS     /* incoming/outgoing netlink messages */
	DEBUG_NL_DUMP_SND /* dump outgoing netlink messages */
	DEBUG_NL_DUMP_RCV /* dump incoming netlink messages */
	DEBUG_NL_PRETTY_MSG
)

func debug_on(debug uint64, bit uint32) bool {
	return (debug & (1 << bit)) != 0
}

const (
	ETH_FLAG_RXCSUM   = (1 << 0)
	ETH_FLAG_TXCSUM   = (1 << 1)
	ETH_FLAG_SG       = (1 << 2)
	ETH_FLAG_TSO      = (1 << 3)
	ETH_FLAG_UFO      = (1 << 4)
	ETH_FLAG_GSO      = (1 << 5)
	ETH_FLAG_GRO      = (1 << 6)
	ETH_FLAG_INT_MASK = (ETH_FLAG_RXCSUM | ETH_FLAG_TXCSUM |
		ETH_FLAG_SG | ETH_FLAG_TSO | ETH_FLAG_UFO |
		ETH_FLAG_GSO | ETH_FLAG_GRO)
	ETH_FLAG_EXT_MASK = (ETH_FLAG_LRO | ETH_FLAG_RXVLAN |
		ETH_FLAG_TXVLAN | ETH_FLAG_NTUPLE |
		ETH_FLAG_RXHASH)
)

type link_mode struct {
	supported      [127]byte
	advertising    [127]byte
	lp_advertising [127]byte
}
type ethtool_link_usettings struct {
	// struct {
	// 	__u8 transceiver;
	// } deprecated;
	base       ethtool_link_settings
	link_modes link_mode
}

type flag_info struct {
	name  string
	value uint32
}

type off_flag_def_t struct {
	short_name  string
	long_name   string
	kernel_name string
	get_cmd     uint32
	set_cmd     uint32
	value       uint32
	/* For features exposed through ETHTOOL_GFLAGS, the oldest
	 * kernel version for which we can trust the result.  Where
	 * the flag was added at the same time the kernel started
	 * supporting the feature, this is 0 (to allow for backports).
	 * Where the feature was supported before the flag was added,
	 * it is the version that introduced the flag.
	 */
	min_kernel_ver uint32
}

const (
	NETIF_MSG_DRV       = 0x0001
	NETIF_MSG_PROBE     = 0x0002
	NETIF_MSG_LINK      = 0x0004
	NETIF_MSG_TIMER     = 0x0008
	NETIF_MSG_IFDOWN    = 0x0010
	NETIF_MSG_IFUP      = 0x0020
	NETIF_MSG_RX_ERR    = 0x0040
	NETIF_MSG_TX_ERR    = 0x0080
	NETIF_MSG_TX_QUEUED = 0x0100
	NETIF_MSG_INTR      = 0x0200
	NETIF_MSG_TX_DONE   = 0x0400
	NETIF_MSG_RX_STATUS = 0x0800
	NETIF_MSG_PKTDATA   = 0x1000
	NETIF_MSG_HW        = 0x2000
	NETIF_MSG_WOL       = 0x4000
)

var flags_msglvl = []flag_info{
	{"drv", NETIF_MSG_DRV},
	{"probe", NETIF_MSG_PROBE},
	{"link", NETIF_MSG_LINK},
	{"timer", NETIF_MSG_TIMER},
	{"ifdown", NETIF_MSG_IFDOWN},
	{"ifup", NETIF_MSG_IFUP},
	{"rx_err", NETIF_MSG_RX_ERR},
	{"tx_err", NETIF_MSG_TX_ERR},
	{"tx_queued", NETIF_MSG_TX_QUEUED},
	{"intr", NETIF_MSG_INTR},
	{"tx_done", NETIF_MSG_TX_DONE},
	{"rx_status", NETIF_MSG_RX_STATUS},
	{"pktdata", NETIF_MSG_PKTDATA},
	{"hw", NETIF_MSG_HW},
	{"wol", NETIF_MSG_WOL},
}

var n_flags_msglvl = len(flags_msglvl)

var off_flag_def = []off_flag_def_t{
	{"rx", "rx-checksumming", "rx-checksum",
		ETHTOOL_GRXCSUM, ETHTOOL_SRXCSUM, ETH_FLAG_RXCSUM, 0},
	{"tx", "tx-checksumming", "tx-checksum-*",
		ETHTOOL_GTXCSUM, ETHTOOL_STXCSUM, ETH_FLAG_TXCSUM, 0},
	{"sg", "scatter-gather", "tx-scatter-gather*",
		ETHTOOL_GSG, ETHTOOL_SSG, ETH_FLAG_SG, 0},
	{"tso", "tcp-segmentation-offload", "tx-tcp*-segmentation",
		ETHTOOL_GTSO, ETHTOOL_STSO, ETH_FLAG_TSO, 0},
	{"ufo", "udp-fragmentation-offload", "tx-udp-fragmentation",
		ETHTOOL_GUFO, ETHTOOL_SUFO, ETH_FLAG_UFO, 0},
	{"gso", "generic-segmentation-offload", "tx-generic-segmentation",
		ETHTOOL_GGSO, ETHTOOL_SGSO, ETH_FLAG_GSO, 0},
	{"gro", "generic-receive-offload", "rx-gro",
		ETHTOOL_GGRO, ETHTOOL_SGRO, ETH_FLAG_GRO, 0},
	{"lro", "large-receive-offload", "rx-lro",
		0, 0, ETH_FLAG_LRO, 0x020618},
	{"rxvlan", "rx-vlan-offload", "rx-vlan-hw-parse",
		0, 0, ETH_FLAG_RXVLAN, 0x020625},
	{"txvlan", "tx-vlan-offload", "tx-vlan-hw-insert",
		0, 0, ETH_FLAG_TXVLAN, 0x020625},
	{"ntuple", "ntuple-filters", "rx-ntuple-filter",
		0, 0, ETH_FLAG_NTUPLE, 0},
	{"rxhash", "receive-hashing", "rx-hashing",
		0, 0, ETH_FLAG_RXHASH, 0},
}

func print_flags(info []flag_info, n_info uint32, value uint32) {

	var sep string = ""
	var i = 0
	for n_info > 0 {
		if value&info[i].value > 0 {
			fmt.Printf("%s%s", sep, string(info[i].name))
			sep = " "
			value &= ^info[0].value
		}
		i++
		n_info--
	}

	if value > 0 {
		fmt.Printf("%s%#x", sep, value)
	}
}

func unparse_wolopts(wolopts int) []byte {
	buf := make([]byte, 0)
	if wolopts > 0 {
		if wolopts&WAKE_PHY > 0 {
			buf = append(buf, 'p')
		}
		if wolopts&WAKE_UCAST > 0 {
			buf = append(buf, 'u')
		}
		if wolopts&WAKE_MCAST > 0 {
			buf = append(buf, 'm')
		}
		if wolopts&WAKE_BCAST > 0 {
			buf = append(buf, 'b')
		}
		if wolopts&WAKE_ARP > 0 {
			buf = append(buf, 'a')
		}
		if wolopts&WAKE_MAGIC > 0 {
			buf = append(buf, 'g')
		}
		if wolopts&WAKE_MAGICSECURE > 0 {
			buf = append(buf, 's')
		}
		if wolopts&WAKE_FILTER > 0 {
			buf = append(buf, 'f')
		}
	} else {
		buf = append(buf, 'd')
	}

	return buf
}

func dump_wol(wol *ethtool_wolinfo) {
	fmt.Printf("	Supports Wake-on: %s\n", unparse_wolopts(wol.supported))
	fmt.Printf("	Wake-on: %s\n", unparse_wolopts(wol.wolopts))
	if wol.supported&WAKE_MAGICSECURE > 0 {
		delim := 0
		fmt.Printf("        SecureOn password: ")
		for i := 0; i < SOPASS_MAX; i++ {
			if delim == 0 {
				fmt.Printf("%s%02x", "", wol.sopass[i])
			} else {
				fmt.Printf("%s%02x", ":", wol.sopass[i])
			}
			delim = 1
		}
		fmt.Printf("\n")
	}
}

func dump_mdix(mdix uint8, mdix_ctrl uint8) {
	fmt.Printf("	MDI-X: ")
	if mdix_ctrl == ETH_TP_MDI {
		fmt.Printf("off (forced)\n")
	} else if mdix_ctrl == ETH_TP_MDI_X {
		fmt.Printf("on (forced)\n")
	} else {
		switch mdix {
		case ETH_TP_MDI:
			fmt.Printf("off")
			break
		case ETH_TP_MDI_X:
			fmt.Printf("on")
			break
		default:
			fmt.Printf("Unknown")
			break
		}
		if mdix_ctrl == ETH_TP_MDI_AUTO {
			fmt.Printf(" (auto)")
		}
		fmt.Printf("\n")
	}
}

//from internal
type cmd_context struct {
	devname    string /* net device name */
	fd         int    /* socket suitable for ethtool ioctl */
	ifr        ifreq  /* ifreq suitable for ethtool ioctl */
	argc       int    /* number of arguments to the sub-command */
	argp       **byte /* arguments to the sub-command */
	debug      uint64 /* debugging mask */
	json       bool   /* Output JSON, if supported */
	show_stats bool   /* include command-specific stats */
	//netlink
}
