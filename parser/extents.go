package parser

import (
	"io"
)

func (self *ExtentHeader) LeafEntries() []*ExtentEntry {
	return ParseArray_ExtentEntry(self.Profile, self.Reader,
		self.Offset+self.Profile.Off_ExtentHeader__EntryArray,
		int(self.EntryCount()))
}

func (self *ExtentHeader) IndexEntries() []*ExtentIndex {
	return ParseArray_ExtentIndex(self.Profile, self.Reader,
		self.Offset+self.Profile.Off_ExtentHeader__EntryArray,
		int(self.EntryCount()))
}

type ExtentReader struct {
	Reader io.ReaderAt
	Size   int64

	Runs []Run
}

func (self *ExtentReader) ReadAt(buf []byte, offset int64) (int, error) {
	buf_offset := 0
	for buf_offset < len(buf) {
		n, err := self.readParial(buf[buf_offset:], offset+int64(buf_offset))
		if err != nil {
			return buf_offset, err
		}
		if n == 0 {
			break
		}

		buf_offset += n
	}
	return buf_offset, nil
}

func (self *ExtentReader) getRun(offset int64) (*Run, error) {
	for _, r := range self.Runs {
		if r.FileOffset <= offset &&
			r.FileOffset+r.Length > offset {
			return &r, nil
		}
	}

	return nil, io.EOF
}

func (self *ExtentReader) readParial(buf []byte, offset int64) (int, error) {
	run, err := self.getRun(offset)
	if err != nil {
		return 0, err
	}

	available_bytes := run.FileOffset + run.Length - offset
	run_offset := offset - run.FileOffset

	// Do not read past the file end.
	if offset+available_bytes > self.Size {
		available_bytes = self.Size - offset
	}

	if available_bytes > int64(len(buf)) {
		available_bytes = int64(len(buf))
	}

	return self.Reader.ReadAt(
		buf[0:available_bytes], run.DiskOffset+run_offset)
}
