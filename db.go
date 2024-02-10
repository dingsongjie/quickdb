package quickdb

import (
	"os"
	"unsafe"
)

const maxAllocSize = 0x7FFFFFFF
const version = 2

var (
	dbNamePrefix = "qdb"
)

type Database struct {
	dir      *os.DirEntry
	path     string
	file     *os.File
	pageSize int
}

func (db *Database) close() error {
	return nil
}

func (db *Database) init() error {
	db.pageSize = os.Getpagesize()

	// Create two meta pages on a buffer.
	buf := make([]byte, db.pageSize*4)
	for i := 0; i < 2; i++ {
		p := db.pageInBuffer(buf[:], uint(i))
		p.id = uint(i)
		p.flags = uint16(metaPageFlag)

		// Initialize the meta page.
		m := p.meta()
		m.version = version
		m.pageSize = uint32(db.pageSize)
		m.root = bucket{root: 3}
		m.pgid = 4
		m.checksum = m.sum64()
	}
	return nil
}

func (db *Database) pageInBuffer(b []byte, id uint) *page {
	return (*page)(unsafe.Pointer(&b[id*uint(db.pageSize)]))
}
