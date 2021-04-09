package ethtool

import (
	"fmt"
	"net"
	"unsafe"
)

func invert_flow_mask(fsp *ethtool_rx_flow_spec) {

	for i := 0; i < int(unsafe.Sizeof(fsp.m_u)); i++ {
		fsp.m_u.hdata[i] ^= 0xFF
	}
}

func uint32ToIP(intIP uint32) net.IP {
	var bytes [4]byte
	bytes[0] = byte(intIP & 0xFF)
	bytes[1] = byte((intIP >> 8) & 0xFF)
	bytes[2] = byte((intIP >> 16) & 0xFF)
	bytes[3] = byte((intIP >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

func rxclass_print_ipv4_rule(sip uint32, sipm uint32, dip uint32,
	dipm uint32, tos uint8, tosm uint8) {
	fmt.Printf(
		"\tSrc IP addr: %s mask: %s\n"+
			"\tDest IP addr: %s mask: %s\n"+
			"\tTOS: 0x%x mask: 0x%x\n",
		uint32ToIP(sip).String(), uint32ToIP(sipm).String(),
		uint32ToIP(dip).String(), uint32ToIP(dipm).String(),
		tos, tosm)
}

func rxclass_print_ipv6_rule(sip uint32, sipm uint32, dip uint32,
	dipm uint32, tclass uint8, tclassm uint8) {

	// fmt.Printf(
	// 	"\tSrc IP addr: %s mask: %s\n"+
	// 	"\tDest IP addr: %s mask: %s\n"+
	// 	"\tTraffic Class: 0x%x mask: 0x%x\n",
	// 	inet_ntop(AF_INET6, sip, sip_str, INET6_ADDRSTRLEN),
	// 	inet_ntop(AF_INET6, sipm, sipm_str, INET6_ADDRSTRLEN),
	// 	inet_ntop(AF_INET6, dip, dip_str, INET6_ADDRSTRLEN),
	// 	inet_ntop(AF_INET6, dipm, dipm_str, INET6_ADDRSTRLEN),
	// 	tclass, tclassm);
}

func rxclass_print_nfc_spec_ext(fsp *ethtool_rx_flow_spec) {
	if (fsp.flow_type & FLOW_EXT) != 0 {
		// var data, datam uint64
		// var etype, etypem, tci, tcim uint16
		// etype = ntohs(fsp.h_ext.vlan_etype);
		// etypem = ntohs(^fsp.m_ext.vlan_etype);
		// tci = ntohs(fsp.h_ext.vlan_tci);
		// tcim = ntohs(^fsp.m_ext.vlan_tci);
		// data = (u64)ntohl(fsp.h_ext.data[0]) << 32;
		// data |= (u64)ntohl(fsp.h_ext.data[1]);
		// datam = (u64)ntohl(~fsp.m_ext.data[0]) << 32;
		// datam |= (u64)ntohl(~fsp.m_ext.data[1]);

		// fmt.Printf(
		// 	"\tVLAN EtherType: 0x%x mask: 0x%x\n"+
		// 	"\tVLAN: 0x%x mask: 0x%x\n"+
		// 	"\tUser-defined: 0x%llx mask: 0x%llx\n",
		// 	etype, etypem, tci, tcim, data, datam);
	}

	if (fsp.flow_type & FLOW_MAC_EXT) != 0 {

		dmac := fsp.h_ext.h_dest
		dmacm := fsp.m_ext.h_dest

		fmt.Printf(
			"\tDest MAC addr: %02X:%02X:%02X:%02X:%02X:%02X"+
				" mask: %02X:%02X:%02X:%02X:%02X:%02X\n",
			dmac[0], dmac[1], dmac[2], dmac[3], dmac[4],
			dmac[5], dmacm[0], dmacm[1], dmacm[2], dmacm[3],
			dmacm[4], dmacm[5])
	}
}

func rxclass_print_nfc_rule(fsp *ethtool_rx_flow_spec,
	rss_context uint32) {
	// unsigned char	*smac, *smacm, *dmac, *dmacm;

	fmt.Printf("Filter: %d\n", fsp.location)

	flow_type := fsp.flow_type & (^(uint32(FLOW_EXT) | FLOW_MAC_EXT | FLOW_RSS))

	invert_flow_mask(fsp)

	switch flow_type {
	case TCP_V4_FLOW:
	case UDP_V4_FLOW:
	case SCTP_V4_FLOW:
		if flow_type == TCP_V4_FLOW {
			fmt.Printf("\tRule Type: TCP over IPv4\n")
		} else if flow_type == UDP_V4_FLOW {
			fmt.Printf("\tRule Type: UDP over IPv4\n")
		} else {
			fmt.Printf("\tRule Type: SCTP over IPv4\n")
		}
		// rxclass_print_ipv4_rule(fsp.h_u.tcp_ip4_spec.ip4src,
		// 	fsp.m_u.tcp_ip4_spec.ip4src,
		// 	fsp.h_u.tcp_ip4_spec.ip4dst,
		// 	fsp.m_u.tcp_ip4_spec.ip4dst,
		// 	fsp.h_u.tcp_ip4_spec.tos,
		// 	fsp.m_u.tcp_ip4_spec.tos)
		// fmt.Printf(
		// 	"\tSrc port: %d mask: 0x%x\n"
		// 	"\tDest port: %d mask: 0x%x\n",
		// 	binary.littleEndian (fsp.h_u.tcp_ip4_spec.psrc),
		// 	ntohs(fsp.m_u.tcp_ip4_spec.psrc),
		// 	ntohs(fsp.h_u.tcp_ip4_spec.pdst),
		// 	ntohs(fsp.m_u.tcp_ip4_spec.pdst));
		break
	case AH_V4_FLOW:
	case ESP_V4_FLOW:
		if flow_type == AH_V4_FLOW {
			fmt.Printf("\tRule Type: IPSEC AH over IPv4\n")
		} else {
			fmt.Printf("\tRule Type: IPSEC ESP over IPv4\n")
		}
		// rxclass_print_ipv4_rule(fsp.h_u.ah_ip4_spec.ip4src,
		// 	fsp.m_u.ah_ip4_spec.ip4src,
		// 	fsp.h_u.ah_ip4_spec.ip4dst,
		// 	fsp.m_u.ah_ip4_spec.ip4dst,
		// 	fsp.h_u.ah_ip4_spec.tos,
		// 	fsp.m_u.ah_ip4_spec.tos)
		// fmt.Printf(
		// 	"\tSPI: %d mask: 0x%x\n",
		// 	ntohl(fsp.h_u.esp_ip4_spec.spi),
		// 	ntohl(fsp.m_u.esp_ip4_spec.spi));
		break
	case IPV4_USER_FLOW:
		fmt.Printf("\tRule Type: Raw IPv4\n")
		// rxclass_print_ipv4_rule(fsp.h_u.usr_ip4_spec.ip4src,
		// 	fsp.m_u.usr_ip4_spec.ip4src,
		// 	fsp.h_u.usr_ip4_spec.ip4dst,
		// 	fsp.m_u.usr_ip4_spec.ip4dst,
		// 	fsp.h_u.usr_ip4_spec.tos,
		// 	fsp.m_u.usr_ip4_spec.tos)
		// fmt.Printf(
		// 	"\tProtocol: %d mask: 0x%x\n"
		// 	"\tL4 bytes: 0x%x mask: 0x%x\n",
		// 	fsp.h_u.usr_ip4_spec.proto,
		// 	fsp.m_u.usr_ip4_spec.proto,
		// 	ntohl(fsp.h_u.usr_ip4_spec.l4_4_bytes),
		// 	ntohl(fsp.m_u.usr_ip4_spec.l4_4_bytes));
		break
	case TCP_V6_FLOW:
	case UDP_V6_FLOW:
	case SCTP_V6_FLOW:
		if flow_type == TCP_V6_FLOW {
			fmt.Printf("\tRule Type: TCP over IPv6\n")
		} else if flow_type == UDP_V6_FLOW {
			fmt.Printf("\tRule Type: UDP over IPv6\n")
		} else {
			fmt.Printf("\tRule Type: SCTP over IPv6\n")
		}
		// rxclass_print_ipv6_rule(fsp.h_u.tcp_ip6_spec.ip6src,
		// 	fsp.m_u.tcp_ip6_spec.ip6src,
		// 	fsp.h_u.tcp_ip6_spec.ip6dst,
		// 	fsp.m_u.tcp_ip6_spec.ip6dst,
		// 	fsp.h_u.tcp_ip6_spec.tclass,
		// 	fsp.m_u.tcp_ip6_spec.tclass)
		// fmt.Printf(
		// 	"\tSrc port: %d mask: 0x%x\n"
		// 	"\tDest port: %d mask: 0x%x\n",
		// 	ntohs(fsp.h_u.tcp_ip6_spec.psrc),
		// 	ntohs(fsp.m_u.tcp_ip6_spec.psrc),
		// 	ntohs(fsp.h_u.tcp_ip6_spec.pdst),
		// 	ntohs(fsp.m_u.tcp_ip6_spec.pdst));
		break
	case AH_V6_FLOW:
	case ESP_V6_FLOW:
		if flow_type == AH_V6_FLOW {
			fmt.Printf("\tRule Type: IPSEC AH over IPv6\n")
		} else {
			fmt.Printf("\tRule Type: IPSEC ESP over IPv6\n")
		}
		// rxclass_print_ipv6_rule(fsp.h_u.ah_ip6_spec.ip6src,
		// 	fsp.m_u.ah_ip6_spec.ip6src,
		// 	fsp.h_u.ah_ip6_spec.ip6dst,
		// 	fsp.m_u.ah_ip6_spec.ip6dst,
		// 	fsp.h_u.ah_ip6_spec.tclass,
		// 	fsp.m_u.ah_ip6_spec.tclass)
		// fmt.Printf(
		// 	"\tSPI: %d mask: 0x%x\n",
		// 	ntohl(fsp.h_u.esp_ip6_spec.spi),
		// 	ntohl(fsp.m_u.esp_ip6_spec.spi))
		break
	case IPV6_USER_FLOW:
		fmt.Printf("\tRule Type: Raw IPv6\n")
		// rxclass_print_ipv6_rule(fsp.h_u.usr_ip6_spec.ip6src,
		// 	fsp.m_u.usr_ip6_spec.ip6src,
		// 	fsp.h_u.usr_ip6_spec.ip6dst,
		// 	fsp.m_u.usr_ip6_spec.ip6dst,
		// 	fsp.h_u.usr_ip6_spec.tclass,
		// 	fsp.m_u.usr_ip6_spec.tclass)
		// fmt.Printf(
		// 	"\tProtocol: %d mask: 0x%x\n"+
		// 		"\tL4 bytes: 0x%x mask: 0x%x\n",
		// 	fsp.h_u.usr_ip6_spec.l4_proto,
		// 	fsp.m_u.usr_ip6_spec.l4_proto,
		// 	ntohl(fsp.h_u.usr_ip6_spec.l4_4_bytes),
		// 	ntohl(fsp.m_u.usr_ip6_spec.l4_4_bytes))
		break
	case ETHER_FLOW:
		// dmac := fsp.h_u.ether_spec.h_dest
		// dmacm := fsp.m_u.ether_spec.h_dest
		// smac := fsp.h_u.ether_spec.h_source
		// smacm := fsp.m_u.ether_spec.h_source

		// fmt.Printf(
		// 	"\tFlow Type: Raw Ethernet\n"+
		// 		"\tSrc MAC addr: %02X:%02X:%02X:%02X:%02X:%02X"+
		// 		" mask: %02X:%02X:%02X:%02X:%02X:%02X\n"+
		// 		"\tDest MAC addr: %02X:%02X:%02X:%02X:%02X:%02X"+
		// 		" mask: %02X:%02X:%02X:%02X:%02X:%02X\n"+
		// 		"\tEthertype: 0x%X mask: 0x%X\n",
		// 	smac[0], smac[1], smac[2], smac[3], smac[4], smac[5],
		// 	smacm[0], smacm[1], smacm[2], smacm[3], smacm[4],
		// 	smacm[5], dmac[0], dmac[1], dmac[2], dmac[3], dmac[4],
		// 	dmac[5], dmacm[0], dmacm[1], dmacm[2], dmacm[3],
		// 	dmacm[4], dmacm[5],
		// 	ntohs(fsp.h_u.ether_spec.h_proto),
		// 	ntohs(fsp.m_u.ether_spec.h_proto))
		break
	default:
		fmt.Printf("\tUnknown Flow type: %d\n", flow_type)
		break
	}

	rxclass_print_nfc_spec_ext(fsp)

	if fsp.flow_type&FLOW_RSS != 0 {
		fmt.Printf("\tRSS Context ID: %d\n", rss_context)
	}
	if fsp.ring_cookie == RX_CLS_FLOW_DISC {
		fmt.Printf("\tAction: Drop\n")
	} else if fsp.ring_cookie == RX_CLS_FLOW_WAKE {
		fmt.Printf("\tAction: Wake-on-LAN\n")
	} else {
		// vf := ethtool_get_flow_spec_ring_vf(fsp.ring_cookie)
		// queue := ethtool_get_flow_spec_ring(fsp.ring_cookie)

		// /* A value of zero indicates that this rule targeted the main
		//  * function. A positive value indicates which virtual function
		//  * was targeted, so we'll subtract 1 in order to show the
		//  * correct VF index
		//  */
		// if vf {
		// 	fmt.Printf("\tAction: Direct to VF %lld queue %llu\n",
		// 		vf-1, queue)
		// } else {
		// 	fmt.Printf("\tAction: Direct to queue %lld\n", queue)
		// }
	}

	fmt.Printf("\n")
}

func rxclass_print_rule(fsp *ethtool_rx_flow_spec, rss_context uint32) {
	/* print the rule in this location */
	switch fsp.flow_type & ^(uint32(FLOW_EXT) | FLOW_MAC_EXT | FLOW_RSS) {
	case TCP_V4_FLOW:
	case UDP_V4_FLOW:
	case SCTP_V4_FLOW:
	case AH_V4_FLOW:
	case ESP_V4_FLOW:
	case TCP_V6_FLOW:
	case UDP_V6_FLOW:
	case SCTP_V6_FLOW:
	case AH_V6_FLOW:
	case ESP_V6_FLOW:
	case IPV6_USER_FLOW:
	case ETHER_FLOW:
		rxclass_print_nfc_rule(fsp, rss_context)
		break
	case IPV4_USER_FLOW:
		// if fsp.h_u.usr_ip4_spec.ip_ver == ETH_RX_NFC_IP4 {
		// 	rxclass_print_nfc_rule(fsp, rss_context)
		// } else { /* IPv6 uses IPV6_USER_FLOW */
		// 	fmt.Printf("IPV4_USER_FLOW with wrong ip_ver\n")
		// }
		break
	default:
		fmt.Printf("rxclass: Unknown flow type\n")
		break
	}
}

func rxclass_get_dev_info(ctx *cmd_context, count *uint32,
	driver_select *int) error {

	nfccmd := ethtool_rxnfc{
		cmd:  ETHTOOL_GRXCLSRLCNT,
		data: 0,
	}
	err := send_ioctl(ctx, uintptr(unsafe.Pointer(&nfccmd)))
	*count = nfccmd.rule_cnt
	if driver_select != nil {
		*driver_select = int(nfccmd.data & RX_CLS_LOC_SPECIAL)
	}
	if err != nil {
		fmt.Printf("rxclass: Cannot get RX class rule count")
	}
	return err
}

func rxclass_rule_get(ctx *cmd_context, loc uint32) error {

	/* fetch rule from netdev */
	nfccmd := ethtool_rxnfc{cmd: ETHTOOL_GRXCLSRULE}
	// memset(&nfccmd.fs, 0, sizeof(struct ethtool_rx_flow_spec));
	nfccmd.fs.location = loc
	err := send_ioctl(ctx, uintptr(unsafe.Pointer(&nfccmd)))
	if err != nil {
		fmt.Printf("rxclass: Cannot get RX class rule")
		return err
	}

	/* display rule */
	rxclass_print_rule(&nfccmd.fs, nfccmd.rule_cnt)
	return err
}

func rxclass_rule_getall(ctx *cmd_context) error {
	var count uint32
	/* determine rule count */
	err := rxclass_get_dev_info(ctx, &count, nil)
	if err != nil {
		return err
	}

	fmt.Printf("Total %d rules\n\n", count)

	/* request location list */
	nfccmd := ethtool_rxnfc{
		cmd:      ETHTOOL_GRXCLSRLALL,
		rule_cnt: count,
	}
	err = send_ioctl(ctx, uintptr(unsafe.Pointer(&nfccmd)))
	if err != nil {
		fmt.Printf("rxclass: Cannot get RX class rules")
		return err
	}

	/* write locations to bitmap */
	rule_locs := nfccmd.rule_locs
	for i := uint32(0); i < count; i++ {
		err = rxclass_rule_get(ctx, rule_locs[i])
		if err != nil {
			break
		}
	}

	return err
}
