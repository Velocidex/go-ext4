package parser

import (
	"io/fs"
	"strings"
	"time"

	"github.com/Velocidex/ordereddict"
)

type FileInfo struct {
	inode      int64
	name       string
	components []string

	mtime time.Time
	atime time.Time
	ctime time.Time
	btime time.Time

	size int64
	mode fs.FileMode

	uid, gid int64
}

func (self *FileInfo) Name() string {
	return self.name
}

func (self *FileInfo) FullPath() string {
	return strings.Join(self.components, "/")
}

func (self *FileInfo) Components() []string {
	return self.components
}

func (self *FileInfo) Inode() int64 {
	return self.inode
}

func (self *FileInfo) Uid() int64 {
	return self.uid
}

func (self *FileInfo) Gid() int64 {
	return self.gid
}

func (self *FileInfo) ModTime() time.Time {
	return self.mtime
}

func (self *FileInfo) Mtime() time.Time {
	return self.mtime
}

func (self *FileInfo) Btime() time.Time {
	return self.btime
}

func (self *FileInfo) Atime() time.Time {
	return self.atime
}

func (self *FileInfo) Ctime() time.Time {
	return self.ctime
}

func (self *FileInfo) Size() int64 {
	return self.size
}

func (self *FileInfo) Mode() fs.FileMode {
	return self.mode
}

func (self *FileInfo) Dict() *ordereddict.Dict {
	return ordereddict.NewDict().
		Set("Name", self.Name()).
		Set("FullPath", self.FullPath()).
		Set("Size", self.Size()).
		Set("Mode", self.Mode()).
		Set("ModeStr", self.Mode().String()).
		Set("ModTime", self.ModTime()).
		Set("Mtime", self.Mtime()).
		Set("Atime", self.Atime()).
		Set("Ctime", self.Ctime()).
		Set("Btime", self.Btime()).
		Set("Data", ordereddict.NewDict().
			Set("Inode", self.inode).
			Set("Uid", self.uid).
			Set("Gid", self.gid))
}

func (self *FileInfo) MarshalJSON() ([]byte, error) {
	return self.Dict().MarshalJSON()
}
