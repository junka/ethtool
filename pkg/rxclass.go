package ethtool

import (
	"encoding/binary"
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

func rxclass_print_ipv6_rule(sip [16]byte, sipm [16]byte, dip [16]byte,
	dipm [16]byte, tclass uint8, tclassm uint8) {

	fmt.Printf(
		"\tSrc IP addr: %s mask: %s\n"+
			"\tDest IP addr: %s mask: %s\n"+
			"\tTraffic Class: 0x%x mask: 0x%x\n",
		net.IP(sip[:]).String(), net.IP(sipm[:]).String(),
		net.IP(dip[:]).String(), net.IP(dipm[:]).String(),
		tclass, tclassm)
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
		fallthrough
	case UDP_V4_FLOW:
		fallthrough
	case SCTP_V4_FLOW:
		if flow_type == TCP_V4_FLOW {
			fmt.Printf("\tRule Type: TCP over IPv4\n")
		} else if flow_type == UDP_V4_FLOW {
			fmt.Printf("\tRule Type: UDP over IPv4\n")
		} else {
			fmt.Printf("\tRule Type: SCTP over IPv4\n")
		}
		tcp_ip4_spec := ethtool_tcpip4_spec{
			ip4src: binary.LittleEndian.Uint32(fsp.h_u.hdata[0:3]),
			ip4dst: binary.LittleEndian.Uint32(fsp.h_u.hdata[4:7]),
			psrc:   binary.LittleEndian.Uint16(fsp.h_u.hdata[8:9]),
			pdst:   binary.LittleEndian.Uint16(fsp.h_u.hdata[10:11]),
			tos:    fsp.h_u.hdata[12],
		}
		tcp_ip4_mask := ethtool_tcpip4_spec{
			ip4src: binary.LittleEndian.Uint32(fsp.m_u.hdata[0:3]),
			ip4dst: binary.LittleEndian.Uint32(fsp.m_u.hdata[4:7]),
			psrc:   binary.LittleEndian.Uint16(fsp.m_u.hdata[8:9]),
			pdst:   binary.LittleEndian.Uint16(fsp.m_u.hdata[10:11]),
			tos:    fsp.m_u.hdata[12],
		}
		rxclass_print_ipv4_rule(tcp_ip4_spec.ip4src,
			tcp_ip4_mask.ip4src,
			tcp_ip4_spec.ip4dst,
			tcp_ip4_mask.ip4dst,
			tcp_ip4_spec.tos,
			tcp_ip4_mask.tos)
		fmt.Printf(
			"\tSrc port: %d mask: 0x%x\n"+
				"\tDest port: %d mask: 0x%x\n",
			tcp_ip4_spec.psrc,
			tcp_ip4_mask.psrc,
			tcp_ip4_spec.pdst,
			tcp_ip4_mask.pdst)

	case AH_V4_FLOW:
		fallthrough
	case ESP_V4_FLOW:
		if flow_type == AH_V4_FLOW {
			fmt.Printf("\tRule Type: IPSEC AH over IPv4\n")
		} else {
			fmt.Printf("\tRule Type: IPSEC ESP over IPv4\n")
		}
		ah_ip4_spec := ethtool_ah_espip4_spec{
			ip4src: binary.LittleEndian.Uint32(fsp.h_u.hdata[0:3]),
			ip4dst: binary.LittleEndian.Uint32(fsp.h_u.hdata[4:7]),
			spi:    binary.LittleEndian.Uint32(fsp.h_u.hdata[8:11]),
			tos:    fsp.h_u.hdata[12],
		}
		ah_ip4_mask := ethtool_ah_espip4_spec{
			ip4src: binary.LittleEndian.Uint32(fsp.m_u.hdata[0:3]),
			ip4dst: binary.LittleEndian.Uint32(fsp.m_u.hdata[4:7]),
			spi:    binary.LittleEndian.Uint32(fsp.m_u.hdata[8:11]),
			tos:    fsp.m_u.hdata[12],
		}
		rxclass_print_ipv4_rule(ah_ip4_spec.ip4src,
			ah_ip4_mask.ip4src,
			ah_ip4_spec.ip4dst,
			ah_ip4_mask.ip4dst,
			ah_ip4_spec.tos,
			ah_ip4_mask.tos)

		fmt.Printf("\tSPI: %d mask: 0x%x\n",
			ah_ip4_spec.spi,
			ah_ip4_mask.spi)

	case IPV4_USER_FLOW:
		fmt.Printf("\tRule Type: Raw IPv4\n")
		usr_ip4_spec := ethtool_usrip4_spec{
			ip4src:     binary.LittleEndian.Uint32(fsp.h_u.hdata[0:3]),
			ip4dst:     binary.LittleEndian.Uint32(fsp.h_u.hdata[4:7]),
			l4_4_bytes: binary.LittleEndian.Uint32(fsp.h_u.hdata[8:11]),
			tos:        fsp.h_u.hdata[12],
			ip_ver:     fsp.h_u.hdata[13],
			proto:      fsp.h_u.hdata[14],
		}
		usr_ip4_mask := ethtool_usrip4_spec{
			ip4src:     binary.LittleEndian.Uint32(fsp.m_u.hdata[0:3]),
			ip4dst:     binary.LittleEndian.Uint32(fsp.m_u.hdata[4:7]),
			l4_4_bytes: binary.LittleEndian.Uint32(fsp.m_u.hdata[8:11]),
			tos:        fsp.m_u.hdata[12],
			ip_ver:     fsp.m_u.hdata[13],
			proto:      fsp.m_u.hdata[14],
		}
		rxclass_print_ipv4_rule(usr_ip4_spec.ip4src,
			usr_ip4_mask.ip4src,
			usr_ip4_spec.ip4dst,
			usr_ip4_mask.ip4dst,
			usr_ip4_spec.tos,
			usr_ip4_mask.tos)
		fmt.Printf(
			"\tProtocol: %d mask: 0x%x\n"+
				"\tL4 bytes: 0x%x mask: 0x%x\n",
			usr_ip4_spec.proto,
			usr_ip4_mask.proto,
			usr_ip4_spec.l4_4_bytes,
			usr_ip4_mask.l4_4_bytes)

	case TCP_V6_FLOW:
		fallthrough
	case UDP_V6_FLOW:
		fallthrough
	case SCTP_V6_FLOW:
		if flow_type == TCP_V6_FLOW {
			fmt.Printf("\tRule Type: TCP over IPv6\n")
		} else if flow_type == UDP_V6_FLOW {
			fmt.Printf("\tRule Type: UDP over IPv6\n")
		} else {
			fmt.Printf("\tRule Type: SCTP over IPv6\n")
		}
		tcp_ip6_spec := ethtool_tcpip6_spec{
			psrc:   binary.LittleEndian.Uint16(fsp.h_u.hdata[32:33]),
			pdst:   binary.LittleEndian.Uint16(fsp.h_u.hdata[34:35]),
			tclass: fsp.h_u.hdata[36],
		}
		copy(tcp_ip6_spec.ip6src[:], fsp.h_u.hdata[0:15])
		copy(tcp_ip6_spec.ip6dst[:], fsp.h_u.hdata[16:31])
		tcp_ip6_mask := ethtool_tcpip6_spec{
			psrc:   binary.LittleEndian.Uint16(fsp.m_u.hdata[32:33]),
			pdst:   binary.LittleEndian.Uint16(fsp.m_u.hdata[34:35]),
			tclass: fsp.m_u.hdata[36],
		}
		copy(tcp_ip6_mask.ip6src[:], fsp.m_u.hdata[0:15])
		copy(tcp_ip6_mask.ip6dst[:], fsp.m_u.hdata[16:31])
		rxclass_print_ipv6_rule(tcp_ip6_spec.ip6src,
			tcp_ip6_mask.ip6src,
			tcp_ip6_spec.ip6dst,
			tcp_ip6_mask.ip6dst,
			tcp_ip6_spec.tclass,
			tcp_ip6_mask.tclass)
		fmt.Printf(
			"\tSrc port: %d mask: 0x%x\n"+
				"\tDest port: %d mask: 0x%x\n",
			tcp_ip6_spec.psrc,
			tcp_ip6_mask.psrc,
			tcp_ip6_spec.pdst,
			tcp_ip6_mask.pdst)

	case AH_V6_FLOW:
		fallthrough
	case ESP_V6_FLOW:
		if flow_type == AH_V6_FLOW {
			fmt.Printf("\tRule Type: IPSEC AH over IPv6\n")
		} else {
			fmt.Printf("\tRule Type: IPSEC ESP over IPv6\n")
		}
		ah_ip6_spec := ethtool_ah_espip6_spec{
			spi:    binary.LittleEndian.Uint32(fsp.h_u.hdata[32:35]),
			tclass: fsp.h_u.hdata[36],
		}
		copy(ah_ip6_spec.ip6src[:], fsp.h_u.hdata[0:15])
		copy(ah_ip6_spec.ip6dst[:], fsp.h_u.hdata[16:31])
		ah_ip6_mask := ethtool_ah_espip6_spec{

			spi:    binary.LittleEndian.Uint32(fsp.m_u.hdata[32:35]),
			tclass: fsp.m_u.hdata[36],
		}
		copy(ah_ip6_mask.ip6src[:], fsp.m_u.hdata[0:15])
		copy(ah_ip6_mask.ip6dst[:], fsp.m_u.hdata[16:31])
		rxclass_print_ipv6_rule(ah_ip6_spec.ip6src,
			ah_ip6_mask.ip6src,
			ah_ip6_spec.ip6dst,
			ah_ip6_mask.ip6dst,
			ah_ip6_spec.tclass,
			ah_ip6_mask.tclass)
		fmt.Printf("\tSPI: %d mask: 0x%x\n",
			ah_ip6_spec.spi,
			ah_ip6_mask.spi)

	case IPV6_USER_FLOW:
		fmt.Printf("\tRule Type: Raw IPv6\n")
		usr_ip6_spec := ethtool_usrip6_spec{

			l4_4_bytes: binary.LittleEndian.Uint32(fsp.h_u.hdata[32:35]),
			tclass:     fsp.h_u.hdata[36],
			l4_proto:   fsp.h_u.hdata[37],
		}
		copy(usr_ip6_spec.ip6src[:], fsp.h_u.hdata[0:15])
		copy(usr_ip6_spec.ip6dst[:], fsp.h_u.hdata[16:31])

		usr_ip6_mask := ethtool_usrip6_spec{

			l4_4_bytes: binary.LittleEndian.Uint32(fsp.m_u.hdata[32:35]),
			tclass:     fsp.m_u.hdata[36],
			l4_proto:   fsp.m_u.hdata[37],
		}

		copy(usr_ip6_mask.ip6src[:], fsp.m_u.hdata[0:15])
		copy(usr_ip6_mask.ip6dst[:], fsp.m_u.hdata[16:31])
		rxclass_print_ipv6_rule(usr_ip6_spec.ip6src,
			usr_ip6_mask.ip6src,
			usr_ip6_spec.ip6dst,
			usr_ip6_mask.ip6dst,
			usr_ip6_spec.tclass,
			usr_ip6_mask.tclass)
		fmt.Printf("\tProtocol: %d mask: 0x%x\n"+
			"\tL4 bytes: 0x%x mask: 0x%x\n",
			usr_ip6_spec.l4_proto,
			usr_ip6_mask.l4_proto,
			(usr_ip6_spec.l4_4_bytes),
			(usr_ip6_mask.l4_4_bytes))

	case ETHER_FLOW:
		dmac := [6]byte{fsp.h_u.hdata[0], fsp.h_u.hdata[1], fsp.h_u.hdata[2],
			fsp.h_u.hdata[3], fsp.h_u.hdata[4], fsp.h_u.hdata[5]}
		dmacm := [6]byte{fsp.m_u.hdata[0], fsp.m_u.hdata[1], fsp.m_u.hdata[2],
			fsp.m_u.hdata[3], fsp.m_u.hdata[4], fsp.h_u.hdata[5]}
		smac := [6]byte{fsp.h_u.hdata[6], fsp.h_u.hdata[7], fsp.h_u.hdata[8],
			fsp.h_u.hdata[9], fsp.h_u.hdata[10], fsp.h_u.hdata[11]}
		smacm := [6]byte{fsp.m_u.hdata[6], fsp.m_u.hdata[7], fsp.m_u.hdata[8],
			fsp.m_u.hdata[9], fsp.m_u.hdata[10], fsp.m_u.hdata[11]}
		proto := binary.LittleEndian.Uint32(fsp.h_u.hdata[12:13])
		protom := binary.LittleEndian.Uint32(fsp.m_u.hdata[12:13])
		fmt.Printf(
			"\tFlow Type: Raw Ethernet\n"+
				"\tSrc MAC addr: %02X:%02X:%02X:%02X:%02X:%02X"+
				" mask: %02X:%02X:%02X:%02X:%02X:%02X\n"+
				"\tDest MAC addr: %02X:%02X:%02X:%02X:%02X:%02X"+
				" mask: %02X:%02X:%02X:%02X:%02X:%02X\n"+
				"\tEthertype: 0x%X mask: 0x%X\n",
			smac[0], smac[1], smac[2], smac[3], smac[4], smac[5],
			smacm[0], smacm[1], smacm[2], smacm[3], smacm[4],
			smacm[5], dmac[0], dmac[1], dmac[2], dmac[3], dmac[4],
			dmac[5], dmacm[0], dmacm[1], dmacm[2], dmacm[3],
			dmacm[4], dmacm[5],
			proto,
			protom)

	default:
		fmt.Printf("\tUnknown Flow type: %d\n", flow_type)

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
		vf := ethtool_get_flow_spec_ring_vf(fsp.ring_cookie)
		queue := ethtool_get_flow_spec_ring(fsp.ring_cookie)

		/* A value of zero indicates that this rule targeted the main
		 * function. A positive value indicates which virtual function
		 * was targeted, so we'll subtract 1 in order to show the
		 * correct VF index
		 */
		if vf != 0 {
			fmt.Printf("\tAction: Direct to VF %d queue %d\n", vf-1, queue)
		} else {
			fmt.Printf("\tAction: Direct to queue %d\n", queue)
		}
	}

	fmt.Printf("\n")
}

func rxclass_print_rule(fsp *ethtool_rx_flow_spec, rss_context uint32) {
	/* print the rule in this location */
	switch fsp.flow_type & ^(uint32(FLOW_EXT) | FLOW_MAC_EXT | FLOW_RSS) {
	case TCP_V4_FLOW:
		fallthrough
	case UDP_V4_FLOW:
		fallthrough
	case SCTP_V4_FLOW:
		fallthrough
	case AH_V4_FLOW:
		fallthrough
	case ESP_V4_FLOW:
		fallthrough
	case TCP_V6_FLOW:
		fallthrough
	case UDP_V6_FLOW:
		fallthrough
	case SCTP_V6_FLOW:
		fallthrough
	case AH_V6_FLOW:
		fallthrough
	case ESP_V6_FLOW:
		fallthrough
	case IPV6_USER_FLOW:
		fallthrough
	case ETHER_FLOW:
		rxclass_print_nfc_rule(fsp, rss_context)

	case IPV4_USER_FLOW:
		usr_ip4_spec := ethtool_usrip4_spec{
			ip4src:     binary.LittleEndian.Uint32(fsp.h_u.hdata[0:3]),
			ip4dst:     binary.LittleEndian.Uint32(fsp.h_u.hdata[4:7]),
			l4_4_bytes: binary.LittleEndian.Uint32(fsp.h_u.hdata[8:11]),
			tos:        fsp.h_u.hdata[12],
			ip_ver:     fsp.h_u.hdata[13],
			proto:      fsp.h_u.hdata[14],
		}
		if usr_ip4_spec.ip_ver == ETH_RX_NFC_IP4 {
			rxclass_print_nfc_rule(fsp, rss_context)
		} else { /* IPv6 uses IPV6_USER_FLOW */
			fmt.Printf("IPV4_USER_FLOW with wrong ip_ver\n")
		}

	default:
		fmt.Printf("rxclass: Unknown flow type\n")
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
