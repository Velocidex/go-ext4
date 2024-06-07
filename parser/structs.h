// From fs/ext4/ext4.h

/*
 * Structure of an inode on the disk
 */
struct ext4_inode {
    __le16i_mode;/* File mode */
    __le16i_uid;/* Low 16 bits of Owner Uid */
    __le32i_size_lo;/* Size in bytes */
    __le32i_atime;/* Access time */
    __le32i_ctime;/* Inode Change time */
    __le32i_mtime;/* Modification time */
    __le32i_dtime;/* Deletion Time */
    __le16i_gid;/* Low 16 bits of Group Id */
    __le16i_links_count;/* Links count */
    __le32i_blocks_lo;/* Blocks count */
    __le32i_flags;/* File flags */
    union {
        struct {
            __le32  l_i_version;
        } linux1;
        struct {
            __u32  h_i_translator;
        } hurd1;
        struct {
            __u32  m_i_reserved1;
        } masix1;
    } osd1;/* OS dependent 1 */
    __le32i_block[EXT4_N_BLOCKS];/* Pointers to blocks */
    __le32i_generation;/* File version (for NFS) */
    __le32i_file_acl_lo;/* File ACL */
    __le32i_size_high;
    __le32i_obso_faddr;/* Obsoleted fragment address */
    union {
        struct {
            __le16l_i_blocks_high; /* were l_i_reserved1 */
            __le16l_i_file_acl_high;
            __le16l_i_uid_high;/* these 2 fields */
            __le16l_i_gid_high;/* were reserved2[0] */
            __le16l_i_checksum_lo;/* crc32c(uuid+inum+inode) LE */
            __le16l_i_reserved;
        } linux2;
        struct {
            __le16h_i_reserved1;/* Obsoleted fragment number/size which are removed in ext4 */
            __u16h_i_mode_high;
            __u16h_i_uid_high;
            __u16h_i_gid_high;
            __u32h_i_author;
        } hurd2;
        struct {
            __le16h_i_reserved1;/* Obsoleted fragment number/size which are removed in ext4 */
            __le16m_i_file_acl_high;
            __u32m_i_reserved2[2];
        } masix2;
    } osd2;/* OS dependent 2 */
    __le16i_extra_isize;
    __le16i_checksum_hi;/* crc32c(uuid+inum+inode) BE */
    __le32  i_ctime_extra;  /* extra Change time      (nsec << 2 | epoch) */
    __le32  i_mtime_extra;  /* extra Modification time(nsec << 2 | epoch) */
    __le32  i_atime_extra;  /* extra Access time      (nsec << 2 | epoch) */
    __le32  i_crtime;       /* File Creation time */
    __le32  i_crtime_extra; /* extra FileCreationtime (nsec << 2 | epoch) */
    __le32  i_version_hi;/* high 32 bits for 64-bit version */
};



/*
 * Each block (leaves and indexes), even inode-stored has header.
 */
struct ext4_extent_header {
    __le16eh_magic;/* probably will support different formats */
    __le16eh_entries;/* number of valid entries */
    __le16eh_max;/* capacity of store in entries */
    __le16eh_depth;/* has tree real underlying blocks? */
    __le32eh_generation;/* generation of the tree */
};


/*
 * This is index on-disk structure.
 * It's used at all the levels except the bottom.
 */
struct ext4_extent_idx {
    __le32ei_block;/* index covers logical blocks from 'block' */
    __le32ei_leaf_lo;/* pointer to the physical block of the next *
                      * level. leaf or next index could be there */
    __le16ei_leaf_hi;/* high 16 bits of physical block */
    __u16ei_unused;
};

/*
 * This is the extent on-disk structure.
 * It's used at the bottom of the tree.
 */
struct ext4_extent {
    __le32ee_block;/* first logical block extent covers */
    __le16ee_len;/* number of blocks covered by extent */
    __le16ee_start_hi;/* high 16 bits of physical block */
    __le32ee_start_lo;/* low 32 bits of physical block */
};


/*
 * The new version of the directory entry.  Since EXT4 structures are
 * stored in intel byte order, and the name_len field could never be
 * bigger than 255 chars, it's safe to reclaim the extra byte for the
 * file_type field.
 */
struct ext4_dir_entry_2 {
    __le32inode;/* Inode number */
    __le16rec_len;/* Directory entry length */
    __u8name_len;/* Name length */
    __u8file_type;
    charname[EXT4_NAME_LEN];/* File name */
};
