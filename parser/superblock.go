package parser

import (
	"fmt"
	"unsafe"
)

type SuperblockStruct struct {
	InodeCount           uint32     `struc:"uint32,little"`
	BlockCountLo         uint32     `struc:"uint32,little"`
	RBlockCountLo        uint32     `struc:"uint32,little"`
	FreeBlockCountLo     uint32     `struc:"uint32,little"`
	FreeInodeCount       uint32     `struc:"uint32,little"`
	FirstDataBlock       uint32     `struc:"uint32,little"`
	LogBlockSize         uint32     `struc:"uint32,little"`
	LogClusterSize       uint32     `struc:"uint32,little"`
	BlockPerGroup        uint32     `struc:"uint32,little"`
	ClusterPerGroup      uint32     `struc:"uint32,little"` // 10
	InodePerGroup        uint32     `struc:"uint32,little"`
	Mtime                uint32     `struc:"uint32,little"`
	Wtime                uint32     `struc:"uint32,little"`
	MntCount             uint16     `struc:"uint16,little"`
	MaxMntCount          uint16     `struc:"uint16,little"`
	Magic                uint16     `struc:"uint16,little"`
	State                uint16     `struc:"uint16,little"`
	Errors               uint16     `struc:"uint16,little"`
	MinorRevLevel        uint16     `struc:"uint16,little"`
	Lastcheck            uint32     `struc:"uint32,little"`
	Checkinterval        uint32     `struc:"uint32,little"`
	CreatorOs            uint32     `struc:"uint32,little"`
	RevLevel             uint32     `struc:"uint32,little"` // 20
	DefResuid            uint16     `struc:"uint16,little"`
	DefResgid            uint16     `struc:"uint16,little"`
	FirstIno             uint32     `struc:"uint32,little"`
	InodeSize            uint16     `struc:"uint16,little"`
	BlockGroupNr         uint16     `struc:"uint16,little"`
	FeatureCompat        uint32     `struc:"uint32,little"`
	FeatureIncompat      uint32     `struc:"uint32,little"`
	FeatureRoCompat      uint32     `struc:"uint32,little"`
	UUID                 [16]byte   `struc:"[16]byte"` // 30
	VolumeName           [16]byte   `struc:"[16]byte"`
	LastMounted          [64]byte   `struc:"[64]byte"` // 42
	AlgorithmUsageBitmap uint32     `struc:"uint32,little"`
	PreallocBlocks       byte       `struc:"byte"`
	PreallocDirBlocks    byte       `struc:"byte"`
	ReservedGdtBlocks    uint16     `struc:"uint16,little"`
	JournalUUID          [16]byte   `struc:"[16]byte"`
	JournalInum          uint32     `struc:"uint32,little"`
	JournalDev           uint32     `struc:"uint32,little"` // 50
	LastOrphan           uint32     `struc:"uint32,little"`
	HashSeed             [4]uint32  `struc:"[4]uint32,little"`
	DefHashVersion       byte       `struc:"byte"`
	JnlBackupType        byte       `struc:"byte"`
	DescSize             uint16     `struc:"uint16,little"`
	DefaultMountOpts     uint32     `struc:"uint32,little"`
	FirstMetaBg          uint32     `struc:"uint32,little"`
	MkfTime              uint32     `struc:"uint32,little"` // 59
	JnlBlocks            [17]uint32 `struc:"[17]uint32,little"`
	BlockCountHi         uint32     `struc:"uint32,little"`
	RBlockCountHi        uint32     `struc:"uint32,little"`
	FreeBlockCountHi     uint32     `struc:"uint32,little"`
	MinExtraIsize        uint16     `struc:"uint16,little"`
	WantExtraIsize       uint16     `struc:"uint16,little"` // 80
	Flags                uint32     `struc:"uint32,little"`
	RaidStride           uint16     `struc:"uint16,little"`
	MmpUpdateInterval    uint16     `struc:"uint16,little"`
	MmpBlock             uint64     `struc:"uint64,little"`
	RaidStripeWidth      uint32     `struc:"uint32,little"` // 85
	LogGroupPerFlex      byte       `struc:"byte"`
	ChecksumType         byte       `struc:"byte"`
	EncryptionLevel      byte       `struc:"byte"`
	ReservedPad          byte       `struc:"byte"`
	KbyteWritten         uint64     `struc:"uint64,little"`
	SnapshotInum         uint32     `struc:"uint32,little"`
	SnapshotID           uint32     `struc:"uint32,little"` // 90
	SnapshotRBlockCount  uint64     `struc:"uint64,little"`
	SnapshotList         uint32     `struc:"uint32,little"`
	ErrorCount           uint32     `struc:"uint32,little"`
	FirstErrorTime       uint32     `struc:"uint32,little"`
	FirstErrorIno        uint32     `struc:"uint32,little"`
	FirstErrorBlock      uint64     `struc:"uint64,little"` // 98
	FirstErrorFunc       [32]byte   `struc:"[32]pad"`
	FirstErrorLine       uint32     `struc:"uint32,little"`
	LastErrorTime        uint32     `struc:"uint32,little"` // 108
	LastErrorIno         uint32     `struc:"uint32,little"`
	LastErrorLine        uint32     `struc:"uint32,little"`
	LastErrorBlock       uint64     `struc:"uint64,little"` // 112
	LastErrorFunc        [32]byte   `struc:"[32]pad"`
	MountOpts            [64]byte   `struc:"[64]pad"` // 136
	UsrQuotaInum         uint32     `struc:"uint32,little"`
	GrpQuotaInum         uint32     `struc:"uint32,little"`
	OverheadClusters     uint32     `struc:"uint32,little"`
	BackupBgs            [2]uint32  `struc:"[2]uint32,little"`
	EncryptAlgos         [4]byte    `struc:"[4]pad"`
	EncryptPwSalt        [16]byte   `struc:"[16]pad"` // 146
	LpfIno               uint32     `struc:"uint32,little"`
	PrjQuotaInum         uint32     `struc:"uint32,little"`
	ChecksumSeed         uint32     `struc:"uint32,little"`
	Reserved             [98]uint32 `struc:"[98]uint32,little"`
	Checksum             uint32     `struc:"uint32,little"`
}

// GroupDescriptor32 is 32 byte
type GroupDescriptor32_ struct {
	BlockBitmapLo     uint32 `struc:"uint32,little"`
	InodeBitmapLo     uint32 `struc:"uint32,little"`
	InodeTableLo      uint32 `struc:"uint32,little"`
	FreeBlocksCountLo uint16 `struc:"uint16,little"`
	FreeInodesCountLo uint16 `struc:"uint16,little"`
	UsedDirsCountLo   uint16 `struc:"uint16,little"`
	Flags             uint16 `struc:"uint16,little"`
	ExcludeBitmapLo   uint32 `struc:"uint32,little"`
	BlockBitmapCsumLo uint16 `struc:"uint16,little"`
	InodeBitmapCsumLo uint16 `struc:"uint16,little"`
	ItableUnusedLo    uint16 `struc:"uint16,little"`
	Checksum          uint16 `struc:"uint16,little"`
}

// GroupDescriptor is 64 byte
type GroupDescriptor64_ struct {
	GroupDescriptor32_
	BlockBitmapHi     uint32 `struc:"uint32,little"`
	InodeBitmapHi     uint32 `struc:"uint32,little"`
	InodeTableHi      uint32 `struc:"uint32,little"`
	FreeBlocksCountHi uint16 `struc:"uint16,little"`
	FreeInodesCountHi uint16 `struc:"uint16,little"`
	UsedDirsCountHi   uint16 `struc:"uint16,little"`
	ItableUnusedHi    uint16 `struc:"uint16,little"`
	ExcludeBitmapHi   uint32 `struc:"uint32,little"`
	BlockBitmapCsumHi uint16 `struc:"uint16,little"`
	InodeBitmapCsumHi uint16 `struc:"uint16,little"`
	Reserved          uint32 `struc:"uint32,little"`
}

type InodeData struct {
	IMode       uint16 /* File mode */
	IUid        uint16 /* Low 16 bits of Owner Uid */
	ISizeLo     uint32 /* Size in bytes */
	IAtime      uint32 /* Access time */
	ICtime      uint32 /* Inode Change time */
	IMtime      uint32 /* Modification time */
	IDtime      uint32 /* Deletion Time */
	IGid        uint16 /* Low 16 bits of Group Id */
	ILinksCount uint16 /* Links count */
	IBlocksLo   uint32 /* Blocks count */
	IFlags      uint32 /* File flags */

	// union {
	//     struct {
	//         __le32  l_i_version;
	//     } linux1;
	//     struct {
	//         __u32  h_i_translator;
	//     } hurd1;
	//     struct {
	//         __u32  m_i_reserved1;
	//     } masix1;
	// } osd1;             /* OS dependent 1 */
	Osd1 [4]byte

	/*
		IBlock is a general buffer for our data, which can have various
		interpretations. `Ext4NBlocks` comes from the kernel where it is a count in
		terms of uint32's, which is then cast as a struct. However, it works better
		for us as an array of bytes.
	*/
	IBlock [15 * 4]byte

	IGeneration uint32 /* File version (for NFS) */
	IFileAclLo  uint32 /* File ACL */
	ISizeHigh   uint32
	IObsoFaddr  uint32 /* Obsoleted fragment address */

	// union {
	//     struct {
	//         __le16  l_i_blocks_high; /* were l_i_reserved1 */
	//         __le16  l_i_file_acl_high;
	//         __le16  l_i_uid_high;   /* these 2 fields */
	//         __le16  l_i_gid_high;   /* were reserved2[0] */
	//         __le16  l_i_checksum_lo;/* crc32c(uuid+inum+inode) LE */
	//         __le16  l_i_reserved;
	//     } linux2;
	//     struct {
	//         __le16  h_i_reserved1;   Obsoleted fragment number/size which are removed in ext4
	//         __u16   h_i_mode_high;
	//         __u16   h_i_uid_high;
	//         __u16   h_i_gid_high;
	//         __u32   h_i_author;
	//     } hurd2;
	//     struct {
	//         __le16  h_i_reserved1;  /* Obsoleted fragment number/size which are removed in ext4 */
	//         __le16  m_i_file_acl_high;
	//         __u32   m_i_reserved2[2];
	//     } masix2;
	// } osd2;             /* OS dependent 2 */
	Osd2 [12]byte

	IExtraIsize  uint16
	IChecksumHi  uint16 /* crc32c(uuid+inum+inode) BE */
	ICtimeExtra  uint32 /* extra Change time      (nsec << 2 | epoch) */
	IMtimeExtra  uint32 /* extra Modification time(nsec << 2 | epoch) */
	IAtimeExtra  uint32 /* extra Access time      (nsec << 2 | epoch) */
	ICrtime      uint32 /* File Creation time */
	ICrtimeExtra uint32 /* extra FileCreationtime (nsec << 2 | epoch) */
	IVersionHi   uint32 /* high 32 bits for 64-bit version */
	IProjid      uint32 /* Project ID */
}

func DebugStruct() {
	a := SuperblockStruct{}
	fmt.Printf("Offset %v\n", unsafe.Offsetof(a.LogGroupPerFlex))
}
