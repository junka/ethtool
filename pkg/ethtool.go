package ethtool

const (
	ETH_MDIO_SUPPORTS_C22 = 1
	ETH_MDIO_SUPPORTS_C45 = 2

	ETHTOOL_FWVERS_LEN   = 32
	ETHTOOL_BUSINFO_LEN  = 32
	ETHTOOL_EROMVERS_LEN = 32
	ETH_GSTRING_LEN      = 32

	SOPASS_MAX = 6

	PFC_STORM_PREVENTION_AUTO    = 0xffff
	PFC_STORM_PREVENTION_DISABLE = 0

	DOWNSHIFT_DEV_DEFAULT_COUNT = 0xff
	DOWNSHIFT_DEV_DISABLE       = 0
	/* Time in msecs after which link is reported as down
	 * 0 = lowest time supported by the PHY
	 * 0xff = off, link down detection according to standard
	 */
	ETHTOOL_PHY_FAST_LINK_DOWN_ON  = 0
	ETHTOOL_PHY_FAST_LINK_DOWN_OFF = 0xff
	/* Energy Detect Power Down (EDPD) is a feature supported by some PHYs, where
	 * the PHY's RX & TX blocks are put into a low-power mode when there is no
	 * link detected (typically cable is un-plugged). For RX, only a minimal
	 * link-detection is available, and for TX the PHY wakes up to send link pulses
	 * to avoid any lock-ups in case the peer PHY may also be running in EDPD mode.
	 *
	 * Some PHYs may support configuration of the wake-up interval for TX pulses,
	 * and some PHYs may support only disabling TX pulses entirely. For the latter
	 * a special value is required (ETHTOOL_PHY_EDPD_NO_TX) so that this can be
	 * configured from userspace (should the user want it).
	 *
	 * The interval units for TX wake-up are in milliseconds, since this should
	 * cover a reasonable range of intervals:
	 *  - from 1 millisecond, which does not sound like much of a power-saver
	 *  - to ~65 seconds which is quite a lot to wait for a link to come up when
	 *    plugging a cable
	 */
	ETHTOOL_PHY_EDPD_DFLT_TX_MSECS = 0xffff
	ETHTOOL_PHY_EDPD_NO_TX         = 0xfffe
	ETHTOOL_PHY_EDPD_DISABLE       = 0

	ETH_RX_NFC_IP4 = 1
)

type ethtool_cmd struct {
	cmd              uint32
	supported        uint32
	advertising      uint32
	speed            uint16
	duplex           uint8
	port             uint8
	phy_address      uint8
	transceiver      uint8
	autoneg          uint8
	mdio_support     uint8
	maxtxpkt         uint32
	maxrxpkt         uint32
	speed_hi         uint16
	eth_tp_mdix      uint16
	eth_tp_mdix_ctrl uint16
	lp_advertising   uint32
	reserverd        [2]uint32
}

type ethtool_drvinfo struct {
	cmd          uint32
	driver       [32]byte
	version      [32]byte
	fw_version   [ETHTOOL_FWVERS_LEN]byte
	bus_info     [ETHTOOL_BUSINFO_LEN]byte
	erom_version [ETHTOOL_EROMVERS_LEN]byte
	reserved2    [12]byte
	n_priv_flags uint32
	n_stats      uint32
	testinfo_len uint32
	eedump_len   uint32
	regdump_len  uint32
}

type ethtool_wolinfo struct {
	cmd       uint32
	supported int
	wolopts   int
	sopass    [SOPASS_MAX]uint8
}

type ethtool_value struct {
	cmd  uint32
	data uint32
}

const (
	ETHTOOL_ID_UNSPEC uint32 = iota
	ETHTOOL_RX_COPYBREAK
	ETHTOOL_TX_COPYBREAK
	ETHTOOL_PFC_PREVENTION_TOUT
	__ETHTOOL_TUNABLE_COUNT
)

const (
	ETHTOOL_TUNABLE_UNSPEC uint32 = iota
	ETHTOOL_TUNABLE_U8
	ETHTOOL_TUNABLE_U16
	ETHTOOL_TUNABLE_U32
	ETHTOOL_TUNABLE_U64
	ETHTOOL_TUNABLE_STRING
	ETHTOOL_TUNABLE_S8
	ETHTOOL_TUNABLE_S16
	ETHTOOL_TUNABLE_S32
	ETHTOOL_TUNABLE_S64
)

type ethtool_tunable struct {
	cmd     uint32
	id      uint32
	type_id uint32
	len     uint32
	data    [8]uintptr
}

const (
	ETHTOOL_PHY_ID_UNSPEC uint32 = iota
	ETHTOOL_PHY_DOWNSHIFT
	ETHTOOL_PHY_FAST_LINK_DOWN
	ETHTOOL_PHY_EDPD
	__ETHTOOL_PHY_TUNABLE_COUNT
)

type ethtool_regs struct {
	cmd     uint32
	version uint32
	len     uint32
	data    [MAX_DATA_BUF]uint8
}

type ethtool_eeprom struct {
	cmd    uint32
	magic  uint32
	offset uint32
	len    uint32
	data   [MAX_DATA_BUF]uint8
}

type ethtool_eee struct {
	cmd            uint32
	supported      uint32
	advertised     uint32
	lp_advertised  uint32
	eee_active     uint32
	eee_enabled    uint32
	tx_lpi_enabled uint32
	tx_lpi_timer   uint32
	reserved       [2]uint32
}

type ethtool_modinfo struct {
	cmd        uint32
	tp         uint32
	eeprom_len uint32
	reserved   [8]uint32
}

type ethtool_coalesce struct {
	cmd                          uint32
	rx_coalesce_usecs            uint32
	rx_max_coalesced_frames      uint32
	rx_coalesce_usecs_irq        uint32
	rx_max_coalesced_frames_irq  uint32
	tx_coalesce_usecs            uint32
	tx_max_coalesced_frames      uint32
	tx_coalesce_usecs_irq        uint32
	tx_max_coalesced_frames_irq  uint32
	stats_block_coalesce_usecs   uint32
	use_adaptive_rx_coalesce     uint32
	use_adaptive_tx_coalesce     uint32
	pkt_rate_low                 uint32
	rx_coalesce_usecs_low        uint32
	rx_max_coalesced_frames_low  uint32
	tx_coalesce_usecs_low        uint32
	tx_max_coalesced_frames_low  uint32
	pkt_rate_high                uint32
	rx_coalesce_usecs_high       uint32
	rx_max_coalesced_frames_high uint32
	tx_coalesce_usecs_high       uint32
	tx_max_coalesced_frames_high uint32
	rate_sample_interval         uint32
}

type ethtool_ringparam struct {
	cmd                  uint32
	rx_max_pending       uint32
	rx_mini_max_pending  uint32
	rx_jumbo_max_pending uint32
	tx_max_pending       uint32
	rx_pending           uint32
	rx_mini_pending      uint32
	rx_jumbo_pending     uint32
	tx_pending           uint32
}

type ethtool_channels struct {
	cmd            uint32
	max_rx         uint32
	max_tx         uint32
	max_other      uint32
	max_combined   uint32
	rx_count       uint32
	tx_count       uint32
	other_count    uint32
	combined_count uint32
}

type ethtool_pauseparam struct {
	cmd      uint32
	autoneg  uint32
	rx_pause uint32
	tx_pause uint32
}

const (
	ETHTOOL_LINK_EXT_STATE_AUTONEG uint32 = iota
	ETHTOOL_LINK_EXT_STATE_LINK_TRAINING_FAILURE
	ETHTOOL_LINK_EXT_STATE_LINK_LOGICAL_MISMATCH
	ETHTOOL_LINK_EXT_STATE_BAD_SIGNAL_INTEGRITY
	ETHTOOL_LINK_EXT_STATE_NO_CABLE
	ETHTOOL_LINK_EXT_STATE_CABLE_ISSUE
	ETHTOOL_LINK_EXT_STATE_EEPROM_ISSUE
	ETHTOOL_LINK_EXT_STATE_CALIBRATION_FAILURE
	ETHTOOL_LINK_EXT_STATE_POWER_BUDGET_EXCEEDED
	ETHTOOL_LINK_EXT_STATE_OVERHEAT
)

const (
	ETHTOOL_LINK_EXT_SUBSTATE_AN_NO_PARTNER_DETECTED uint32 = iota
	ETHTOOL_LINK_EXT_SUBSTATE_AN_ACK_NOT_RECEIVED
	ETHTOOL_LINK_EXT_SUBSTATE_AN_NEXT_PAGE_EXCHANGE_FAILED
	ETHTOOL_LINK_EXT_SUBSTATE_AN_NO_PARTNER_DETECTED_FORCE_MODE
	ETHTOOL_LINK_EXT_SUBSTATE_AN_FEC_MISMATCH_DURING_OVERRIDE
	ETHTOOL_LINK_EXT_SUBSTATE_AN_NO_HCD
)

const (
	ETHTOOL_LINK_EXT_SUBSTATE_LT_KR_FRAME_LOCK_NOT_ACQUIRED uint32 = iota
	ETHTOOL_LINK_EXT_SUBSTATE_LT_KR_LINK_INHIBIT_TIMEOUT
	ETHTOOL_LINK_EXT_SUBSTATE_LT_KR_LINK_PARTNER_DID_NOT_SET_RECEIVER_READY
	ETHTOOL_LINK_EXT_SUBSTATE_LT_REMOTE_FAULT
)

const (
	ETHTOOL_LINK_EXT_SUBSTATE_LLM_PCS_DID_NOT_ACQUIRE_BLOCK_LOCK uint32 = iota
	ETHTOOL_LINK_EXT_SUBSTATE_LLM_PCS_DID_NOT_ACQUIRE_AM_LOCK
	ETHTOOL_LINK_EXT_SUBSTATE_LLM_PCS_DID_NOT_GET_ALIGN_STATUS
	ETHTOOL_LINK_EXT_SUBSTATE_LLM_FC_FEC_IS_NOT_LOCKED
	ETHTOOL_LINK_EXT_SUBSTATE_LLM_RS_FEC_IS_NOT_LOCKED
)

const (
	ETHTOOL_LINK_EXT_SUBSTATE_BSI_LARGE_NUMBER_OF_PHYSICAL_ERRORS uint32 = iota
	ETHTOOL_LINK_EXT_SUBSTATE_BSI_UNSUPPORTED_RATE
)

const (
	ETHTOOL_LINK_EXT_SUBSTATE_CI_UNSUPPORTED_CABLE uint32 = iota
	ETHTOOL_LINK_EXT_SUBSTATE_CI_CABLE_TEST_FAILURE
)

const (
	ETH_SS_TEST uint32 = iota
	ETH_SS_STATS
	ETH_SS_PRIV_FLAGS
	ETH_SS_NTUPLE_FILTERS
	ETH_SS_FEATURES
	ETH_SS_RSS_HASH_FUNCS
	ETH_SS_TUNABLES
	ETH_SS_PHY_STATS
	ETH_SS_PHY_TUNABLES
	ETH_SS_LINK_MODES
	ETH_SS_MSG_CLASSES
	ETH_SS_WOL_MODES
	ETH_SS_SOF_TIMESTAMPING
	ETH_SS_TS_TX_TYPES
	ETH_SS_TS_RX_FILTERS
	ETH_SS_UDP_TUNNEL_TYPES

	/* add new constants above here */
	ETH_SS_COUNT
)

const MAX_DATA_BUF = 65535

type ethtool_gstrings struct {
	cmd        uint32
	string_set uint32
	len        uint32
	data       [MAX_DATA_BUF]uint8
}

type ethtool_sset_info struct {
	cmd       uint32
	reserved  uint32
	sset_mask uint64

	// data *uint32
}

const (
	ETH_TEST_FL_OFFLINE          = (1 << 0)
	ETH_TEST_FL_FAILED           = (1 << 1)
	ETH_TEST_FL_EXTERNAL_LB      = (1 << 2)
	ETH_TEST_FL_EXTERNAL_LB_DONE = (1 << 3)
)

type ethtool_test struct {
	cmd      uint32
	flags    uint32
	reserved uint32
	len      uint32
	data     [MAX_DATA_BUF]uint64
}

type ethtool_stats struct {
	cmd     uint32
	n_stats uint32
	data    [MAX_DATA_BUF]uint64
}
type ethtool_perm_addr struct {
	cmd  uint32
	size uint32
	data [MAX_DATA_BUF]uint8
}

const (
	ETH_FLAG_TXVLAN = (1 << 7)  /* TX VLAN offload enabled */
	ETH_FLAG_RXVLAN = (1 << 8)  /* RX VLAN offload enabled */
	ETH_FLAG_LRO    = (1 << 15) /* LRO is enabled */
	ETH_FLAG_NTUPLE = (1 << 27) /* N-tuple filters enabled */
	ETH_FLAG_RXHASH = (1 << 28)
)

type ethtool_tcpip4_spec struct {
	ip4src uint32
	ip4dst uint32
	psrc   uint16
	pdst   uint16
	tos    uint8
}

type ethtool_ah_espip4_spec struct {
	ip4src uint32
	ip4dst uint32
	spi    uint32
	tos    uint8
}

type ethtool_usrip4_spec struct {
	ip4src     uint32
	ip4dst     uint32
	l4_4_bytes uint32
	tos        uint8
	ip_ver     uint8
	proto      uint8
}

type ethtool_tcpip6_spec struct {
	ip6src [16]byte
	ip6dst [16]byte
	psrc   uint16
	pdst   uint16
	tclass uint8
}

type ethtool_ah_espip6_spec struct {
	ip6src [16]byte
	ip6dst [16]byte
	spi    uint32
	tclass uint8
}

type ethtool_usrip6_spec struct {
	ip6src     [16]byte
	ip6dst     [16]byte
	l4_4_bytes uint32
	tclass     uint8
	l4_proto   uint8
}

type ethtool_flow_union struct {
	hdata [52]byte
}

type ethtool_flow_ext struct {
	padding    [2]uint8
	h_dest     [6]byte
	vlan_etype uint16
	vlan_tci   uint16
	data       [2]uint32
}

type ethtool_rx_flow_spec struct {
	flow_type   uint32
	h_u         ethtool_flow_union
	h_ext       ethtool_flow_ext
	m_u         ethtool_flow_union
	m_ext       ethtool_flow_ext
	ring_cookie uint64
	location    uint32
}

const (
	ETHTOOL_RX_FLOW_SPEC_RING        = 0x00000000FFFFFFFF
	ETHTOOL_RX_FLOW_SPEC_RING_VF     = 0x000000FF00000000
	ETHTOOL_RX_FLOW_SPEC_RING_VF_OFF = 32
)

func ethtool_get_flow_spec_ring(ring_cookie uint64) uint64 {
	return ETHTOOL_RX_FLOW_SPEC_RING & ring_cookie
}

func ethtool_get_flow_spec_ring_vf(ring_cookie uint64) uint64 {
	return (ETHTOOL_RX_FLOW_SPEC_RING_VF & ring_cookie) >>
		ETHTOOL_RX_FLOW_SPEC_RING_VF_OFF
}

type ethtool_rxnfc struct {
	cmd       uint32
	flow_type uint32
	data      uint64
	fs        ethtool_rx_flow_spec

	rule_cnt uint32 // union with rss_context;

	rule_locs [MAX_DATA_BUF]uint32
}
type ethtool_rxfh_indir struct {
	cmd        uint32
	size       uint32
	ring_index [MAX_DATA_BUF]uint32
}

type ethtool_rxfh struct {
	cmd         uint32
	rss_context uint32
	indir_size  uint32
	key_size    uint32
	hfunc       uint8
	rsvd8       [3]uint8
	rsvd32      uint32
	rss_config  [MAX_DATA_BUF]uint32
}

const (
	ETH_RXFH_CONTEXT_ALLOC   = 0xffffffff
	ETH_RXFH_INDIR_NO_CHANGE = 0xffffffff
)

const (
	ETHTOOL_RXNTUPLE_ACTION_DROP  = -1
	ETHTOOL_RXNTUPLE_ACTION_CLEAR = -2
)

type ethtool_rx_ntuple_flow_spec struct {
	flow_type     uint32
	h_u           [72]byte
	m_u           [72]byte
	vlan_tag      uint16
	vlan_tag_mask uint16
	data          uint64
	data_mask     uint64
	action        int32
}

type ethtool_rx_ntuple struct {
	cmd uint32
	fs  ethtool_rx_ntuple_flow_spec
}

const ETHTOOL_FLASH_MAX_FILENAME = 128
const ETHTOOL_FLASH_ALL_REGIONS = 0

type ethtool_flash struct {
	cmd    uint32
	region uint32
	data   [ETHTOOL_FLASH_MAX_FILENAME]byte
}

type ethtool_dump struct {
	cmd     uint32
	version uint32
	flag    uint32
	len     uint32
	data    [MAX_DATA_BUF]uint8
}

const ETH_FW_DUMP_DISABLE = 0

type ethtool_get_features_block struct {
	available     uint32
	requested     uint32
	active        uint32
	never_changed uint32
}

type ethtool_gfeatures struct {
	cmd      uint32
	size     uint32
	features [1024]ethtool_get_features_block
}

type ethtool_set_features_block struct {
	valid     uint32
	requested uint32
}

type ethtool_sfeatures struct {
	cmd      uint32
	size     uint32
	features [1024]ethtool_set_features_block
}

type ethtool_ts_info struct {
	cmd             uint32
	so_timestamping uint32
	phc_index       int32
	tx_types        uint32
	tx_reserved     [3]uint32
	rx_filters      uint32
	rx_reserved     [3]uint32
}

const (
	ETHTOOL_F_UNSUPPORTED__BIT int = iota
	ETHTOOL_F_WISH__BIT
	ETHTOOL_F_COMPAT__BIT
)

const (
	ETHTOOL_F_UNSUPPORTED = (1 << ETHTOOL_F_UNSUPPORTED__BIT)
	ETHTOOL_F_WISH        = (1 << ETHTOOL_F_WISH__BIT)
	ETHTOOL_F_COMPAT      = (1 << ETHTOOL_F_COMPAT__BIT)
	MAX_NUM_QUEUE         = 4096
)

type ethtool_per_queue_op struct {
	cmd         uint32
	sub_command uint32
	queue_mask  [128]uint32
	data        *byte
}

type ethtool_fecparam struct {
	cmd uint32
	/* bitmask of FEC modes */
	active_fec uint32
	fec        uint32
	reserved   uint32
}

const (
	ETHTOOL_FEC_NONE_BIT int = iota
	ETHTOOL_FEC_AUTO_BIT
	ETHTOOL_FEC_OFF_BIT
	ETHTOOL_FEC_RS_BIT
	ETHTOOL_FEC_BASER_BIT
	ETHTOOL_FEC_LLRS_BIT
)

const (
	ETHTOOL_FEC_NONE  = (1 << ETHTOOL_FEC_NONE_BIT)
	ETHTOOL_FEC_AUTO  = (1 << ETHTOOL_FEC_AUTO_BIT)
	ETHTOOL_FEC_OFF   = (1 << ETHTOOL_FEC_OFF_BIT)
	ETHTOOL_FEC_RS    = (1 << ETHTOOL_FEC_RS_BIT)
	ETHTOOL_FEC_BASER = (1 << ETHTOOL_FEC_BASER_BIT)
	ETHTOOL_FEC_LLRS  = (1 << ETHTOOL_FEC_LLRS_BIT)
)
const (
	/* CMDs currently supported */
	ETHTOOL_GSET = 0x00000001 /* DEPRECATED, Get settings.
	* Please use ETHTOOL_GLINKSETTINGS
	 */
	ETHTOOL_SSET = 0x00000002 /* DEPRECATED, Set settings.
	* Please use ETHTOOL_SLINKSETTINGS
	 */
	ETHTOOL_GDRVINFO = 0x00000003 /* Get driver info. */
	ETHTOOL_GREGS    = 0x00000004 /* Get NIC registers. */
	ETHTOOL_GWOL     = 0x00000005 /* Get wake-on-lan options. */
	ETHTOOL_SWOL     = 0x00000006 /* Set wake-on-lan options. */
	ETHTOOL_GMSGLVL  = 0x00000007 /* Get driver message level */
	ETHTOOL_SMSGLVL  = 0x00000008 /* Set driver msg level. */
	ETHTOOL_NWAY_RST = 0x00000009 /* Restart autonegotiation. */
	/* Get link status for host, i.e. whether the interface *and* the
	 * physical port (if there is one) are up (ethtool_value). */
	ETHTOOL_GLINK       = 0x0000000a
	ETHTOOL_GEEPROM     = 0x0000000b /* Get EEPROM data */
	ETHTOOL_SEEPROM     = 0x0000000c /* Set EEPROM data. */
	ETHTOOL_GCOALESCE   = 0x0000000e /* Get coalesce config */
	ETHTOOL_SCOALESCE   = 0x0000000f /* Set coalesce config. */
	ETHTOOL_GRINGPARAM  = 0x00000010 /* Get ring parameters */
	ETHTOOL_SRINGPARAM  = 0x00000011 /* Set ring parameters. */
	ETHTOOL_GPAUSEPARAM = 0x00000012 /* Get pause parameters */
	ETHTOOL_SPAUSEPARAM = 0x00000013 /* Set pause parameters. */
	ETHTOOL_GRXCSUM     = 0x00000014 /* Get RX hw csum enable (ethtool_value) */
	ETHTOOL_SRXCSUM     = 0x00000015 /* Set RX hw csum enable (ethtool_value) */
	ETHTOOL_GTXCSUM     = 0x00000016 /* Get TX hw csum enable (ethtool_value) */
	ETHTOOL_STXCSUM     = 0x00000017 /* Set TX hw csum enable (ethtool_value) */
	ETHTOOL_GSG         = 0x00000018 /* Get scatter-gather enable
	* (ethtool_value) */
	ETHTOOL_SSG = 0x00000019 /* Set scatter-gather enable
	* (ethtool_value). */
	ETHTOOL_TEST      = 0x0000001a /* execute NIC self-test. */
	ETHTOOL_GSTRINGS  = 0x0000001b /* get specified string set */
	ETHTOOL_PHYS_ID   = 0x0000001c /* identify the NIC */
	ETHTOOL_GSTATS    = 0x0000001d /* get NIC-specific statistics */
	ETHTOOL_GTSO      = 0x0000001e /* Get TSO enable (ethtool_value) */
	ETHTOOL_STSO      = 0x0000001f /* Set TSO enable (ethtool_value) */
	ETHTOOL_GPERMADDR = 0x00000020 /* Get permanent hardware address */
	ETHTOOL_GUFO      = 0x00000021 /* Get UFO enable (ethtool_value) */
	ETHTOOL_SUFO      = 0x00000022 /* Set UFO enable (ethtool_value) */
	ETHTOOL_GGSO      = 0x00000023 /* Get GSO enable (ethtool_value) */
	ETHTOOL_SGSO      = 0x00000024 /* Set GSO enable (ethtool_value) */
	ETHTOOL_GFLAGS    = 0x00000025 /* Get flags bitmap(ethtool_value) */
	ETHTOOL_SFLAGS    = 0x00000026 /* Set flags bitmap(ethtool_value) */
	ETHTOOL_GPFLAGS   = 0x00000027 /* Get driver-private flags bitmap */
	ETHTOOL_SPFLAGS   = 0x00000028 /* Set driver-private flags bitmap */

	ETHTOOL_GRXFH       = 0x00000029 /* Get RX flow hash configuration */
	ETHTOOL_SRXFH       = 0x0000002a /* Set RX flow hash configuration */
	ETHTOOL_GGRO        = 0x0000002b /* Get GRO enable (ethtool_value) */
	ETHTOOL_SGRO        = 0x0000002c /* Set GRO enable (ethtool_value) */
	ETHTOOL_GRXRINGS    = 0x0000002d /* Get RX rings available for LB */
	ETHTOOL_GRXCLSRLCNT = 0x0000002e /* Get RX class rule count */
	ETHTOOL_GRXCLSRULE  = 0x0000002f /* Get RX classification rule */
	ETHTOOL_GRXCLSRLALL = 0x00000030 /* Get all RX classification rule */
	ETHTOOL_SRXCLSRLDEL = 0x00000031 /* Delete RX classification rule */
	ETHTOOL_SRXCLSRLINS = 0x00000032 /* Insert RX classification rule */
	ETHTOOL_FLASHDEV    = 0x00000033 /* Flash firmware to device */
	ETHTOOL_RESET       = 0x00000034 /* Reset hardware */
	ETHTOOL_SRXNTUPLE   = 0x00000035 /* Add an n-tuple filter to device */
	ETHTOOL_GRXNTUPLE   = 0x00000036 /* deprecated */
	ETHTOOL_GSSET_INFO  = 0x00000037 /* Get string set info */
	ETHTOOL_GRXFHINDIR  = 0x00000038 /* Get RX flow hash indir'n table */
	ETHTOOL_SRXFHINDIR  = 0x00000039 /* Set RX flow hash indir'n table */

	ETHTOOL_GFEATURES     = 0x0000003a /* Get device offload settings */
	ETHTOOL_SFEATURES     = 0x0000003b /* Change device offload settings */
	ETHTOOL_GCHANNELS     = 0x0000003c /* Get no of channels */
	ETHTOOL_SCHANNELS     = 0x0000003d /* Set no of channels */
	ETHTOOL_SET_DUMP      = 0x0000003e /* Set dump settings */
	ETHTOOL_GET_DUMP_FLAG = 0x0000003f /* Get dump settings */
	ETHTOOL_GET_DUMP_DATA = 0x00000040 /* Get dump data */
	ETHTOOL_GET_TS_INFO   = 0x00000041 /* Get time stamping and PHC info */
	ETHTOOL_GMODULEINFO   = 0x00000042 /* Get plug-in module information */
	ETHTOOL_GMODULEEEPROM = 0x00000043 /* Get plug-in module eeprom */
	ETHTOOL_GEEE          = 0x00000044 /* Get EEE settings */
	ETHTOOL_SEEE          = 0x00000045 /* Set EEE settings */

	ETHTOOL_GRSSH     = 0x00000046 /* Get RX flow hash configuration */
	ETHTOOL_SRSSH     = 0x00000047 /* Set RX flow hash configuration */
	ETHTOOL_GTUNABLE  = 0x00000048 /* Get tunable configuration */
	ETHTOOL_STUNABLE  = 0x00000049 /* Set tunable configuration */
	ETHTOOL_GPHYSTATS = 0x0000004a /* get PHY-specific statistics */

	ETHTOOL_PERQUEUE = 0x0000004b /* Set per queue options */

	ETHTOOL_GLINKSETTINGS = 0x0000004c /* Get ethtool_link_settings */
	ETHTOOL_SLINKSETTINGS = 0x0000004d /* Set ethtool_link_settings */
	ETHTOOL_PHY_GTUNABLE  = 0x0000004e /* Get PHY tunable configuration */
	ETHTOOL_PHY_STUNABLE  = 0x0000004f /* Set PHY tunable configuration */
	ETHTOOL_GFECPARAM     = 0x00000050 /* Get FEC settings */
	ETHTOOL_SFECPARAM     = 0x00000051 /* Set FEC settings */

	/* compatibility with older code */
	SPARC_ETH_GSET = ETHTOOL_GSET
	SPARC_ETH_SSET = ETHTOOL_SSET
)

/* Link mode bit indices */
const (
	ETHTOOL_LINK_MODE_10baseT_Half_BIT       = 0
	ETHTOOL_LINK_MODE_10baseT_Full_BIT       = 1
	ETHTOOL_LINK_MODE_100baseT_Half_BIT      = 2
	ETHTOOL_LINK_MODE_100baseT_Full_BIT      = 3
	ETHTOOL_LINK_MODE_1000baseT_Half_BIT     = 4
	ETHTOOL_LINK_MODE_1000baseT_Full_BIT     = 5
	ETHTOOL_LINK_MODE_Autoneg_BIT            = 6
	ETHTOOL_LINK_MODE_TP_BIT                 = 7
	ETHTOOL_LINK_MODE_AUI_BIT                = 8
	ETHTOOL_LINK_MODE_MII_BIT                = 9
	ETHTOOL_LINK_MODE_FIBRE_BIT              = 10
	ETHTOOL_LINK_MODE_BNC_BIT                = 11
	ETHTOOL_LINK_MODE_10000baseT_Full_BIT    = 12
	ETHTOOL_LINK_MODE_Pause_BIT              = 13
	ETHTOOL_LINK_MODE_Asym_Pause_BIT         = 14
	ETHTOOL_LINK_MODE_2500baseX_Full_BIT     = 15
	ETHTOOL_LINK_MODE_Backplane_BIT          = 16
	ETHTOOL_LINK_MODE_1000baseKX_Full_BIT    = 17
	ETHTOOL_LINK_MODE_10000baseKX4_Full_BIT  = 18
	ETHTOOL_LINK_MODE_10000baseKR_Full_BIT   = 19
	ETHTOOL_LINK_MODE_10000baseR_FEC_BIT     = 20
	ETHTOOL_LINK_MODE_20000baseMLD2_Full_BIT = 21
	ETHTOOL_LINK_MODE_20000baseKR2_Full_BIT  = 22
	ETHTOOL_LINK_MODE_40000baseKR4_Full_BIT  = 23
	ETHTOOL_LINK_MODE_40000baseCR4_Full_BIT  = 24
	ETHTOOL_LINK_MODE_40000baseSR4_Full_BIT  = 25
	ETHTOOL_LINK_MODE_40000baseLR4_Full_BIT  = 26
	ETHTOOL_LINK_MODE_56000baseKR4_Full_BIT  = 27
	ETHTOOL_LINK_MODE_56000baseCR4_Full_BIT  = 28
	ETHTOOL_LINK_MODE_56000baseSR4_Full_BIT  = 29
	ETHTOOL_LINK_MODE_56000baseLR4_Full_BIT  = 30
	ETHTOOL_LINK_MODE_25000baseCR_Full_BIT   = 31

	/* Last allowed bit for __ETHTOOL_LINK_MODE_LEGACY_MASK is bit
	 * 31. Please do NOT define any SUPPORTED_* or ADVERTISED_*
	 * macro for bits > 31. The only way to use indices > 31 is to
	 * use the new ETHTOOL_GLINKSETTINGS/ETHTOOL_SLINKSETTINGS API.
	 */

	ETHTOOL_LINK_MODE_25000baseKR_Full_BIT       = 32
	ETHTOOL_LINK_MODE_25000baseSR_Full_BIT       = 33
	ETHTOOL_LINK_MODE_50000baseCR2_Full_BIT      = 34
	ETHTOOL_LINK_MODE_50000baseKR2_Full_BIT      = 35
	ETHTOOL_LINK_MODE_100000baseKR4_Full_BIT     = 36
	ETHTOOL_LINK_MODE_100000baseSR4_Full_BIT     = 37
	ETHTOOL_LINK_MODE_100000baseCR4_Full_BIT     = 38
	ETHTOOL_LINK_MODE_100000baseLR4_ER4_Full_BIT = 39
	ETHTOOL_LINK_MODE_50000baseSR2_Full_BIT      = 40
	ETHTOOL_LINK_MODE_1000baseX_Full_BIT         = 41
	ETHTOOL_LINK_MODE_10000baseCR_Full_BIT       = 42
	ETHTOOL_LINK_MODE_10000baseSR_Full_BIT       = 43
	ETHTOOL_LINK_MODE_10000baseLR_Full_BIT       = 44
	ETHTOOL_LINK_MODE_10000baseLRM_Full_BIT      = 45
	ETHTOOL_LINK_MODE_10000baseER_Full_BIT       = 46
	ETHTOOL_LINK_MODE_2500baseT_Full_BIT         = 47
	ETHTOOL_LINK_MODE_5000baseT_Full_BIT         = 48

	ETHTOOL_LINK_MODE_FEC_NONE_BIT                   = 49
	ETHTOOL_LINK_MODE_FEC_RS_BIT                     = 50
	ETHTOOL_LINK_MODE_FEC_BASER_BIT                  = 51
	ETHTOOL_LINK_MODE_50000baseKR_Full_BIT           = 52
	ETHTOOL_LINK_MODE_50000baseSR_Full_BIT           = 53
	ETHTOOL_LINK_MODE_50000baseCR_Full_BIT           = 54
	ETHTOOL_LINK_MODE_50000baseLR_ER_FR_Full_BIT     = 55
	ETHTOOL_LINK_MODE_50000baseDR_Full_BIT           = 56
	ETHTOOL_LINK_MODE_100000baseKR2_Full_BIT         = 57
	ETHTOOL_LINK_MODE_100000baseSR2_Full_BIT         = 58
	ETHTOOL_LINK_MODE_100000baseCR2_Full_BIT         = 59
	ETHTOOL_LINK_MODE_100000baseLR2_ER2_FR2_Full_BIT = 60
	ETHTOOL_LINK_MODE_100000baseDR2_Full_BIT         = 61
	ETHTOOL_LINK_MODE_200000baseKR4_Full_BIT         = 62
	ETHTOOL_LINK_MODE_200000baseSR4_Full_BIT         = 63
	ETHTOOL_LINK_MODE_200000baseLR4_ER4_FR4_Full_BIT = 64
	ETHTOOL_LINK_MODE_200000baseDR4_Full_BIT         = 65
	ETHTOOL_LINK_MODE_200000baseCR4_Full_BIT         = 66
	ETHTOOL_LINK_MODE_100baseT1_Full_BIT             = 67
	ETHTOOL_LINK_MODE_1000baseT1_Full_BIT            = 68
	ETHTOOL_LINK_MODE_400000baseKR8_Full_BIT         = 69
	ETHTOOL_LINK_MODE_400000baseSR8_Full_BIT         = 70
	ETHTOOL_LINK_MODE_400000baseLR8_ER8_FR8_Full_BIT = 71
	ETHTOOL_LINK_MODE_400000baseDR8_Full_BIT         = 72
	ETHTOOL_LINK_MODE_400000baseCR8_Full_BIT         = 73
	ETHTOOL_LINK_MODE_FEC_LLRS_BIT                   = 74
	ETHTOOL_LINK_MODE_100000baseKR_Full_BIT          = 75
	ETHTOOL_LINK_MODE_100000baseSR_Full_BIT          = 76
	ETHTOOL_LINK_MODE_100000baseLR_ER_FR_Full_BIT    = 77
	ETHTOOL_LINK_MODE_100000baseCR_Full_BIT          = 78
	ETHTOOL_LINK_MODE_100000baseDR_Full_BIT          = 79
	ETHTOOL_LINK_MODE_200000baseKR2_Full_BIT         = 80
	ETHTOOL_LINK_MODE_200000baseSR2_Full_BIT         = 81
	ETHTOOL_LINK_MODE_200000baseLR2_ER2_FR2_Full_BIT = 82
	ETHTOOL_LINK_MODE_200000baseDR2_Full_BIT         = 83
	ETHTOOL_LINK_MODE_200000baseCR2_Full_BIT         = 84
	ETHTOOL_LINK_MODE_400000baseKR4_Full_BIT         = 85
	ETHTOOL_LINK_MODE_400000baseSR4_Full_BIT         = 86
	ETHTOOL_LINK_MODE_400000baseLR4_ER4_FR4_Full_BIT = 87
	ETHTOOL_LINK_MODE_400000baseDR4_Full_BIT         = 88
	ETHTOOL_LINK_MODE_400000baseCR4_Full_BIT         = 89
	ETHTOOL_LINK_MODE_100baseFX_Half_BIT             = 90
	ETHTOOL_LINK_MODE_100baseFX_Full_BIT             = 91
	/* must be last entry */
	__ETHTOOL_LINK_MODE_MASK_NBITS = iota
)

/* The following are all involved in forcing a particular link
 * mode for the device for setting things.  When getting the
 * devices settings, these indicate the current mode and whether
 * it was forced up into this mode or autonegotiated.
 */

/* The forced speed, in units of 1Mb. All values 0 to INT_MAX are legal.
 * Update drivers/net/phy/phy.c:phy_speed_to_str() and
 * drivers/net/bonding/bond_3ad.c:__get_link_speed() when adding new values.
 */
const (
	SPEED_10     = 10
	SPEED_100    = 100
	SPEED_1000   = 1000
	SPEED_2500   = 2500
	SPEED_5000   = 5000
	SPEED_10000  = 10000
	SPEED_14000  = 14000
	SPEED_20000  = 20000
	SPEED_25000  = 25000
	SPEED_40000  = 40000
	SPEED_50000  = 50000
	SPEED_56000  = 56000
	SPEED_100000 = 100000
	SPEED_200000 = 200000
	SPEED_400000 = 400000

	SPEED_UNKNOWN = -1
)

const (
	DUPLEX_HALF    = 0x00
	DUPLEX_FULL    = 0x01
	DUPLEX_UNKNOWN = 0xff
)

const (
	MASTER_SLAVE_CFG_UNSUPPORTED      = 0
	MASTER_SLAVE_CFG_UNKNOWN          = 1
	MASTER_SLAVE_CFG_MASTER_PREFERRED = 2
	MASTER_SLAVE_CFG_SLAVE_PREFERRED  = 3
	MASTER_SLAVE_CFG_MASTER_FORCE     = 4
	MASTER_SLAVE_CFG_SLAVE_FORCE      = 5
	MASTER_SLAVE_STATE_UNSUPPORTED    = 0
	MASTER_SLAVE_STATE_UNKNOWN        = 1
	MASTER_SLAVE_STATE_MASTER         = 2
	MASTER_SLAVE_STATE_SLAVE          = 3
	MASTER_SLAVE_STATE_ERR            = 4

	/* Which connector port. */
	PORT_TP    = 0x00
	PORT_AUI   = 0x01
	PORT_MII   = 0x02
	PORT_FIBRE = 0x03
	PORT_BNC   = 0x04
	PORT_DA    = 0x05
	PORT_NONE  = 0xef
	PORT_OTHER = 0xff

	/* Which transceiver to use. */
	XCVR_INTERNAL = 0x00 /* PHY and MAC are in the same package */
	XCVR_EXTERNAL = 0x01 /* PHY and MAC are in different packages */
	XCVR_DUMMY1   = 0x02
	XCVR_DUMMY2   = 0x03
	XCVR_DUMMY3   = 0x04

	/* Enable or disable autonegotiation. */
	AUTONEG_DISABLE = 0x00
	AUTONEG_ENABLE  = 0x01

	/* MDI or MDI-X status/control - if MDI/MDI_X/AUTO is set then
	 * the driver is required to renegotiate link
	 */
	ETH_TP_MDI_INVALID = 0x00 /* status: unknown; control: unsupported */
	ETH_TP_MDI         = 0x01 /* status: MDI;     control: force MDI */
	ETH_TP_MDI_X       = 0x02 /* status: MDI-X;   control: force MDI-X */
	ETH_TP_MDI_AUTO    = 0x03 /*                  control: auto-select */

	/* Wake-On-Lan options. */
	WAKE_PHY         = (1 << 0)
	WAKE_UCAST       = (1 << 1)
	WAKE_MCAST       = (1 << 2)
	WAKE_BCAST       = (1 << 3)
	WAKE_ARP         = (1 << 4)
	WAKE_MAGIC       = (1 << 5)
	WAKE_MAGICSECURE = (1 << 6) /* only meaningful if WAKE_MAGIC */
	WAKE_FILTER      = (1 << 7)

	WOL_MODE_COUNT = 8

	/* L2-L4 network traffic flow types */
	TCP_V4_FLOW    = 0x01 /* hash or spec (tcp_ip4_spec) */
	UDP_V4_FLOW    = 0x02 /* hash or spec (udp_ip4_spec) */
	SCTP_V4_FLOW   = 0x03 /* hash or spec (sctp_ip4_spec) */
	AH_ESP_V4_FLOW = 0x04 /* hash only */
	TCP_V6_FLOW    = 0x05 /* hash or spec (tcp_ip6_spec; nfc only) */
	UDP_V6_FLOW    = 0x06 /* hash or spec (udp_ip6_spec; nfc only) */
	SCTP_V6_FLOW   = 0x07 /* hash or spec (sctp_ip6_spec; nfc only) */
	AH_ESP_V6_FLOW = 0x08 /* hash only */
	AH_V4_FLOW     = 0x09 /* hash or spec (ah_ip4_spec) */
	ESP_V4_FLOW    = 0x0a /* hash or spec (esp_ip4_spec) */
	AH_V6_FLOW     = 0x0b /* hash or spec (ah_ip6_spec; nfc only) */
	ESP_V6_FLOW    = 0x0c /* hash or spec (esp_ip6_spec; nfc only) */
	IPV4_USER_FLOW = 0x0d /* spec only (usr_ip4_spec) */
	IP_USER_FLOW   = IPV4_USER_FLOW
	IPV6_USER_FLOW = 0x0e /* spec only (usr_ip6_spec; nfc only) */
	IPV4_FLOW      = 0x10 /* hash only */
	IPV6_FLOW      = 0x11 /* hash only */
	ETHER_FLOW     = 0x12 /* spec only (ether_spec) */
	/* Flag to enable additional fields in struct ethtool_rx_flow_spec */
	FLOW_EXT     = 0x80000000
	FLOW_MAC_EXT = 0x40000000
	/* Flag to enable RSS spreading of traffic matching rule (nfc only) */
	FLOW_RSS = 0x20000000

	/* L3-L4 network traffic flow hash options */
	RXH_L2DA     = (1 << 1)
	RXH_VLAN     = (1 << 2)
	RXH_L3_PROTO = (1 << 3)
	RXH_IP_SRC   = (1 << 4)
	RXH_IP_DST   = (1 << 5)
	RXH_L4_B_0_1 = (1 << 6) /* src port in case of TCP/UDP/SCTP */
	RXH_L4_B_2_3 = (1 << 7) /* dst port in case of TCP/UDP/SCTP */
	RXH_DISCARD  = (1 << 31)

	RX_CLS_FLOW_DISC = 0xffffffffffffffff
	RX_CLS_FLOW_WAKE = 0xfffffffffffffffe

	/* Special RX classification rule insert location values */
	RX_CLS_LOC_SPECIAL = 0x80000000 /* flag */
	RX_CLS_LOC_ANY     = 0xffffffff
	RX_CLS_LOC_FIRST   = 0xfffffffe
	RX_CLS_LOC_LAST    = 0xfffffffd

	/* EEPROM Standards for plug in modules */
	ETH_MODULE_SFF_8079     = 0x1
	ETH_MODULE_SFF_8079_LEN = 256
	ETH_MODULE_SFF_8472     = 0x2
	ETH_MODULE_SFF_8472_LEN = 512
	ETH_MODULE_SFF_8636     = 0x3
	ETH_MODULE_SFF_8636_LEN = 256
	ETH_MODULE_SFF_8436     = 0x4
	ETH_MODULE_SFF_8436_LEN = 256

	ETH_MODULE_SFF_8636_MAX_LEN = 640
	ETH_MODULE_SFF_8436_MAX_LEN = 640
)

const (
	ETH_RESET_MGMT    = 1 << 0 /* Management processor */
	ETH_RESET_IRQ     = 1 << 1 /* Interrupt requester */
	ETH_RESET_DMA     = 1 << 2 /* DMA engine */
	ETH_RESET_FILTER  = 1 << 3 /* Filtering/flow direction */
	ETH_RESET_OFFLOAD = 1 << 4 /* Protocol offload */
	ETH_RESET_MAC     = 1 << 5 /* Media access controller */
	ETH_RESET_PHY     = 1 << 6 /* Transceiver/PHY */
	ETH_RESET_RAM     = 1 << 7 /* RAM shared between
	 * multiple components */
	ETH_RESET_AP = 1 << 8 /* Application processor */

	ETH_RESET_DEDICATED = 0x0000ffff /* All components dedicated to
	 * this interface */
	ETH_RESET_ALL = 0xffffffff /* All components used by this
	 * interface, even if shared */
)

const ETH_RESET_SHARED_SHIFT = 16

type ethtool_link_settings struct {
	cmd                    uint32
	speed                  uint32
	duplex                 uint8
	port                   uint8
	phy_address            uint8
	autoneg                uint8
	mdio_support           uint8
	eth_tp_mdix            uint8
	eth_tp_mdix_ctrl       uint8
	link_mode_masks_nwords int8
	transceiver            uint8
	master_slave_cfg       uint8
	master_slave_state     uint8
	reserved1              [1]uint8
	reserved               [7]uint32
	link_mode_masks        [0]uint32
}
