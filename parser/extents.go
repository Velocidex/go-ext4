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

	if buf_offset == 0 {
		return 0, io.EOF
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

// nextRun returns the run with the smallest FileOffset that is still greater
// than offset, or nil when no such run exists.
func (self *ExtentReader) nextRun(offset int64) *Run {
	var next *Run
	for i := range self.Runs {
		r := &self.Runs[i]
		if r.FileOffset > offset {
			if next == nil || r.FileOffset < next.FileOffset {
				next = r
			}
		}
	}
	return next
}

func (self *ExtentReader) readParial(buf []byte, offset int64) (int, error) {
	if offset >= self.Size {
		return 0, io.EOF
	}

	run, err := self.getRun(offset)
	if err == nil {
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

	// Sparse hole: fill with zeros up to the next extent or file end.
	holeEnd := self.Size
	if next := self.nextRun(offset); next != nil && next.FileOffset < holeEnd {
		holeEnd = next.FileOffset
	}

	available_bytes := holeEnd - offset
	if available_bytes > int64(len(buf)) {
		available_bytes = int64(len(buf))
	}

	hole := buf[0:available_bytes]
	for i := range hole {
		hole[i] = 0
	}
	return int(available_bytes), nil
}
