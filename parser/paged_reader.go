// This reader is needed for reading raw windows devices, such as
// \\.\c: On windows, such devices may only be read using sector
// alignment in whole sector numbers. This reader implements page
// aligned reading and adds pages to an LRU cache to make accessing
// various field members faster.

package parser

import (
	"io"
	"sync"
)

// Keep pages in a free list to avoid allocations.
type FreeList struct {
	mu       sync.Mutex
	pagesize int64

	freelist [][]byte
}

func (self *FreeList) Get() []byte {
	self.mu.Lock()
	defer self.mu.Unlock()

	if len(self.freelist) == 0 {
		return make([]byte, self.pagesize)
	}

	// Take the page off the end of the list
	result := self.freelist[len(self.freelist)-1]
	self.freelist = self.freelist[:len(self.freelist)-1]

	return result
}

func (self *FreeList) Put(in []byte) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.freelist = append(self.freelist, in)
}

type PagedReader struct {
	mu sync.Mutex

	reader   io.ReaderAt
	pagesize int64
	lru      *LRU

	size int64

	freelist *FreeList

	Hits int64
	Miss int64
}

func (self *PagedReader) SetSize(size int64) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.size = size
}

func (self *PagedReader) ReadAt(buf []byte, offset int64) (int, error) {
	self.mu.Lock()
	defer self.mu.Unlock()

	if self.size > 0 {
		if offset >= self.size {
			return 0, io.EOF
		}

		if int64(len(buf))+offset > self.size {
			buf = buf[:self.size-offset]
		}
	}

	buf_idx := 0
	for {
		// How much is left in this page to read?
		to_read := int(self.pagesize - offset%self.pagesize)

		// How much do we need to read into the buffer?
		if to_read > len(buf)-buf_idx {
			to_read = len(buf) - buf_idx
		}

		// Are we done?
		if to_read == 0 {
			return buf_idx, nil
		}

		var page_buf []byte

		page := offset - offset%self.pagesize
		cached_page_buf, pres := self.lru.Get(int(page))
		if !pres {
			self.Miss += 1

			// Read this page into memory - already holding the lock.
			page_buf = self.freelist.Get()
			n, err := self.reader.ReadAt(page_buf, page)
			if err != nil && err != io.EOF {
				return buf_idx, err
			}

			// Only cache full pages.
			if n == int(self.pagesize) {
				self.lru.Add(int(page), page_buf)
			}
		} else {
			self.Hits += 1
			page_buf = cached_page_buf.([]byte)
		}

		// Copy the relevant data from the page.
		page_offset := int(offset % self.pagesize)
		copy(buf[buf_idx:buf_idx+to_read],
			page_buf[page_offset:page_offset+to_read])

		offset += int64(to_read)
		buf_idx += to_read
	}
}

func (self *PagedReader) Flush() {
	self.lru.Purge()

	flusher, ok := self.reader.(Flusher)
	if ok {
		flusher.Flush()
	}
}

func NewPagedReader(reader io.ReaderAt, pagesize int64, cache_size int) (*PagedReader, error) {
	self := &PagedReader{
		reader:   reader,
		pagesize: pagesize,
		freelist: &FreeList{
			pagesize: pagesize,
		},
	}

	// By default 10mb cache.
	cache, err := NewLRU(cache_size, func(key int, value interface{}) {
		// Put the page back on the free list
		self.freelist.Put(value.([]byte))
	}, "NewPagedReader")
	if err != nil {
		return nil, err
	}

	self.lru = cache

	return self, nil
}

// Invalidate the disk cache
type Flusher interface {
	Flush()
}
