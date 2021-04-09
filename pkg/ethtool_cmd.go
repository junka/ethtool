package ethtool

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/junka/ioctl"
)

const (
	CMDL_NONE = iota
	CMDL_BOOL
	CMDL_S32
	CMDL_U8
	CMDL_U16
	CMDL_U32
	CMDL_U64
	CMDL_BE16
	CMDL_IP4
	CMDL_STR
	CMDL_FLAG
	CMDL_MAC
)

type cmdline_info struct {
	name string
	tp   int
	/* Points to int (BOOL), s32, u16, u32 (U32/FLAG/IP4), u64,
	 * char * (STR) or u8[6] (MAC).  For FLAG, the value accumulates
	 * all flags to be set. */
	wanted_val uintptr
	ioctl_val  uintptr
	/* For FLAG, the flag value to be set/cleared */
	flag_val uint32
	/* For FLAG, points to u32 and accumulates all flags seen.
	 * For anything else, points to int and is set if the option is
	 * seen. */
	seen_val uintptr
}

type feature_def struct {
	name           [ETH_GSTRING_LEN]byte
	off_flag_index int /* index in off_flag_def; negative if none match */
}

type feature_defs struct {
	n_features uint64
	/* Number of features each offload flag is associated with */
	off_flag_matched [OFF_FLAG_DEF_SIZE]uint64
	/* Name and offload flag index for each feature */
	def [0]feature_def
}

const SIOCETHTOOL = 0x8946

func send_ioctl(ctx *cmd_context, cmd uintptr) error {
	ctx.ifr.ifr_data = cmd

	return ioctl.Ioctl(ctx.fd, SIOCETHTOOL, uintptr(unsafe.Pointer(&ctx.ifr)))
}

func init_ioctl(ctx *cmd_context, no_dev bool) int {
	if no_dev {
		ctx.fd = -1
		return 0
	}
	if len(ctx.devname) > IF_NAMESIZE {
		fmt.Println("Device name longer than %u characters")
		return -1
	}
	copy(ctx.ifr.ifr_name[:], ctx.devname)
	var err error
	ctx.fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if ctx.fd < 0 || err != nil {
		ctx.fd, err = syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW, 16)
	}
	if ctx.fd < 0 || err != nil {
		fmt.Printf("cannot get control socket")
		return 70
	}
	return 0
}

//
func dump_drvinfo(info *ethtool_drvinfo) int {
	stats_sp, test_info_sp, eedump_sp, regdump_sp, priv_sp := "no", "no", "no", "no", "no"
	if info.n_stats > 0 {
		stats_sp = "yes"
	}
	if info.testinfo_len > 0 {
		test_info_sp = "yes"
	}
	if info.eedump_len > 0 {
		eedump_sp = "yes"
	}
	if info.regdump_len > 0 {
		regdump_sp = "yes"
	}
	if info.n_priv_flags > 0 {
		priv_sp = "yes"
	}
	fmt.Printf("driver: %.*s\n"+
		"version: %.*s\n"+
		"firmware-version: %.*s\n"+
		"expansion-rom-version: %.*s\n"+
		"bus-info: %.*s\n"+
		"supports-statistics: %s\n"+
		"supports-test: %s\n"+
		"supports-eeprom-access: %s\n"+
		"supports-register-dump: %s\n"+
		"supports-priv-flags: %s\n",
		len(string(info.driver[:])), info.driver,
		len(info.version), info.version,
		len(info.fw_version), info.fw_version,
		len(info.erom_version), info.erom_version,
		len(info.bus_info), info.bus_info, stats_sp,
		test_info_sp, eedump_sp, regdump_sp, priv_sp)

	return 0
}

type driver_dump struct {
	name       string
	regdump_fn func(info *ethtool_drvinfo, regs *ethtool_regs)
}

var driver_list = []driver_dump{
	// {"8139cp", realtek_dump_regs},
	// {"8139too", realtek_dump_regs},
	// {"r8169", realtek_dump_regs},
	// {"de2104x", de2104x_dump_regs},
	// {"e1000", e1000_dump_regs},
	// {"e1000e", e1000_dump_regs},
	// {"igb", igb_dump_regs},
	// {"ixgb", ixgb_dump_regs},
	// {"ixgbe", ixgbe_dump_regs},
	// {"ixgbevf", ixgbevf_dump_regs},
	// {"natsemi", natsemi_dump_regs},
	// {"e100", e100_dump_regs},
	// {"amd8111e", amd8111e_dump_regs},
	// {"pcnet32", pcnet32_dump_regs},
	// {"fec_8xx", fec_8xx_dump_regs},
	// {"ibm_emac", ibm_emac_dump_regs},
	// {"tg3", tg3_dump_regs},
	// {"skge", skge_dump_regs},
	// {"sky2", sky2_dump_regs},
	// {"vioc", vioc_dump_regs},
	// {"smsc911x", smsc911x_dump_regs},
	// {"at76c50x-usb", at76c50x_usb_dump_regs},
	// {"sfc", sfc_dump_regs},
	// {"st_mac100", st_mac100_dump_regs},
	// {"st_gmac", st_gmac_dump_regs},
	// {"et131x", et131x_dump_regs},
	// {"altera_tse", altera_tse_dump_regs},
	// {"vmxnet3", vmxnet3_dump_regs},
	// {"fjes", fjes_dump_regs},
	// {"lan78xx", lan78xx_dump_regs},
	// {"dsa", dsa_dump_regs},
	// {"fec", fec_dump_regs},
	// {"igc", igc_dump_regs},
	// {"bnxt_en", bnxt_dump_regs},
}

func dump_hex(file *os.File, data []uint8, len uint32, offset uint32) {

	fmt.Fprintf(file, "Offset\t\tValues\n")
	fmt.Fprintf(file, "------\t\t------")
	for i := uint32(0); i < len; i++ {
		if i%16 == 0 {
			fmt.Fprintf(file, "\n0x%04x:\t\t", i+offset)
		}
		fmt.Fprintf(file, "%02x ", data[i])
	}
	fmt.Fprintf(file, "\n")
}

func dump_regs(gregs_dump_raw int, gregs_dump_hex int,
	info *ethtool_drvinfo, regs *ethtool_regs) int {

	// if gregs_dump_raw != 0 {
	// 	// fwrite(regs.data, regs.len, 1, stdout);
	// 	if info.regdump_len > regs.len+uint32(unsafe.Sizeof(*info)+unsafe.Sizeof(*regs)) {
	// 		info = (*ethtool_drvinfo)(&regs.data[regs.len])
	// 		regs = (*ethtool_regs)(&regs.data[regs.len+uint32(unsafe.Sizeof(*info))])

	// 		return dump_regs(gregs_dump_raw, gregs_dump_hex, info, regs)
	// 	}
	// 	return 0
	// }

	// if gregs_dump_hex == 0 {
	// 	for i := 0; i < len(driver_list); i++ {
	// 		if driver_list[i].name == string(info.driver[:]) {
	// 			if driver_list[i].regdump_fn(info, regs) == 0 {
	// 				/* Recurse dump if some drvinfo and regs structures are nested */
	// 				if info.regdump_len > regs.len+unsafe.Sizeof(*info)+unsafe.Sizeof(*regs) {
	// 					info = (&regs.data[regs.len])
	// 					regs = (&regs.data[regs.len+uint32(unsafe.Sizeof(*info))])

	// 					return dump_regs(gregs_dump_raw, gregs_dump_hex, info, regs)
	// 				}
	// 			}
	// 			/* This version (or some other
	// 			 * variation in the dump format) is
	// 			 * not handled; fall back to hex
	// 			 */
	// 			break
	// 		}
	// 	}
	// }

	// dump_hex(os.Stdout, regs.data, regs.len, 0)

	// /* Recurse dump if some drvinfo and regs structures are nested */
	// if info.regdump_len > regs.len+unsafe.Sizeof(*info)+unsafe.Sizeof(*regs) {
	// 	info = (&regs.data[0] + regs.len)
	// 	regs = (&regs.data[0] + regs.len + unsafe.Sizeof(*info))

	// 	return dump_regs(gregs_dump_raw, gregs_dump_hex, info, regs)
	// }

	return 0
}

func dump_eeprom(geeprom_dump_raw int,
	info *ethtool_drvinfo,
	ee *ethtool_eeprom) int {
	if geeprom_dump_raw != 0 {
		// fwrite(ee.data, 1, ee.len, stdout);
		return 0
	}
	// #ifdef ETHTOOL_ENABLE_PRETTY_DUMP
	// 	if (!strncmp("natsemi", info.driver, ETHTOOL_BUSINFO_LEN)) {
	// 		return natsemi_dump_eeprom(info, ee);
	// 	} else if (!strncmp("tg3", info.driver, ETHTOOL_BUSINFO_LEN)) {
	// 		return tg3_dump_eeprom(info, ee);
	// 	}
	// #endif
	// (*[4]uint8)(unsafe.Pointer(&ee.data[]))[:]
	// dump_hex(os.Stdout, ee.data, ee.len, ee.offset)

	return 0
}

func dump_test(test *ethtool_test,
	strings *ethtool_gstrings) int {
	res := "PASS"
	rc := test.flags & ETH_TEST_FL_FAILED
	if rc != 0 {
		res = "FAIL"
	}
	fmt.Printf("The test result is %s\n", res)

	res = "not "
	if test.flags&ETH_TEST_FL_EXTERNAL_LB != 0 {
		res = ""
	}
	fmt.Printf("External loopback test was %sexecuted\n", res)

	if strings.len > 0 {
		fmt.Printf("The test extra info:\n")
	}
	for i := 0; i < int(strings.len); i++ {
		// fmt.Printf("%s\t %d\n",
		// 	string(&strings.data[i*ETH_GSTRING_LEN]),
		// 	test.data[i])
	}

	fmt.Printf("\n")
	return int(rc)
}

func dump_pause(epause *ethtool_pauseparam,
	advertising uint32, lp_advertising uint32) int {
	neg, rx, tx, arx, atx := "off", "off", "off", "off", "off"
	if epause.autoneg > 0 {
		neg = "on"
	}
	if epause.rx_pause > 0 {
		rx = "on"
	}
	if epause.tx_pause > 0 {
		tx = "on"
	}
	fmt.Printf(
		"Autonegotiate:	%s\n"+
			"RX:		%s\n"+
			"TX:		%s\n",
		neg, rx, tx)

	if lp_advertising != 0 {
		an_rx, an_tx := 0, 0

		/* Work out negotiated pause frame usage per
		 * IEEE 802.3-2005 table 28B-3.
		 */
		if (advertising & lp_advertising & 1 << ETHTOOL_LINK_MODE_Pause_BIT) != 0 {
			an_tx = 1
			an_rx = 1
		} else if (advertising & lp_advertising &
			1 << ETHTOOL_LINK_MODE_Asym_Pause_BIT) != 0 {
			if (advertising & 1 << ETHTOOL_LINK_MODE_Pause_BIT) != 0 {
				an_rx = 1
			} else if (lp_advertising & 1 << ETHTOOL_LINK_MODE_Pause_BIT) != 0 {
				an_tx = 1
			}
		}
		if an_rx == 1 {
			arx = "on"
		}
		if an_tx == 1 {
			atx = "on"
		}
		fmt.Printf(
			"RX negotiated:	%s\n"+
				"TX negotiated:	%s\n",
			arx, atx)
	}

	fmt.Printf("\n")
	return 0
}

func dump_ring(ering *ethtool_ringparam) int {
	fmt.Printf(
		"Pre-set maximums:\n"+
			"RX:		%d\n"+
			"RX Mini:	%d\n"+
			"RX Jumbo:	%d\n"+
			"TX:		%d\n",
		ering.rx_max_pending,
		ering.rx_mini_max_pending,
		ering.rx_jumbo_max_pending,
		ering.tx_max_pending)

	fmt.Printf(
		"Current hardware settings:\n"+
			"RX:		%d\n"+
			"RX Mini:	%d\n"+
			"RX Jumbo:	%d\n"+
			"TX:		%d\n",
		ering.rx_pending,
		ering.rx_mini_pending,
		ering.rx_jumbo_pending,
		ering.tx_pending)

	fmt.Printf("\n")
	return 0
}

func dump_channels(echannels *ethtool_channels) int {
	fmt.Printf(
		"Pre-set maximums:\n"+
			"RX:		%d\n"+
			"TX:		%d\n"+
			"Other:		%d\n"+
			"Combined:	%d\n",
		echannels.max_rx, echannels.max_tx,
		echannels.max_other,
		echannels.max_combined)

	fmt.Printf(
		"Current hardware settings:\n"+
			"RX:		%d\n"+
			"TX:		%d\n"+
			"Other:		%d\n"+
			"Combined:	%d\n",
		echannels.rx_count, echannels.tx_count,
		echannels.other_count,
		echannels.combined_count)

	fmt.Printf("\n")
	return 0
}

func dump_coalesce(ecoal *ethtool_coalesce) int {
	urx, utx := "off", "off"
	if ecoal.use_adaptive_rx_coalesce != 0 {
		urx = "on"
	}
	if ecoal.use_adaptive_tx_coalesce != 0 {
		utx = "on"
	}
	fmt.Printf("Adaptive RX: %s  TX: %s\n", urx, utx)

	fmt.Printf(
		"stats-block-usecs: %d\n"+
			"sample-interval: %d\n"+
			"pkt-rate-low: %d\n"+
			"pkt-rate-high: %d\n"+
			"\n"+
			"rx-usecs: %d\n"+
			"rx-frames: %d\n"+
			"rx-usecs-irq: %d\n"+
			"rx-frames-irq: %d\n"+
			"\n"+
			"tx-usecs: %d\n"+
			"tx-frames: %d\n"+
			"tx-usecs-irq: %d\n"+
			"tx-frames-irq: %d\n"+
			"\n"+
			"rx-usecs-low: %d\n"+
			"rx-frames-low: %d\n"+
			"tx-usecs-low: %d\n"+
			"tx-frames-low: %d\n"+
			"\n"+
			"rx-usecs-high: %d\n"+
			"rx-frames-high: %d\n"+
			"tx-usecs-high: %d\n"+
			"tx-frames-high: %d\n"+
			"\n",
		ecoal.stats_block_coalesce_usecs,
		ecoal.rate_sample_interval,
		ecoal.pkt_rate_low,
		ecoal.pkt_rate_high,

		ecoal.rx_coalesce_usecs,
		ecoal.rx_max_coalesced_frames,
		ecoal.rx_coalesce_usecs_irq,
		ecoal.rx_max_coalesced_frames_irq,

		ecoal.tx_coalesce_usecs,
		ecoal.tx_max_coalesced_frames,
		ecoal.tx_coalesce_usecs_irq,
		ecoal.tx_max_coalesced_frames_irq,

		ecoal.rx_coalesce_usecs_low,
		ecoal.rx_max_coalesced_frames_low,
		ecoal.tx_coalesce_usecs_low,
		ecoal.tx_max_coalesced_frames_low,

		ecoal.rx_coalesce_usecs_high,
		ecoal.rx_max_coalesced_frames_high,
		ecoal.tx_coalesce_usecs_high,
		ecoal.tx_max_coalesced_frames_high)

	return 0
}

func dump_per_queue_coalesce(per_queue_opt *ethtool_per_queue_op,
	queue_mask *uint32, n_queues int) {
	// struct ethtool_coalesce *ecoal;
	// int i, idx = 0;

	// ecoal = (struct ethtool_coalesce *)(per_queue_opt + 1);
	// for (i = 0; i < __KERNEL_DIV_ROUND_UP(MAX_NUM_QUEUE, 32); i++) {
	// 	int queue = i * 32;
	// 	__u32 mask = queue_mask[i];

	// 	while (mask > 0) {
	// 		if (mask & 0x1) {
	// 			fmt.Printf(  "Queue: %d\n", queue);
	// 			dump_coalesce(ecoal + idx);
	// 			idx++;
	// 		}
	// 		mask = mask >> 1;
	// 		queue++;
	// 	}
	// 	if (idx == n_queues)
	// 		break;
	// }
}

type feature_state struct {
	off_flags uint32
	features  ethtool_gfeatures
}

func dump_one_feature(indent string, name string,
	state *feature_state,
	ref_state *feature_state,
	index uint32) {
	// if (ref_state &&
	// 	!(FEATURE_BIT_IS_SET(state.features.features, index, active) ^
	// 	FEATURE_BIT_IS_SET(ref_state.features.features, index, active))){
	// 	return;
	// }

	// fmt.Printf("%s%s: %s%s\n",
	// indent, name,
	// FEATURE_BIT_IS_SET(state.features.features, index, active) ?
	// "on" : "off",
	// (!FEATURE_BIT_IS_SET(state.features.features, index, available)
	// || FEATURE_BIT_IS_SET(state.features.features, index,
	// 		 never_changed))
	// ? " [fixed]"
	// : (FEATURE_BIT_IS_SET(state.features.features, index, requested)
	// ^ FEATURE_BIT_IS_SET(state.features.features, index, active))
	// ? (FEATURE_BIT_IS_SET(state.features.features, index, requested)
	// ? " [requested on]" : " [requested off]")
	// : "");
}

func linux_version_code() uint32 {
	utsname := syscall.Utsname{}
	var version, patchlevel, sublevel uint32
	err := syscall.Uname(&utsname)
	if err != nil {
		return math.MaxUint32
	}
	release := ""
	for _, r := range utsname.Release {
		if r == 0 {
			break
		}
		release += string(r)
	}
	_, err = fmt.Sscanf(release, "%u.%u.%u", &version, &patchlevel, &sublevel)
	if err != nil {
		return math.MaxUint32
	}
	return version<<16 | patchlevel<<8 | sublevel
}

func dump_features(defs *feature_defs,
	state *feature_state,
	ref_state *feature_state) {
	kernel_ver := linux_version_code()
	var indent, i int

	for i = 0; i < OFF_FLAG_DEF_SIZE; i++ {
		/* Don't show features whose state is unknown on this
		 * kernel version
		 */
		if defs.off_flag_matched[i] == 0 &&
			((off_flag_def[i].get_cmd == 0 &&
				kernel_ver < off_flag_def[i].min_kernel_ver) ||
				(off_flag_def[i].get_cmd == ETHTOOL_GUFO)) {
			continue
		}

		value := off_flag_def[i].value

		/* If this offload flag matches exactly one generic
		 * feature then it's redundant to show the flag and
		 * feature states separately.  Otherwise, show the
		 * flag state first.
		 */
		off_str := "off"
		if (state.off_flags & value) != 0 {
			off_str = "on"
		}
		if defs.off_flag_matched[i] != 1 &&
			(ref_state == nil || (state.off_flags^ref_state.off_flags)&value != 0) {
			fmt.Printf("%s: %s\n", off_flag_def[i].long_name, off_str)
			indent = 1
		} else {
			indent = 0
		}
		/* Show matching features */
		for j := uint32(0); uint64(j) < defs.n_features; j++ {
			if defs.def[j].off_flag_index != i {
				continue
			}
			ind_str := ""
			if indent != 0 {
				ind_str = "\t"
			}
			if defs.off_flag_matched[i] != 1 {
				/* Show all matching feature states */
				dump_one_feature(ind_str,
					string(defs.def[j].name[:]),
					state, ref_state, j)
			} else {
				/* Show full state with the old flag name */
				dump_one_feature("", off_flag_def[i].long_name,
					state, ref_state, j)
			}
		}
	}

	/* Show all unmatched features that have non-null names */
	for j := uint32(0); uint64(j) < defs.n_features; j++ {
		if defs.def[j].off_flag_index < 0 && defs.def[j].name[0] != 0 {
			dump_one_feature("", string(defs.def[j].name[:]),
				state, ref_state, j)
		}
	}
}
func dump_rxfhash(fhash int, val uint64) int {
	switch fhash & ^FLOW_RSS {
	case TCP_V4_FLOW:
		fmt.Printf("TCP over IPV4 flows")
		break
	case UDP_V4_FLOW:
		fmt.Printf("UDP over IPV4 flows")
		break
	case SCTP_V4_FLOW:
		fmt.Printf("SCTP over IPV4 flows")
		break
	case AH_ESP_V4_FLOW:
	case AH_V4_FLOW:
	case ESP_V4_FLOW:
		fmt.Printf("IPSEC AH/ESP over IPV4 flows")
		break
	case TCP_V6_FLOW:
		fmt.Printf("TCP over IPV6 flows")
		break
	case UDP_V6_FLOW:
		fmt.Printf("UDP over IPV6 flows")
		break
	case SCTP_V6_FLOW:
		fmt.Printf("SCTP over IPV6 flows")
		break
	case AH_ESP_V6_FLOW:
	case AH_V6_FLOW:
	case ESP_V6_FLOW:
		fmt.Printf("IPSEC AH/ESP over IPV6 flows")
		break
	default:
		break
	}

	if val&RXH_DISCARD != 0 {
		fmt.Printf(" - All matching flows discarded on RX\n")
		return 0
	}
	fmt.Printf(" use these fields for computing Hash flow key:\n")

	// fmt.Printf("%s\n", unparse_rxfhashopts(val))

	return 0
}

func dump_eeecmd(ep *ethtool_eee) {
	var link_mode [127]uint32

	fmt.Printf("	EEE status: ")
	if ep.supported == 0 {
		fmt.Printf("not supported\n")
		return
	} else if ep.eee_enabled == 0 {
		fmt.Printf("disabled\n")
	} else {
		fmt.Printf("enabled - ")
		if ep.eee_active != 0 {
			fmt.Printf("active\n")
		} else {
			fmt.Printf("inactive\n")
		}
	}

	fmt.Printf("	Tx LPI:")
	if ep.tx_lpi_enabled != 0 {
		fmt.Printf(" %d (us)\n", ep.tx_lpi_timer)
	} else {
		fmt.Printf(" disabled\n")
	}
	// ethtool_link_mode_zero(link_mode)

	link_mode[0] = ep.supported
	// dump_link_caps("Supported EEE", "", link_mode, 1)

	link_mode[0] = ep.advertised
	// dump_link_caps("Advertised EEE", "", link_mode, 1)

	link_mode[0] = ep.lp_advertised
	// dump_link_caps("Link partner advertised EEE", "", link_mode, 1)
}

func dump_fec(fec uint32) {
	if (fec & ETHTOOL_FEC_NONE) != 0 {
		fmt.Printf(" None")
	}
	if (fec & ETHTOOL_FEC_AUTO) != 0 {
		fmt.Printf(" Auto")
	}

	if (fec & ETHTOOL_FEC_OFF) != 0 {
		fmt.Printf(" Off")
	}
	if (fec & ETHTOOL_FEC_BASER) != 0 {
		fmt.Printf(" BaseR")
	}
	if (fec & ETHTOOL_FEC_RS) != 0 {
		fmt.Printf(" RS")
	}
	if (fec & ETHTOOL_FEC_LLRS) != 0 {
		fmt.Printf(" LLRS")
	}
}

const N_SOTS = 7

var so_timestamping_labels = [N_SOTS]string{
	"hardware-transmit     (SOF_TIMESTAMPING_TX_HARDWARE)",
	"software-transmit     (SOF_TIMESTAMPING_TX_SOFTWARE)",
	"hardware-receive      (SOF_TIMESTAMPING_RX_HARDWARE)",
	"software-receive      (SOF_TIMESTAMPING_RX_SOFTWARE)",
	"software-system-clock (SOF_TIMESTAMPING_SOFTWARE)",
	"hardware-legacy-clock (SOF_TIMESTAMPING_SYS_HARDWARE)",
	"hardware-raw-clock    (SOF_TIMESTAMPING_RAW_HARDWARE)",
}

const N_TX_TYPES = 3

var tx_type_labels = [N_TX_TYPES]string{
	"off                   (HWTSTAMP_TX_OFF)",
	"on                    (HWTSTAMP_TX_ON)",
	"one-step-sync         (HWTSTAMP_TX_ONESTEP_SYNC)",
}

const N_RX_FILTERS = 16

var rx_filter_labels = [N_RX_FILTERS]string{
	"none                  (HWTSTAMP_FILTER_NONE)",
	"all                   (HWTSTAMP_FILTER_ALL)",
	"some                  (HWTSTAMP_FILTER_SOME)",
	"ptpv1-l4-event        (HWTSTAMP_FILTER_PTP_V1_L4_EVENT)",
	"ptpv1-l4-sync         (HWTSTAMP_FILTER_PTP_V1_L4_SYNC)",
	"ptpv1-l4-delay-req    (HWTSTAMP_FILTER_PTP_V1_L4_DELAY_REQ)",
	"ptpv2-l4-event        (HWTSTAMP_FILTER_PTP_V2_L4_EVENT)",
	"ptpv2-l4-sync         (HWTSTAMP_FILTER_PTP_V2_L4_SYNC)",
	"ptpv2-l4-delay-req    (HWTSTAMP_FILTER_PTP_V2_L4_DELAY_REQ)",
	"ptpv2-l2-event        (HWTSTAMP_FILTER_PTP_V2_L2_EVENT)",
	"ptpv2-l2-sync         (HWTSTAMP_FILTER_PTP_V2_L2_SYNC)",
	"ptpv2-l2-delay-req    (HWTSTAMP_FILTER_PTP_V2_L2_DELAY_REQ)",
	"ptpv2-event           (HWTSTAMP_FILTER_PTP_V2_EVENT)",
	"ptpv2-sync            (HWTSTAMP_FILTER_PTP_V2_SYNC)",
	"ptpv2-delay-req       (HWTSTAMP_FILTER_PTP_V2_DELAY_REQ)",
	"ntp-all               (HWTSTAMP_FILTER_NTP_ALL)",
}

func dump_tsinfo(info *ethtool_ts_info) int {

	fmt.Printf("Capabilities:\n")

	for i := 0; i < N_SOTS; i++ {
		if (info.so_timestamping & (1 << i)) != 0 {
			fmt.Printf("\t%s\n", so_timestamping_labels[i])
		}
	}

	fmt.Printf("PTP Hardware Clock: ")

	if info.phc_index < 0 {
		fmt.Printf("none\n")
	} else {
		fmt.Printf("%d\n", info.phc_index)
	}
	fmt.Printf("Hardware Transmit Timestamp Modes:")

	if info.tx_types == 0 {
		fmt.Printf(" none\n")
	} else {
		fmt.Printf("\n")
	}
	for i := 0; i < N_TX_TYPES; i++ {
		if (info.tx_types & (1 << i)) != 0 {
			fmt.Printf("\t%s\n", tx_type_labels[i])
		}
	}

	fmt.Printf("Hardware Receive Filter Modes:")

	if info.rx_filters == 0 {
		fmt.Printf(" none\n")
	} else {
		fmt.Printf("\n")
	}
	for i := 0; i < N_RX_FILTERS; i++ {
		if (info.rx_filters & (1 << i)) != 0 {
			fmt.Printf("\t%s\n", rx_filter_labels[i])
		}
	}

	return 0
}

func get_stringset(ctx *cmd_context, set_id uint32,
	drvinfo_offset uintptr, null_terminate int) *ethtool_gstrings {
	type set_info struct {
		hdr ethtool_sset_info
		buf [1]uint32
	}
	var sset_info set_info
	len := uint32(0)

	sset_info.hdr.cmd = ETHTOOL_GSSET_INFO
	sset_info.hdr.reserved = 0
	sset_info.hdr.sset_mask = (1 << set_id)
	err := send_ioctl(ctx, uintptr(unsafe.Pointer(&sset_info)))
	if err == nil {
		if sset_info.hdr.sset_mask > 0 {
			len = sset_info.buf[0]
		}
	} else if err == syscall.EOPNOTSUPP && drvinfo_offset != 0 {
		/* Fallback for old kernel versions */
		drvinfo := ethtool_drvinfo{cmd: ETHTOOL_GDRVINFO}
		if send_ioctl(ctx, uintptr(unsafe.Pointer(&drvinfo))) != nil {
			return nil
		}
		len = *(*uint32)(unsafe.Pointer((uintptr(unsafe.Pointer(&drvinfo)) + drvinfo_offset)))
	} else {
		return nil
	}
	strings := ethtool_gstrings{
		cmd:        ETHTOOL_GSTRINGS,
		string_set: set_id,
		len:        len,
	}
	if len != 0 && send_ioctl(ctx, uintptr(unsafe.Pointer(&strings))) != nil {
		return nil
	}

	if null_terminate != 0 {
		for i := uint32(0); i < len; i++ {
			strings.data[(i+1)*ETH_GSTRING_LEN-1] = 0
		}
	}
	return &strings
}

func do_gdrv(ctx *cmd_context) int {
	drvinfo := ethtool_drvinfo{
		cmd: ETHTOOL_GDRVINFO,
	}
	err := send_ioctl(ctx, uintptr(unsafe.Pointer(&drvinfo)))
	if err != nil {
		fmt.Printf("Cannot get driver information: %v\n", err)
		return 71
	}
	return dump_drvinfo(&drvinfo)
}

func do_gpause(ctx *cmd_context) int {

	epause := ethtool_pauseparam{
		cmd: ETHTOOL_GPAUSEPARAM,
	}
	err := send_ioctl(ctx, uintptr(unsafe.Pointer(&epause)))
	if err != nil {
		fmt.Printf("Cannot get device pause settings: %v\n", err)
		return 76
	}
	if epause.autoneg != 0 {
		ecmd := ethtool_cmd{cmd: ETHTOOL_GSET}
		err = send_ioctl(ctx, uintptr(unsafe.Pointer(&ecmd)))
		if err != nil {
			fmt.Printf("Cannot get device settings: %v\n", err)
			return 1
		}
		dump_pause(&epause, ecmd.advertising, ecmd.lp_advertising)
	} else {
		dump_pause(&epause, 0, 0)
	}
	return 0
}

func do_gcoalesce(ctx *cmd_context) int {
	fmt.Printf("Coalesce parameters for %s:\n", ctx.devname)

	ecoal := ethtool_coalesce{cmd: ETHTOOL_GCOALESCE}
	err := send_ioctl(ctx, uintptr(unsafe.Pointer(&ecoal)))
	if err == nil {
		dump_coalesce(&ecoal)

	} else {
		fmt.Printf("Cannot get device coalesce settings: %v\n", err)
		return 82
	}

	return 0
}

func get_features(ctx *cmd_context, defs *feature_defs) *feature_state {

	var eval ethtool_value
	allfail := 1

	state := &feature_state{}
	// malloc(sizeof(*state) +
	// 	       FEATURE_BITS_TO_BLOCKS(defs.n_features) *
	// 	       sizeof(state.features.features[0]));

	state.off_flags = 0

	for i := 0; i < OFF_FLAG_DEF_SIZE; i++ {
		value := off_flag_def[i].value
		if off_flag_def[i].get_cmd == 0 {
			continue
		}
		eval.cmd = off_flag_def[i].get_cmd
		err := send_ioctl(ctx, uintptr(unsafe.Pointer(&eval)))
		if err != nil {
			if err == syscall.EOPNOTSUPP &&
				off_flag_def[i].get_cmd == ETHTOOL_GUFO {
				continue
			}

			fmt.Printf("Cannot get device %s settings: %v\n",
				off_flag_def[i].long_name, err)
		} else {
			if eval.data != 0 {
				state.off_flags |= value
			}
			allfail = 0
		}
	}

	eval.cmd = ETHTOOL_GFLAGS
	err := send_ioctl(ctx, uintptr(unsafe.Pointer(&eval)))
	if err != nil {
		fmt.Printf("Cannot get device flags: %v\n", err)
	} else {
		state.off_flags |= eval.data & ETH_FLAG_EXT_MASK
		allfail = 0
	}

	if defs.n_features > 0 {
		state.features.cmd = ETHTOOL_GFEATURES
		state.features.size = uint32((defs.n_features + 32 - 1) / 32)
		err = send_ioctl(ctx, uintptr(unsafe.Pointer(&state.features)))
		if err != nil {
			fmt.Printf("Cannot get device generic features: %v\n", err)
		} else {
			allfail = 0
		}
	}

	if allfail != 0 {
		return nil
	}

	return state
}
func get_feature_defs(ctx *cmd_context) *feature_defs {

	defs := &feature_defs{}
	// var n_features uint32

	// names := get_stringset(ctx, ETH_SS_FEATURES, 0, 1);
	// if names!=nil {
	// 	n_features = names.len;
	// } else if (errno == EOPNOTSUPP || errno == EINVAL) {
	// 	/* Kernel doesn't support named features; not an error */
	// 	n_features = 0;
	// } else if (errno == EPERM) {
	// 	/* Kernel bug: ETHTOOL_GSSET_INFO was privileged.
	// 	 * Work around it. */
	// 	n_features = 0;
	// } else {
	// 	return NULL;
	// }

	// if (!defs) {
	// 	return NULL;
	// }

	// defs.n_features = n_features;
	// // memset(defs.off_flag_matched, 0, sizeof(defs.off_flag_matched));

	// /* Copy out feature names and find those associated with legacy flags */
	// for i := 0; i < defs.n_features; i++ {
	// 	memcpy(defs.def[i].name, names.data + i * ETH_GSTRING_LEN,
	// 	       ETH_GSTRING_LEN);
	// 	defs.def[i].off_flag_index = -1;

	// 	for j = 0;
	// 	     j < OFF_FLAG_DEF_SIZE &&
	// 		     defs.def[i].off_flag_index < 0;
	// 	     j++ {
	// 		const char *pattern =
	// 			off_flag_def[j].kernel_name;
	// 		const char *name = defs.def[i].name;
	// 		for (;;) {
	// 			if (*pattern == '*') {
	// 				/* There is only one wildcard; so
	// 				 * switch to a suffix comparison */
	// 				size_t pattern_len =
	// 					strlen(pattern + 1);
	// 				size_t name_len = strlen(name);
	// 				if (name_len < pattern_len)
	// 					break; /* name is too short */
	// 				name += name_len - pattern_len;
	// 				++pattern;
	// 			} else if (*pattern != *name) {
	// 				break; /* mismatch */
	// 			} else if (*pattern == 0) {
	// 				defs.def[i].off_flag_index = j;
	// 				defs.off_flag_matched[j]++;
	// 				break;
	// 			} else {
	// 				++name;
	// 				++pattern;
	// 			}
	// 		}
	// 	}
	// }

	return defs
}

func do_gfeatures(ctx *cmd_context) int {

	// if (ctx.argc != 0)
	// 	exit_bad_args();

	defs := get_feature_defs(ctx)
	if defs == nil {
		fmt.Printf("Cannot get device feature names\n")
		return 1
	}

	fmt.Printf("Features for %s:\n", ctx.devname)

	features := get_features(ctx, defs)
	if features == nil {
		fmt.Printf("no feature info available\n")
		return 1
	}

	dump_features(defs, features, nil)

	return 0
}

func do_gring(ctx *cmd_context) int {

	// if (ctx.argc != 0)
	// 	exit_bad_args();

	fmt.Printf("Ring parameters for %s:\n", ctx.devname)

	ering := ethtool_ringparam{
		cmd: ETHTOOL_GRINGPARAM,
	}
	err := send_ioctl(ctx, uintptr(unsafe.Pointer(&ering)))
	if err == nil {
		dump_ring(&ering)
	} else {
		fmt.Printf("Cannot get device ring settings: %v\n", err)
		return 76
	}

	return 0
}

func do_gregs(ctx *cmd_context) int {
	// int gregs_changed = 0;
	gregs_dump_raw := 1
	gregs_dump_hex := 1

	// parse_generic_cmdline(ctx, &gregs_changed,
	// 		      cmdline_gregs, ARRAY_SIZE(cmdline_gregs));

	drvinfo := ethtool_drvinfo{cmd: ETHTOOL_GDRVINFO}
	err := send_ioctl(ctx, uintptr(unsafe.Pointer(&drvinfo)))
	if err != nil {
		fmt.Printf("Cannot get driver information")
		return 72
	}

	// regs = calloc(1, sizeof(*regs)+drvinfo.regdump_len);
	// if (!regs) {
	// 	fmt.Printf("Cannot allocate memory for register dump");
	// 	return 73;
	// }
	regs := ethtool_regs{
		cmd: ETHTOOL_GREGS,
		len: drvinfo.regdump_len,
	}
	err = send_ioctl(ctx, uintptr(unsafe.Pointer(&regs)))
	if err != nil {
		fmt.Printf("Cannot get register dump: %v\n", err)
		return 74
	}

	// if (!gregs_dump_raw && gregs_dump_file != NULL) {
	// 	/* overwrite reg values from file dump */
	// 	FILE *f = fopen(gregs_dump_file, "r");
	// 	struct ethtool_regs *nregs;
	// 	struct stat st;
	// 	size_t nread;

	// 	if (!f || fstat(fileno(f), &st) < 0) {
	// 		fmt.Printf( "Can't open '%s': %s\n",
	// 			gregs_dump_file, strerror(errno));
	// 		if (f)
	// 			fclose(f);

	// 		return 75;
	// 	}

	// 	nregs = realloc(regs, sizeof(*regs) + st.st_size);
	// 	if (!nregs) {
	// 		fmt.Printf("Cannot allocate memory for register dump");
	// 		 /* was not freed by realloc */
	// 		return 73;
	// 	}
	// 	regs = nregs;
	// 	regs.len = st.st_size;
	// 	nread = fread(regs.data, regs.len, 1, f);
	// 	fclose(f);
	// 	if (nread != 1) {

	// 		return 75;
	// 	}
	// }

	if dump_regs(gregs_dump_raw, gregs_dump_hex,
		&drvinfo, &regs) < 0 {
		fmt.Printf("Cannot dump registers\n")

		return 75
	}
	return 0
}

func do_geeprom(ctx *cmd_context) int {

	// int geeprom_changed = 0;
	geeprom_dump_raw := 0
	geeprom_offset := uint32(0)
	geeprom_length := uint32(0)
	geeprom_length_seen := 0

	// parse_generic_cmdline(ctx, &geeprom_changed,
	// 		      cmdline_geeprom, ARRAY_SIZE(cmdline_geeprom));

	drvinfo := ethtool_drvinfo{cmd: ETHTOOL_GDRVINFO}
	err := send_ioctl(ctx, uintptr(unsafe.Pointer(&drvinfo)))
	if err != nil {
		fmt.Printf("Cannot get driver information: %v\n", err)
		return 74
	}

	if geeprom_length_seen == 0 {
		geeprom_length = drvinfo.eedump_len
	}

	if drvinfo.eedump_len < geeprom_offset+geeprom_length {
		geeprom_length = drvinfo.eedump_len - geeprom_offset
	}

	eeprom := ethtool_eeprom{
		cmd:    ETHTOOL_GEEPROM,
		len:    geeprom_length,
		offset: geeprom_offset,
	}
	err = send_ioctl(ctx, uintptr(unsafe.Pointer(&eeprom)))
	if err != nil {
		fmt.Printf("Cannot get EEPROM data: %v\n", err)
		return 74
	}
	ret := dump_eeprom(geeprom_dump_raw, &drvinfo, &eeprom)

	return ret

}

func do_phys_id(ctx *cmd_context) int {

	var edata ethtool_value
	var phys_id_time int

	// if (ctx.argc > 1)
	// 	exit_bad_args();
	// if (ctx.argc == 1)
	// 	phys_id_time = get_int(*ctx.argp, 0);
	// else
	phys_id_time = 0

	edata.cmd = ETHTOOL_PHYS_ID
	edata.data = uint32(phys_id_time)
	err := send_ioctl(ctx, uintptr(unsafe.Pointer(&edata)))
	if err != nil {
		fmt.Printf("Cannot identify NIC: %v", err)
		return -1
	}

	return 0
}

func do_gstats(ctx *cmd_context, cmd uint32, stringset uint32, name string) int {
	stats := &ethtool_stats{}

	// if (ctx.argc != 0)
	// 	exit_bad_args();
	drvinfo := ethtool_drvinfo{}

	strings := get_stringset(ctx, stringset,
		unsafe.Offsetof(drvinfo.n_stats), 0)
	if strings == nil {
		fmt.Printf("Cannot get stats strings information")
		return 96
	}

	n_stats := strings.len
	if n_stats < 1 {
		fmt.Printf("no stats available\n")
		return 94
	}

	stats.cmd = cmd
	stats.n_stats = n_stats
	err := send_ioctl(ctx, uintptr(unsafe.Pointer(stats)))
	if err != nil {
		fmt.Printf("Cannot get stats information")
		return 97
	}

	/* todo - pretty-print the strings per-driver */
	fmt.Printf("%s statistics:\n", name)
	for i := 0; uint32(i) < n_stats; i++ {
		fmt.Printf("     %.*s: %d\n",
			ETH_GSTRING_LEN,
			string(strings.data[i*ETH_GSTRING_LEN:(i+1)*ETH_GSTRING_LEN]),
			stats.data[i])
	}

	return 0
}

func do_gnicstats(ctx *cmd_context) int {
	return do_gstats(ctx, ETHTOOL_GSTATS, ETH_SS_STATS, "NIC")
}

func do_gphystats(ctx *cmd_context) int {
	return do_gstats(ctx, ETHTOOL_GPHYSTATS, ETH_SS_PHY_STATS, "PHY")
}

func rxflow_str_to_type(str string) int {
	flow_type := 0

	if str == "tcp4" {
		flow_type = TCP_V4_FLOW
	} else if str == "udp4" {
		flow_type = UDP_V4_FLOW
	} else if (str == "ah4") || (str == "esp4") {
		flow_type = AH_ESP_V4_FLOW
	} else if str == "sctp4" {
		flow_type = SCTP_V4_FLOW
	} else if str == "tcp6" {
		flow_type = TCP_V6_FLOW
	} else if str == "udp6" {
		flow_type = UDP_V6_FLOW
	} else if (str == "ah6") || (str == "esp6") {
		flow_type = AH_ESP_V6_FLOW
	} else if str == "sctp6" {
		flow_type = SCTP_V6_FLOW
	} else if str == "ether" {
		flow_type = ETHER_FLOW
	}
	return flow_type
}

func do_grxclass(ctx *cmd_context) int {
	var err error
	nfccmd := ethtool_rxnfc{}

	if ctx.argc > 0 && ctx.argp[0] == "rx-flow-hash" {
		var rx_fhash_get int
		flow_rss := false

		if ctx.argc == 4 {
			if ctx.argp[2] == "context" {
				// exit_bad_args();
			}
			flow_rss = true
			val, _ := strconv.ParseUint(ctx.argp[3], 10, 32)
			nfccmd.rule_cnt = uint32(val)
		} else if ctx.argc != 2 {
			// exit_bad_args();
		}

		rx_fhash_get = rxflow_str_to_type(ctx.argp[1])
		if rx_fhash_get == 0 {
			// exit_bad_args()
		}

		nfccmd.cmd = ETHTOOL_GRXFH
		nfccmd.flow_type = uint32(rx_fhash_get)

		if !flow_rss {
			nfccmd.flow_type |= FLOW_RSS
		}
		err := send_ioctl(ctx, uintptr(unsafe.Pointer(&nfccmd)))
		if err != nil {
			fmt.Print("Cannot get RX network flow hashing options\n")
		} else {
			if !flow_rss {
				fmt.Printf("For RSS context %d:\n", nfccmd.rule_cnt)
			}
			dump_rxfhash(rx_fhash_get, nfccmd.data)
		}
	} else if ctx.argc == 2 && (ctx.argp[0] == "rule") {
		rx_class_rule_get, _ := strconv.ParseUint(ctx.argp[1], 10, 32)

		err := rxclass_rule_get(ctx, uint32(rx_class_rule_get))
		if err != nil {
			fmt.Printf("Cannot get RX classification rule\n")
		}
	} else if ctx.argc == 0 {
		nfccmd.cmd = ETHTOOL_GRXRINGS
		err = send_ioctl(ctx, uintptr(unsafe.Pointer(&nfccmd)))
		if err != nil {
			fmt.Printf("Cannot get RX rings: %v\n", err)
		} else {
			fmt.Printf("%d RX rings available\n", int(nfccmd.data))
		}
		err = rxclass_rule_getall(ctx)
		if err != nil {
			fmt.Printf("RX classification rule retrieval failed\n")
		}

	} else {
		// exit_bad_args();
	}

	if err != nil {
		return 1
	}
	return 0
}

func do_tsinfo(ctx *cmd_context) int {

	// if (ctx.argc != 0)
	// 	exit_bad_args();

	fmt.Printf("Time stamping parameters for %s:\n", ctx.devname)
	info := ethtool_ts_info{
		cmd: ETHTOOL_GET_TS_INFO,
	}

	err := send_ioctl(ctx, uintptr(unsafe.Pointer(&info)))
	if err != nil {
		fmt.Printf("Cannot get device time stamping settings: %v\n", err)
		return -1
	}
	dump_tsinfo(&info)
	return 0
}

func print_indir_table(ctx *cmd_context,
	ring_count *ethtool_rxnfc,
	indir_size uint32, indir *uint32) {

	fmt.Printf("RX flow hash indirection table for %s with %d RX ring(s):\n",
		ctx.devname, ring_count.data)

	if indir_size == 0 {
		fmt.Printf("Operation not supported\n")
	}
	slic_indir := *(*[MAX_DATA_BUF]uint32)(unsafe.Pointer(indir))
	for i := uint32(0); i < indir_size; i++ {
		if i%8 == 0 {
			fmt.Printf("%5d: ", i)
		}
		fmt.Printf(" %5d", slic_indir[i])
		if i%8 == 7 || i == indir_size-1 {
			fmt.Printf("\n")
		}
	}
}

func do_grxfhindir(ctx *cmd_context,
	ring_count *ethtool_rxnfc) int {

	indir_head := ethtool_rxfh_indir{
		cmd:  ETHTOOL_GRXFHINDIR,
		size: 0,
	}
	err := send_ioctl(ctx, uintptr(unsafe.Pointer(&indir_head)))
	if err != nil {
		fmt.Printf("Cannot get RX flow hash indirection table size")
		return 1
	}

	indir := ethtool_rxfh_indir{
		cmd:  ETHTOOL_GRXFHINDIR,
		size: indir_head.size,
	}

	err = send_ioctl(ctx, uintptr(unsafe.Pointer(&indir)))
	if err != nil {
		fmt.Printf("Cannot get RX flow hash indirection table")
		return 1
	}
	print_indir_table(ctx, ring_count, indir.size, &indir.ring_index[0])

	return 0
}
func do_grxfh(ctx *cmd_context) int {

	rss_context := uint32(0)
	// unsigned int arg_num = 0;
	// char *hkey;

	// while (arg_num < ctx.argc) {
	// 	if (!strcmp(ctx.argp[arg_num], "context")) {
	// 		++arg_num;
	// 		rss_context = get_int_range(ctx.argp[arg_num], 0, 1,
	// 					    ETH_RXFH_CONTEXT_ALLOC - 1);
	// 		++arg_num;
	// 	} else {
	// 		exit_bad_args();
	// 	}
	// }

	ring_count := ethtool_rxnfc{
		cmd: ETHTOOL_GRXRINGS,
	}
	err := send_ioctl(ctx, uintptr(unsafe.Pointer(&ring_count)))
	if err != nil {
		fmt.Printf("Cannot get RX ring count")
		return 1
	}

	rss_head := ethtool_rxfh{
		cmd:         ETHTOOL_GRSSH,
		rss_context: rss_context,
	}
	err = send_ioctl(ctx, uintptr(unsafe.Pointer(&rss_head)))
	if err != nil && err == syscall.EOPNOTSUPP && rss_context == 0 {
		return do_grxfhindir(ctx, &ring_count)
	} else if err != nil {
		fmt.Printf("Cannot get RX flow hash indir size and/or key size")
		return 1
	}

	rss := ethtool_rxfh{
		cmd:         ETHTOOL_GRSSH,
		rss_context: rss_context,
		indir_size:  rss_head.indir_size,
		key_size:    rss_head.key_size,
	}

	// rssv := make([]uint32, rss_head.indir_size+rss_head.key_size/4)
	// rss.rss_config = (*uint32)(unsafe.Pointer(&rssv[0]))
	err = send_ioctl(ctx, uintptr(unsafe.Pointer(&rss)))
	if err != nil {
		fmt.Printf("Cannot get RX flow hash configuration")
		return 1
	}

	print_indir_table(ctx, &ring_count, rss.indir_size, &rss.rss_config[0])

	// cdata := *(*[]uint32)(unsafe.Pointer(rss.rss_config))

	indir_bytes := unsafe.Pointer(&rss.rss_config[rss.indir_size])
	hkey := *(*[]byte)(indir_bytes)

	fmt.Printf("RSS hash key:\n")
	if rss.key_size == 0 {
		fmt.Printf("Operation not supported\n")
	}
	for i := uint32(0); i < rss.key_size; i++ {
		if i == (rss.key_size - 1) {
			fmt.Printf("%02x\n", hkey[i])
		} else {
			fmt.Printf("%02x:", hkey[i])
		}
	}

	fmt.Printf("RSS hash function:\n")
	if rss.hfunc == 0 {
		fmt.Printf("    Operation not supported\n")
		return 0
	}

	hfuncs := get_stringset(ctx, ETH_SS_RSS_HASH_FUNCS, 0, 1)
	if hfuncs == nil {
		fmt.Printf("Cannot get hash functions names")
		return 1
	}

	// for i := 0; i < hfuncs.len; i++{
	// 	fmt.Printf("    %s: %s\n",
	// 	       (const char *)hfuncs.data + i * ETH_GSTRING_LEN,
	// 	       (rss.hfunc & (1 << i)) ? "on" : "off")
	// }

	return 0
}

func do_permaddr(ctx *cmd_context) int {
	epaddr := ethtool_perm_addr{
		cmd:  ETHTOOL_GPERMADDR,
		size: 32,
	}
	err := send_ioctl(ctx, uintptr(unsafe.Pointer(&epaddr)))
	if err != nil {
		fmt.Printf("Cannot read permanent address: %v\n", err)
	} else {
		fmt.Printf("Permanent address:")
		for i := uint32(0); i < epaddr.size; i++ {
			if i == 0 {
				fmt.Printf("%c%02x", ' ', epaddr.data[i])
			}
			fmt.Printf("%c%02x", ':', epaddr.data[i])
		}
		fmt.Printf("\n")
	}
	return 0
}

func do_writefwdump(dump *ethtool_dump, dump_file string) error {

	f, err := os.OpenFile(dump_file, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Printf("Can't open file %s: %v\n", dump_file, err)
		return errors.New("Can't open file")
	}
	defer f.Close()
	_, err = f.Write(dump.data[:])
	if err != nil {
		fmt.Printf("Can not write all of dump data\n")
	}
	return err
}

func do_getfwdump(ctx *cmd_context) int {
	dump_flag := ETHTOOL_GET_DUMP_DATA
	dump_file := ""

	if ctx.argc == 2 && ctx.argp[0] == "data" {
		dump_flag = ETHTOOL_GET_DUMP_DATA
		dump_file = string(ctx.argp[1])
	} else if ctx.argc == 0 {
		dump_flag = 0
		dump_file = ""
	}
	// else {
	// 	exit_bad_args();
	// }

	edata := ethtool_dump{cmd: ETHTOOL_GET_DUMP_FLAG}

	err := send_ioctl(ctx, uintptr(unsafe.Pointer(&edata)))
	if err != nil {
		fmt.Printf("Can not get dump level: %v\n", err)
		return 1
	}
	if dump_flag != ETHTOOL_GET_DUMP_DATA {
		fmt.Printf("flag: %d, version: %d, length: %d\n",
			edata.flag, edata.version, edata.len)
		return 0
	}
	// data = calloc(1, offsetof(struct ethtool_dump, data) + edata.len);
	// if (!data) {
	// 	perror("Can not allocate enough memory\n");
	// 	return 1;
	// }
	data := ethtool_dump{
		cmd: ETHTOOL_GET_DUMP_DATA,
		len: edata.len,
	}
	err = send_ioctl(ctx, uintptr(unsafe.Pointer(&data)))
	if err != nil {
		fmt.Printf("Can not get dump data: %v\n", err)
		return 1
	}
	err = do_writefwdump(&data, dump_file)
	if err != nil {
		return 1
	}
	return 0
}

func do_gchannels(ctx *cmd_context) int {

	// if (ctx.argc != 0)
	// 	exit_bad_args();

	fmt.Printf("Channel parameters for %s:\n", ctx.devname)

	echannels := ethtool_channels{cmd: ETHTOOL_GCHANNELS}
	err := send_ioctl(ctx, uintptr(unsafe.Pointer(&echannels)))
	if err == nil {
		dump_channels(&echannels)
	} else {
		fmt.Printf("Cannot get device channel parameters %v\n", err)
		return 1
	}
	return 0

}

func do_gprivflags(ctx *cmd_context) int {

	max_len, cur_len := 0, 0

	// if (ctx.argc != 0)
	// 	exit_bad_args()
	var drvinfo ethtool_drvinfo
	strings := get_stringset(ctx, ETH_SS_PRIV_FLAGS,
		unsafe.Offsetof(drvinfo.n_priv_flags), 1)
	if strings == nil {
		fmt.Printf("Cannot get private flag names\n")
		return 1
	}
	if strings.len == 0 {
		fmt.Printf("No private flags defined\n")
		return 1
	}
	if strings.len > 32 {
		/* ETHTOOL_GPFLAGS can only cover 32 flags */
		fmt.Printf("Only showing first 32 private flags\n")
		strings.len = 32
	}

	flags := ethtool_value{
		cmd: ETHTOOL_GPFLAGS,
	}
	err := send_ioctl(ctx, uintptr(unsafe.Pointer(&flags)))
	if err != nil {
		fmt.Printf("Cannot get private flags: %v", err)
		return -1
	}

	/* Find longest string and align all strings accordingly */
	for i := uint32(0); i < strings.len; i++ {
		cur_len = len(string(strings.data[i*ETH_GSTRING_LEN]))
		if cur_len > max_len {
			max_len = cur_len
		}
	}

	fmt.Printf("Private flags for %s:\n", ctx.devname)
	flag_str := "off"

	for i := uint32(0); i < strings.len; i++ {
		if (flags.data & (1 << i)) != 0 {
			flag_str = "on"
		}
		fmt.Printf("%-*s: %s\n",
			max_len,
			string(strings.data[i*ETH_GSTRING_LEN]),
			flag_str)
	}

	return 0
}

func do_getmodule(ctx *cmd_context) int {

	geeprom_offset := uint32(0)
	geeprom_length := uint32(0)
	// geeprom_changed := 0
	geeprom_dump_raw := 0
	geeprom_dump_hex := 0
	geeprom_length_seen := 0

	// parse_generic_cmdline(ctx, &geeprom_changed,
	// 		      cmdline_geeprom, ARRAY_SIZE(cmdline_geeprom));

	if geeprom_dump_raw != 0 && geeprom_dump_hex != 0 {
		fmt.Printf("Hex and raw dump cannot be specified together\n")
		return 1
	}

	modinfo := ethtool_modinfo{cmd: ETHTOOL_GMODULEINFO}
	err := send_ioctl(ctx, uintptr(unsafe.Pointer(&modinfo)))
	if err != nil {
		fmt.Printf("Cannot get module EEPROM information")
		return 1
	}

	if geeprom_length_seen == 0 {
		geeprom_length = modinfo.eeprom_len
	}

	if modinfo.eeprom_len < uint32(geeprom_offset+geeprom_length) {
		geeprom_length = modinfo.eeprom_len - geeprom_offset
	}

	// eeprom = calloc(1, sizeof(*eeprom)+geeprom_length);
	// if (!eeprom) {
	// 	perror("Cannot allocate memory for Module EEPROM data");
	// 	return 1;
	// }

	eeprom := ethtool_eeprom{
		cmd:    ETHTOOL_GMODULEEEPROM,
		len:    geeprom_length,
		offset: geeprom_offset,
	}
	err = send_ioctl(ctx, uintptr(unsafe.Pointer(&eeprom)))
	if err != nil {
		fmt.Printf("Cannot get Module EEPROM data: %v\n", err)
		if err == syscall.ENODEV || err == syscall.EIO ||
			err == syscall.ENXIO {
			fmt.Printf("SFP module not in cage?\n")
		}
		return 1
	}

	/*
	 * SFF-8079 EEPROM layout contains the memory available at A0 address on
	 * the PHY EEPROM.
	 * SFF-8472 defines a virtual extension of the EEPROM, where the
	 * microcontroller on the SFP/SFP+ generates a page at the A2 address,
	 * which contains data relative to optical diagnostics.
	 * The current kernel implementation returns a blob, which contains:
	 *  - ETH_MODULE_SFF_8079 => The A0 page only.
	 *  - ETH_MODULE_SFF_8472 => The A0 and A2 page concatenated.
	 */
	if geeprom_dump_raw != 0 {
		// fwrite(eeprom.data, 1, eeprom.len)
	} else {
		if eeprom.offset != 0 ||
			(eeprom.len != modinfo.eeprom_len) {
			geeprom_dump_hex = 1
		} else if geeprom_dump_hex == 0 {
			switch modinfo.tp {

			case ETH_MODULE_SFF_8079:
				//sff8079_show_all(eeprom.data)
				break
			case ETH_MODULE_SFF_8472:
				// sff8079_show_all(eeprom.data)
				// sff8472_show_all(eeprom.data)
				break
			case ETH_MODULE_SFF_8436:
			case ETH_MODULE_SFF_8636:
				// sff8636_show_all(eeprom.data,modinfo.eeprom_len)
				break

			default:
				geeprom_dump_hex = 1
				break
			}
		}
		if geeprom_dump_hex != 0 {
			dump_hex(os.Stdout, eeprom.data[:],
				eeprom.len, eeprom.offset)
		}
	}

	return 0
}

func do_geee(ctx *cmd_context) int {

	// if (ctx.argc != 0)
	// 	exit_bad_args();

	eeecmd := ethtool_eee{cmd: ETHTOOL_GEEE}
	err := send_ioctl(ctx, uintptr(unsafe.Pointer(&eeecmd)))
	if err != nil {
		fmt.Printf("Cannot get EEE settings: %v\n", err)
		return 1
	}

	fmt.Printf("EEE Settings for %s:\n", ctx.devname)
	dump_eeecmd(&eeecmd)

	return 0
}

type func_t func(*cmd_context) int
type options struct {
	opts      *bool /*flag string*/
	no_dev    bool
	ioctlfunc func_t
	nlfunc    func_t
	xhelp     string
}

var (
	opt_args = []options{
		{flag.Bool("s", false, "Change generic options"), false, nil, nil, "		[ speed %d ]\n" +
			"		[ duplex half|full ]\n" +
			"		[ port tp|aui|bnc|mii|fibre|da ]\n" +
			"		[ mdix auto|on|off ]\n" +
			"		[ autoneg on|off ]\n" +
			"		[ advertise %x[/%x] | mode on|off ... [--] ]\n" +
			"		[ phyad %d ]\n" +
			"		[ xcvr internal|external ]\n" +
			"		[ wol %d[/%d] | p|u|m|b|a|g|s|f|d... ]\n" +
			"		[ sopass %x:%x:%x:%x:%x:%x ]\n" +
			"		[ msglvl %d[/%d] | type on|off ... [--] ]\n" +
			"		[ master-slave master-preferred|slave-preferred|master-force|slave-force ]\n"},
		{flag.Bool("a", false, "Show pause options"), false, do_gpause, nil, ""},
		{flag.Bool("A", false, "Set pause options"), false, nil, nil, "		[ autoneg on|off ]\n" +
			"		[ rx on|off ]\n" +
			"		[ tx on|off ]\n"},
		{flag.Bool("c", false, "Show coalesce options"), false, do_gcoalesce, nil, ""},
		{flag.Bool("C", false, "Show pause options"), false, nil, nil, "		[adaptive-rx on|off]\n" +
			"		[adaptive-tx on|off]\n" +
			"		[rx-usecs N]\n" +
			"		[rx-frames N]\n" +
			"		[rx-usecs-irq N]\n" +
			"		[rx-frames-irq N]\n" +
			"		[tx-usecs N]\n" +
			"		[tx-frames N]\n" +
			"		[tx-usecs-irq N]\n" +
			"		[tx-frames-irq N]\n" +
			"		[stats-block-usecs N]\n" +
			"		[pkt-rate-low N]\n" +
			"		[rx-usecs-low N]\n" +
			"		[rx-frames-low N]\n" +
			"		[tx-usecs-low N]\n" +
			"		[tx-frames-low N]\n" +
			"		[pkt-rate-high N]\n" +
			"		[rx-usecs-high N]\n" +
			"		[rx-frames-high N]\n" +
			"		[tx-usecs-high N]\n" +
			"		[tx-frames-high N]\n" +
			"		[sample-interval N]\n"},
		{flag.Bool("g", false, "Query RX/TX ring parameters"), false, do_gring, nil, ""},

		{flag.Bool("G", false, "Set RX/TX ring parameters"), false, nil, nil, "		[ rx N ]\n" +
			"		[ rx-mini N ]\n" +
			"		[ rx-jumbo N ]\n" +
			"		[ tx N ]\n"},
		{flag.Bool("k", false, "Get state of protocol offload and other features"), false, do_gfeatures, nil, ""},
		{flag.Bool("K", false, "Set protocol offload and other features"), false, nil, nil, "		FEATURE on|off ...\n"},

		{flag.Bool("i", false, "Show driver information"), false, do_gdrv, nil, ""},
		{flag.Bool("d", false, "Do a register dump"), false, do_gregs, nil, "		[ raw on|off ]\n" +
			"		[ file FILENAME ]\n"},
		{flag.Bool("e", false, "Do a EEPROM dump"), false, do_geeprom, nil, "		[ raw on|off ]\n" +
			"		[ offset N ]\n" +
			"		[ length N ]\n"},
		{flag.Bool("E", false, "Change bytes in device EEPROM"), false, nil, nil, "		[ magic N ]\n" +
			"		[ offset N ]\n" +
			"		[ length N ]\n" +
			"		[ value N ]\n"},
		{flag.Bool("r", false, "Restart N-WAY negotiation"), false, nil, nil, ""},
		{flag.Bool("p", false, "Show visible port identification (e.g. blinking)"), false, do_phys_id, nil, "               [ TIME-IN-SECONDS ]\n"},
		{flag.Bool("t", false, "Execute adapter self test"), false, nil, nil, "               [ online | offline | external_lb ]\n"},
		{flag.Bool("S", false, "Show adapter statistics"), false, do_gnicstats, nil, ""},
		{flag.Bool("phy-statistics", false, "Show phy statistics"), false, do_gphystats, nil, ""},
		{flag.Bool("n", false, "Show Rx network flow classification options or rules"), false, do_grxclass, nil, "		[ rx-flow-hash tcp4|udp4|ah4|esp4|sctp4|" +
			"tcp6|udp6|ah6|esp6|sctp6 [context %d] |\n" +
			"		  rule %d ]\n"},
		{flag.Bool("N", false, "Configure Rx network flow classification options or rules"), false, nil, nil, "		rx-flow-hash tcp4|udp4|ah4|esp4|sctp4|" +
			"tcp6|udp6|ah6|esp6|sctp6 m|v|t|s|d|f|n|r... [context %d] |\n" +
			"		flow-type ether|ip4|tcp4|udp4|sctp4|ah4|esp4|" +
			"ip6|tcp6|udp6|ah6|esp6|sctp6\n" +
			"			[ src %x:%x:%x:%x:%x:%x [m %x:%x:%x:%x:%x:%x] ]\n" +
			"			[ dst %x:%x:%x:%x:%x:%x [m %x:%x:%x:%x:%x:%x] ]\n" +
			"			[ proto %d [m %x] ]\n" +
			"			[ src-ip IP-ADDRESS [m IP-ADDRESS] ]\n" +
			"			[ dst-ip IP-ADDRESS [m IP-ADDRESS] ]\n" +
			"			[ tos %d [m %x] ]\n" +
			"			[ tclass %d [m %x] ]\n" +
			"			[ l4proto %d [m %x] ]\n" +
			"			[ src-port %d [m %x] ]\n" +
			"			[ dst-port %d [m %x] ]\n" +
			"			[ spi %d [m %x] ]\n" +
			"			[ vlan-etype %x [m %x] ]\n" +
			"			[ vlan %x [m %x] ]\n" +
			"			[ user-def %x [m %x] ]\n" +
			"			[ dst-mac %x:%x:%x:%x:%x:%x [m %x:%x:%x:%x:%x:%x] ]\n" +
			"			[ action %d ] | [ vf %d queue %d ]\n" +
			"			[ context %d ]\n" +
			"			[ loc %d]] |\n" +
			"		delete %d\n"},
		{flag.Bool("T", false, "Show time stamping capabilities"), false, do_tsinfo, nil, ""},
		{flag.Bool("x", false, "Show Rx flow hash indirection table and/or RSS hash key"), false, do_grxfh, nil, "		[ context %d ]\n"},
		{flag.Bool("-X|--set-rxfh-indir|--rxfh", false, "Set Rx flow hash indirection table and/or RSS hash key"), false, nil, nil, "		[ context %d|new ]\n" +
			"		[ equal N | weight W0 W1 ... | default ]\n" +
			"		[ hkey %x:%x:%x:%x:%x:.... ]\n" +
			"		[ hfunc FUNC ]\n" +
			"		[ delete ]\n"},
		{flag.Bool("f", false, "Flash firmware image from the specified file to a region on the device"), false, nil, nil, "               FILENAME [ REGION-NUMBER-TO-FLASH ]\n"},
		{flag.Bool("P", false, "Show permanent hardware address"), false, do_permaddr, nil, ""},
		{flag.Bool("w", false, "Get dump flag, data"), false, do_getfwdump, nil, "		[ data FILENAME ]\n"},
		{flag.Bool("W", false, "Set dump flag of the device"), false, nil, nil, "		N\n"},
		{flag.Bool("l", false, "Query Channels"), false, do_gchannels, nil, ""},
		{flag.Bool("L", false, "Set Channels"), false, nil, nil, "               [ rx N ]\n" +
			"               [ tx N ]\n" +
			"               [ other N ]\n" +
			"               [ combined N ]\n"},
		{flag.Bool("show-priv-flags", false, "Query private flags"), false, do_gprivflags, nil, ""},
		{flag.Bool("set-priv-flag", false, "Set private flags"), false, nil, nil, "		FLAG on|off ...\n"},
		{flag.Bool("m", false, "Query/Decode Module EEPROM information and optical diagnostics if available"), false, do_getmodule, nil, "		[ raw on|off ]\n" +
			"		[ hex on|off ]\n" +
			"		[ offset N ]\n" +
			"		[ length N ]\n"},
		{flag.Bool("show-eee", false, "Show EEE settings"), false, do_geee, nil, ""},
		{flag.Bool("set-eee", false, "Set EEE settings"), false, nil, nil, "		[ eee on|off ]\n" +
			"		[ advertise %x ]\n" +
			"		[ tx-lpi on|off ]\n" +
			"		[ tx-timer %d ]\n"},
		{flag.Bool("set-phy-tunable", false, "Set PHY tunable"), false, nil, nil, "		[ downshift on|off [count N] ]\n" +
			"		[ fast-link-down on|off [msecs N] ]\n" +
			"		[ energy-detect-power-down on|off [msecs N] ]\n"},
		{flag.Bool("get-phy-tunable", false, "Get PHY tunable"), false, nil, nil, "		[ downshift ]\n" +
			"		[ fast-link-down ]\n" +
			"		[ energy-detect-power-down ]\n"},
		{flag.Bool("set-tunable", false, "Set tunable"), false, nil, nil, "		[ rx-copybreak N]\n" +
			"		[ tx-copybreak N]\n" +
			"		[ pfc-precention-tout N]\n"},
		{flag.Bool("get-tunable", false, "Get tunable"), false, nil, nil, "		[ rx-copybreak ]\n" +
			"		[ tx-copybreak ]\n" +
			"		[ pfc-precention-tout ]\n"},
		{flag.Bool("reset", false, "Reset components"), false, nil, nil, "		[ flags %x ]\n" +
			"		[ mgmt ]\n" +
			"		[ mgmt-shared ]\n" +
			"		[ irq ]\n" +
			"		[ irq-shared ]\n" +
			"		[ dma ]\n" +
			"		[ dma-shared ]\n" +
			"		[ filter ]\n" +
			"		[ filter-shared ]\n" +
			"		[ offload ]\n" +
			"		[ offload-shared ]\n" +
			"		[ mac ]\n" +
			"		[ mac-shared ]\n" +
			"		[ phy ]\n" +
			"		[ phy-shared ]\n" +
			"		[ ram ]\n" +
			"		[ ram-shared ]\n" +
			"		[ ap ]\n" +
			"		[ ap-shared ]\n" +
			"		[ dedicated ]\n" +
			"		[ all ]\n"},
		{flag.Bool("-show-fec", false, "Show FEC setting"), false, nil, nil, ""},
		{flag.Bool("-set-fec", false, "Set FEC setting"), false, nil, nil, "		[ encoding auto|off|rs|baser|llrs [...]]\n"},
		{flag.Bool("Q", false, "Apply per-queue command."), false, nil, nil, "The supported sub commands include --show-coalesce, --coalesce" +
			"             [queue_mask %x] SUB_COMMAND\n"},
		{flag.Bool("-cable-test", false, "Perform a cable test"), false, nil, nil, ""},
		{flag.Bool("-cable-test-tdr", false, "Print cable test time domain reflectrometery data"), false, nil, nil, "		[ first N ]\n" +
			"		[ last N ]\n" +
			"		[ step N ]\n" +
			"		[ pair N ]\n"},
		{flag.Bool("-show-tunnels", false, "Show NIC tunnel offload information"), false, nil, nil, ""},
		{flag.Bool("-version", false, "Show version number"), false, nil, nil, ""},
	}
)

// Show_usage
func Show_usage() {
	fmt.Printf("ethtool version %s\n", "1.0.0")
	fmt.Printf("Usage:\n" +
		"        ethtool [ FLAGS ] DEVNAME\t" +
		"Display standard information about device\n")
	flag.PrintDefaults()
	fmt.Printf("\n")
	fmt.Printf("FLAGS:\n")
	fmt.Printf("	--debug MASK	turn on debugging messages\n")
	fmt.Printf("	--json		enable JSON output format (not supported by all commands)\n")
	fmt.Printf("	-I|--include-statistics		request device statistics related to the command (not supported by all commands)\n")

}

// Parse_args parse the args
func Parse_args() error {
	flag.Parse()
	cnt := 0
	for i := 0; i < len(opt_args); i++ {
		if *opt_args[i].opts {
			cnt++
		}
	}
	if cnt > 1 {
		fmt.Printf("ethtool: bad command line argument(s)\n")
		return errors.New("bad command line argument(s)")
	}
	return nil
}

// Do_actions will call ioctl to get or set infos
func Do_actions() int {
	var ctx cmd_context
	no_dev := false
	i := 0
	for ; i < len(opt_args); i++ {
		if *opt_args[i].opts {
			break
		}
	}
	if i >= len(opt_args) {
		return -1
	}
	args := flag.Args()
	ctx.devname = args[0]
	err := init_ioctl(&ctx, no_dev)
	if err != 0 {
		return err
	}
	return opt_args[i].ioctlfunc(&ctx)
}

func find_max_num_queues(ctx *cmd_context) int {
	var echannels ethtool_channels
	echannels.cmd = ETHTOOL_GCHANNELS
	err := send_ioctl(ctx, uintptr(unsafe.Pointer(&echannels)))
	if err != nil {
		return -1
	}
	return int(math.Max(float64(echannels.rx_count), float64(echannels.tx_count))) + int(echannels.combined_count)
}
