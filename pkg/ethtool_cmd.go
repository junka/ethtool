package ethtool

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
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
	dump_hex(os.Stdout, ee.data, ee.len, ee.offset)

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
			"RX:		%u\n"+
			"RX Mini:	%u\n"+
			"RX Jumbo:	%u\n"+
			"TX:		%u\n",
		ering.rx_max_pending,
		ering.rx_mini_max_pending,
		ering.rx_jumbo_max_pending,
		ering.tx_max_pending)

	fmt.Printf(
		"Current hardware settings:\n"+
			"RX:		%u\n"+
			"RX Mini:	%u\n"+
			"RX Jumbo:	%u\n"+
			"TX:		%u\n",
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
			"RX:		%u\n"+
			"TX:		%u\n"+
			"Other:		%u\n"+
			"Combined:	%u\n",
		echannels.max_rx, echannels.max_tx,
		echannels.max_other,
		echannels.max_combined)

	fmt.Printf(
		"Current hardware settings:\n"+
			"RX:		%u\n"+
			"TX:		%u\n"+
			"Other:		%u\n"+
			"Combined:	%u\n",
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
		"stats-block-usecs: %u\n"+
			"sample-interval: %u\n"+
			"pkt-rate-low: %u\n"+
			"pkt-rate-high: %u\n"+
			"\n"+
			"rx-usecs: %u\n"+
			"rx-frames: %u\n"+
			"rx-usecs-irq: %u\n"+
			"rx-frames-irq: %u\n"+
			"\n"+
			"tx-usecs: %u\n"+
			"tx-frames: %u\n"+
			"tx-usecs-irq: %u\n"+
			"tx-frames-irq: %u\n"+
			"\n"+
			"rx-usecs-low: %u\n"+
			"rx-frames-low: %u\n"+
			"tx-usecs-low: %u\n"+
			"tx-frames-low: %u\n"+
			"\n"+
			"rx-usecs-high: %u\n"+
			"rx-frames-high: %u\n"+
			"tx-usecs-high: %u\n"+
			"tx-frames-high: %u\n"+
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
	// 			fprintf(stdout, "Queue: %d\n", queue);
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
func get_features(ctx *cmd_context, defs *feature_defs) *feature_state {

	var eval ethtool_value
	allfail := 1

	state := feature_state{}
	// malloc(sizeof(*state) +
	// 	       FEATURE_BITS_TO_BLOCKS(defs.n_features) *
	// 	       sizeof(state.features.features[0]));

	state.off_flags = 0

	for i := 0; i < OFF_FLAG_DEF_SIZE; i++ {
		value = off_flag_def[i].value
		if off_flag_def[i].get_cmd == 0 {
			continue
		}
		eval.cmd = off_flag_def[i].get_cmd
		err = send_ioctl(ctx, &eval)
		if err {
			if errno == EOPNOTSUPP &&
				off_flag_def[i].get_cmd == ETHTOOL_GUFO {
				continue
			}

			fprintf(stderr,
				"Cannot get device %s settings: %m\n",
				off_flag_def[i].long_name)
		} else {
			if eval.data != 0 {
				state.off_flags |= value
			}
			allfail = 0
		}
	}

	eval.cmd = ETHTOOL_GFLAGS
	err = send_ioctl(ctx, &eval)
	if err {
		perror("Cannot get device flags")
	} else {
		state.off_flags |= eval.data & ETH_FLAG_EXT_MASK
		allfail = 0
	}

	if defs.n_features > 0 {
		state.features.cmd = ETHTOOL_GFEATURES
		state.features.size = FEATURE_BITS_TO_BLOCKS(defs.n_features)
		err = send_ioctl(ctx, &state.features)
		if err != nil {
			fmt.Printf("Cannot get device generic features: %v\n", err)
		} else {
			allfail = 0
		}
	}

	if allfail {
		free(state)
		return NULL
	}

	return state
}

func do_gfeatures(ctx *cmd_context) int {
	feature_defs * defs
	feature_state * features

	// if (ctx.argc != 0)
	// 	exit_bad_args();

	defs = get_feature_defs(ctx)
	if defs == nil {
		fmt.Printf("Cannot get device feature names\n")
		return 1
	}

	fmt.Printf("Features for %s:\n", ctx.devname)

	features = get_features(ctx, defs)
	if !features {
		fmt.Printf("no feature info available\n")
		free(defs)
		return 1
	}

	dump_features(defs, features, NULL)
	free(features)
	free(defs)
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
		{flag.Bool("c", false, "Show coalesce options"), false, nil, nil, ""},
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
		{flag.Bool("g", false, "Query RX/TX ring parameters"), false, nil, nil, ""},

		{flag.Bool("G", false, "Set RX/TX ring parameters"), false, nil, nil, "		[ rx N ]\n" +
			"		[ rx-mini N ]\n" +
			"		[ rx-jumbo N ]\n" +
			"		[ tx N ]\n"},
		{flag.Bool("k", false, "Get state of protocol offload and other features"), false, nil, nil, ""},
		{flag.Bool("K", false, "Set protocol offload and other features"), false, nil, nil, "		FEATURE on|off ...\n"},

		{flag.Bool("i", false, "Show driver information"), false, do_gdrv, nil, ""},
		{flag.Bool("d", false, "Do a register dump"), false, nil, nil, "		[ raw on|off ]\n" +
			"		[ file FILENAME ]\n"},
		{flag.Bool("e", false, "Do a EEPROM dump"), false, nil, nil, "		[ raw on|off ]\n" +
			"		[ offset N ]\n" +
			"		[ length N ]\n"},
		{flag.Bool("E", false, "Change bytes in device EEPROM"), false, nil, nil, "		[ magic N ]\n" +
			"		[ offset N ]\n" +
			"		[ length N ]\n" +
			"		[ value N ]\n"},
		{flag.Bool("r", false, "Restart N-WAY negotiation"), false, nil, nil, ""},
		{flag.Bool("p", false, "Show visible port identification (e.g. blinking)"), false, nil, nil, "               [ TIME-IN-SECONDS ]\n"},
		{flag.Bool("t", false, "Execute adapter self test"), false, nil, nil, "               [ online | offline | external_lb ]\n"},
		{flag.Bool("S", false, "Show adapter statistics"), false, nil, nil, ""},
		{flag.Bool("-phy-statistics", false, "Show phy statistics"), false, nil, nil, ""},
		{flag.Bool("-n|-u|--show-nfc|--show-ntuple", false, "Show Rx network flow classification options or rules"), false, nil, nil, "		[ rx-flow-hash tcp4|udp4|ah4|esp4|sctp4|" +
			"tcp6|udp6|ah6|esp6|sctp6 [context %d] |\n" +
			"		  rule %d ]\n"},
		{flag.Bool("-N|-U|--config-nfc|--config-ntuple", false, "Configure Rx network flow classification options or rules"), false, nil, nil, "		rx-flow-hash tcp4|udp4|ah4|esp4|sctp4|" +
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
		{flag.Bool("-T|--show-time-stamping", false, "Show time stamping capabilities"), false, nil, nil, ""},
		{flag.Bool("-x|--show-rxfh-indir|--show-rxfh", false, "Show Rx flow hash indirection table and/or RSS hash key"), false, nil, nil, "		[ context %d ]\n"},
		{flag.Bool("-X|--set-rxfh-indir|--rxfh", false, "Set Rx flow hash indirection table and/or RSS hash key"), false, nil, nil, "		[ context %d|new ]\n" +
			"		[ equal N | weight W0 W1 ... | default ]\n" +
			"		[ hkey %x:%x:%x:%x:%x:.... ]\n" +
			"		[ hfunc FUNC ]\n" +
			"		[ delete ]\n"},
		{flag.Bool("f", false, "Flash firmware image from the specified file to a region on the device"), false, nil, nil, "               FILENAME [ REGION-NUMBER-TO-FLASH ]\n"},
		{flag.Bool("P", false, "Show permanent hardware address"), false, nil, nil, ""},
		{flag.Bool("w", false, "Get dump flag, data"), false, nil, nil, "		[ data FILENAME ]\n"},
		{flag.Bool("W", false, "Set dump flag of the device"), false, nil, nil, "		N\n"},
		{flag.Bool("l", false, "Query Channels"), false, nil, nil, ""},
		{flag.Bool("L", false, "Set Channels"), false, nil, nil, "               [ rx N ]\n" +
			"               [ tx N ]\n" +
			"               [ other N ]\n" +
			"               [ combined N ]\n"},
		{flag.Bool("-show-priv-flags", false, "Query private flags"), false, nil, nil, ""},
		{flag.Bool("-set-priv-flag", false, "Set private flags"), false, nil, nil, "		FLAG on|off ...\n"},
		{flag.Bool("m", false, "Query/Decode Module EEPROM information and optical diagnostics if available"), false, nil, nil, "		[ raw on|off ]\n" +
			"		[ hex on|off ]\n" +
			"		[ offset N ]\n" +
			"		[ length N ]\n"},
		{flag.Bool("-show-eee", false, "Show EEE settings"), false, nil, nil, ""},
		{flag.Bool("-set-eee", false, "Set EEE settings"), false, nil, nil, "		[ eee on|off ]\n" +
			"		[ advertise %x ]\n" +
			"		[ tx-lpi on|off ]\n" +
			"		[ tx-timer %d ]\n"},
		{flag.Bool("-set-phy-tunable", false, "Set PHY tunable"), false, nil, nil, "		[ downshift on|off [count N] ]\n" +
			"		[ fast-link-down on|off [msecs N] ]\n" +
			"		[ energy-detect-power-down on|off [msecs N] ]\n"},
		{flag.Bool("-get-phy-tunable", false, "Get PHY tunable"), false, nil, nil, "		[ downshift ]\n" +
			"		[ fast-link-down ]\n" +
			"		[ energy-detect-power-down ]\n"},
		{flag.Bool("-set-tunable", false, "Set tunable"), false, nil, nil, "		[ rx-copybreak N]\n" +
			"		[ tx-copybreak N]\n" +
			"		[ pfc-precention-tout N]\n"},
		{flag.Bool("--get-tunable", false, "Get tunable"), false, nil, nil, "		[ rx-copybreak ]\n" +
			"		[ tx-copybreak ]\n" +
			"		[ pfc-precention-tout ]\n"},
		{flag.Bool("-reset", false, "Reset components"), false, nil, nil, "		[ flags %x ]\n" +
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
