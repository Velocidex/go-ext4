package parser

import (
	"errors"
	"fmt"
	"io"
)

var (
	notFoundError = errors.New("Not found")
)

type EXT4Context struct {
	Reader  io.ReaderAt
	Profile *EXT4Profile

	Header *Header

	Is64Bit bool

	BlockCount                int64
	FirstDataBlockOffset      int64
	GroupDescriptorTableCount int64
	GroupDescriptorCount      int64
	BlockSize                 int64
	GroupsPerFlex             uint32
	InodesPerGroup            int64
	InodeSize                 int64

	BlockGroup *BlockGroup
}

func (self *EXT4Context) DebugString() string {
	result := self.Header.DebugString()
	result += fmt.Sprintf("BlockCount %v\n", self.BlockCount)
	result += fmt.Sprintf("GroupDescriptorCount %v\n", self.GroupDescriptorCount)
	result += fmt.Sprintf("BlockSize %v\n", self.BlockSize)
	result += fmt.Sprintf("GroupsPerFlex %v\n", self.GroupsPerFlex)

	return result
}

func (self *EXT4Context) OpenInode(inode int64, components []string) (*Inode, error) {
	res, err := self.BlockGroup.OpenInode(inode)
	if err != nil {
		return nil, err
	}

	res.components = components
	return res, nil
}

func GetEXT4Context(reader io.ReaderAt) (*EXT4Context, error) {
	profile := NewEXT4Profile()

	self := &EXT4Context{
		Reader:  reader,
		Profile: profile,
		Header:  profile.Header(reader, 0),
	}

	sb := self.Header.Superblock()
	if sb.Magic() != 0xEF53 {
		return nil, errors.New("Invalid magic")
	}

	self.BlockSize = int64(1024 << uint(sb.LogBlockSize()))
	self.GroupsPerFlex = 1 << sb.LogGroupPerFlex()
	self.FirstDataBlockOffset = int64(sb.FirstDataBlock()+1) * self.BlockSize
	self.InodesPerGroup = int64(sb.InodePerGroup())
	self.InodeSize = int64(sb.InodeSize())

	self.Is64Bit, _ = sb.FeatureIncompat().Names["FEATURE_INCOMPAT_64BIT"]
	if self.Is64Bit {
		self.BlockCount = (int64(sb.BlockCountHi()) << 32) |
			int64(sb.BlockCountLo())
		self.GroupDescriptorTableCount = self.BlockCount/int64(sb.BlockPerGroup()) + 1
		self.GroupDescriptorCount = (self.GroupDescriptorTableCount * 64 / 1024) + 1

	} else {
		self.BlockCount = int64(sb.BlockCountLo())
		self.GroupDescriptorTableCount = self.BlockCount/int64(sb.BlockPerGroup()) + 1
		self.GroupDescriptorCount = (self.GroupDescriptorTableCount * 32 / 1024) + 1
	}

	bg, err := NewBlockGroup(self, reader)
	if err != nil {
		return nil, err
	}
	self.BlockGroup = bg

	return self, nil
}
