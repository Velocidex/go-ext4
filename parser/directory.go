package parser

import (
	"errors"
	"io/fs"
	"strings"
)

type DirEntry struct {
	Inode    int64
	Name     string
	FileMode fs.FileMode
}

func findMatch(entries []DirEntry, name string) (*DirEntry, bool) {
	for _, e := range entries {
		if strings.EqualFold(e.Name, name) {
			return &e, true
		}
	}

	return nil, false
}

func (self *EXT4Context) OpenInodeWithPath(components []string) (*Inode, error) {
	inode_number := ROOT_INODE
	inode, err := self.BlockGroup.OpenInode(inode_number)
	if err != nil {
		return nil, err
	}
	for idx, component := range components {
		if component == "" || component == "." {
			continue
		}

		dir_entries, err := inode.dir(self)
		if err != nil {
			return nil, err
		}

		next_dir_entry, found := findMatch(dir_entries, component)
		if !found {
			return nil, errors.New("Not found")
		}

		next_inode, err := self.BlockGroup.OpenInode(next_dir_entry.Inode)
		if err != nil {
			return nil, errors.New("Not found")
		}

		next_inode.name = next_dir_entry.Name
		next_inode.components = append(components[:idx], next_dir_entry.Name)
		next_inode.inode = next_dir_entry.Inode

		inode = next_inode
	}

	return inode, nil
}
