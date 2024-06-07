package parser

import (
	"fmt"
	"io"
	"io/fs"
	"time"
)

type Run struct {
	FileOffset int64
	DiskOffset int64
	Length     int64
}

type Inode struct {
	*Inode_

	inode int64

	// The name is only filled when we open the inode through its
	// directory entry.
	name string

	components []string

	file_mode fs.FileMode
}

func (self *Inode) Inode() int64 {
	return self.inode
}

func (self *Inode) DataSize() int64 {
	return int64(self.SizeLo()) | (int64(self.SizeHi()) << 32)
}

func (self *Inode) Stat() *FileInfo {
	return &FileInfo{
		inode:      self.inode,
		name:       self.name,
		components: self.components,

		mtime: time.Unix(int64(self.Mtime()), int64(self.MtimeExtra())>>2),
		atime: time.Unix(int64(self.Atime()), int64(self.AtimeExtra())>>2),
		ctime: time.Unix(int64(self.Ctime()), int64(self.CtimeExtra())>>2),
		btime: time.Unix(int64(self.CRtime()), int64(self.CRtimeExtra())>>2),
		uid:   int64(self.Uid()) | (int64(self.UidHi()) << 16),
		gid:   int64(self.Gid()) | (int64(self.GidHi()) << 16),
		size:  self.DataSize(),
		mode:  fs.FileMode(self.Mode()) | self.file_mode,
	}
}

// Walk the extents tree and collect the runs
func (self *Inode) GetReader(ctx *EXT4Context) (io.ReaderAt, error) {
	runs := self.Runs(ctx)
	reader, err := NewPagedReader(&ExtentReader{
		Reader: self.Reader,
		Size:   self.DataSize(),
		Runs:   runs,
	}, ctx.BlockSize, 100)

	reader.SetSize(self.DataSize())

	return reader, err
}

func (self *Inode) dir(ctx *EXT4Context) ([]DirEntry, error) {
	reader, err := self.GetReader(ctx)
	if err != nil {
		return nil, err
	}

	res := []DirEntry{}
	data_size := self.DataSize()
	for offset := int64(0); offset < data_size; {
		dir_entry := self.Profile.Ext4DirEntry(reader, offset)
		name := ParseString(reader,
			dir_entry.Offset+self.Profile.Off_Ext4DirEntry_Name,
			int64(dir_entry.NameLen()))
		inode_number := int64(dir_entry.Inode())

		file_mode := fs.FileMode(0)
		switch dir_entry.FileTypeInt() {

		// Not a real file skip.
		case 0:
			offset += int64(dir_entry.RecLen())
			continue

		case 2:
			file_mode = fs.ModeDir

		case 3, 4:
			file_mode = fs.ModeDevice

		case 5:
			file_mode = fs.ModeNamedPipe

		case 6:
			file_mode = fs.ModeSocket

		case 7:
			file_mode = fs.ModeSymlink
		}

		res = append(res, DirEntry{
			Name:     name,
			Inode:    inode_number,
			FileMode: file_mode,
		})
		offset += int64(dir_entry.RecLen())
	}

	return res, nil
}

// List the directory
func (self *Inode) Dir(ctx *EXT4Context) ([]*FileInfo, error) {

	dir_entries, err := self.dir(ctx)
	if err != nil {
		return nil, err
	}
	res := []*FileInfo{}
	for _, d := range dir_entries {
		child_inode, err := ctx.BlockGroup.OpenInode(d.Inode)
		if err != nil {
			continue
		}

		file_info := child_inode.Stat()
		file_info.name = d.Name
		file_info.inode = d.Inode
		file_info.components = append(self.components, d.Name)
		file_info.mode |= d.FileMode

		res = append(res, file_info)
	}

	return res, nil
}

func (self *Inode) Runs(ctx *EXT4Context) []Run {
	runs := []Run{}

	// Start the walk at the beginning of the Blocks buffer
	offset := self.Offset + self.Profile.Off_Inode__BlockPointers

	self.visit(offset, &runs, ctx.BlockSize)

	return runs
}

// Walk the extent tree
func (self *Inode) visit(offset int64, runs *[]Run, blocksize int64) {
	extent_header := self.Profile.ExtentHeader(self.Reader, offset)
	if extent_header.Magic() != 0xf30a {
		return
	}

	// Leaf node
	if extent_header.Depth() == 0 {

		for _, e := range extent_header.LeafEntries() {
			*runs = append(*runs, Run{
				FileOffset: int64(e.FirstLogicalBlock()) * blocksize,
				DiskOffset: ((int64(e.StartHi()) << 32) | int64(e.StartLo())) * blocksize,
				Length:     int64(e.Length()) * blocksize,
			})
		}

		// Index nodes - recurse into them
	} else {

		for _, i := range extent_header.IndexEntries() {
			leaf_offset := (int64(i.LeafLo()) | (int64(i.LeafHi()) << 32)) *
				blocksize
			self.visit(leaf_offset, runs, blocksize)
		}
	}
}

func (self *Inode) Debug(ctx *EXT4Context) {
	fmt.Println(self.DebugString())
	for _, r := range self.Runs(ctx) {
		fmt.Printf("Run %#v\n", r)
	}

	dir_entries, err := self.Dir(ctx)
	if err == nil {
		for _, d := range dir_entries {
			fmt.Printf("Child %#v\n", d)
		}
	}
}
