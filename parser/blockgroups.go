package parser

import (
	"errors"
	"io"
)

type BlockGroup struct {
	Reader  io.ReaderAt
	Profile *EXT4Profile

	GroupDescriptors []*GroupDescriptor

	BlockSize      int64
	InodesPerGroup int64
	InodeSize      int64
}

// Which block group we need to look at
func (self *BlockGroup) GroupNumber(inode int64) int64 {
	return (inode - 1) / self.InodesPerGroup
}

// The relative inode number in the block group
func (self *BlockGroup) RelativeIndodeNumber(inode int64) int64 {
	return (inode - 1) % self.InodesPerGroup
}

func (self *BlockGroup) GetDescriptor(inode int64) (*GroupDescriptor, error) {
	idx := int(self.GroupNumber(inode))
	if idx < 0 || idx > len(self.GroupDescriptors) {
		return nil, errors.New("Inode out of range!")
	}
	return self.GroupDescriptors[idx], nil
}

func (self *BlockGroup) OpenInode(inode int64) (*Inode, error) {
	// Find the inode inside this block group.
	bg, err := self.GetDescriptor(inode)
	if err != nil {
		return nil, err
	}

	inode_offset := self.BlockSize*bg.InodeTable +
		self.RelativeIndodeNumber(inode)*self.InodeSize

	return &Inode{
		Inode_: self.Profile.Inode_(self.Reader, inode_offset),
		inode:  inode,
	}, nil
}

func NewBlockGroup(ctx *EXT4Context, r io.ReaderAt) (*BlockGroup, error) {
	result := &BlockGroup{
		Profile:        ctx.Profile,
		Reader:         ctx.Reader,
		InodesPerGroup: ctx.InodesPerGroup,
		BlockSize:      ctx.BlockSize,
		InodeSize:      ctx.InodeSize,
	}

	// Figure out how many group descriptors there are
	if ctx.Is64Bit {
		for _, gd := range ParseArray_GroupDescriptor64(ctx.Profile, ctx.Reader,
			ctx.FirstDataBlockOffset, int(ctx.GroupDescriptorTableCount)) {
			result.GroupDescriptors = append(result.GroupDescriptors,
				NewGroupDescriptor64(gd))
		}
	} else {
		for _, gd := range ParseArray_GroupDescriptor32(ctx.Profile, ctx.Reader,
			ctx.FirstDataBlockOffset, int(ctx.GroupDescriptorTableCount)) {
			result.GroupDescriptors = append(result.GroupDescriptors,
				NewGroupDescriptor32(gd))
		}
	}

	return result, nil
}

type GroupDescriptor struct {
	BlockBitmap int64
	InodeBitmap int64
	InodeTable  int64
}

func NewGroupDescriptor32(gd *GroupDescriptor32) *GroupDescriptor {
	return &GroupDescriptor{
		BlockBitmap: int64(gd.BlockBitmapLo()),
		InodeBitmap: int64(gd.InodeBitmapLo()),
		InodeTable:  int64(gd.InodeTableLo()),
	}
}

func NewGroupDescriptor64(gd *GroupDescriptor64) *GroupDescriptor {
	return &GroupDescriptor{
		BlockBitmap: int64(gd.BlockBitmapLo()) | (int64(gd.BlockBitmapHi()) << 32),
		InodeBitmap: int64(gd.InodeBitmapLo()) | (int64(gd.InodeBitmapHi()) << 32),
		InodeTable:  int64(gd.InodeTableLo()) | (int64(gd.InodeTableHi()) << 32),
	}
}
