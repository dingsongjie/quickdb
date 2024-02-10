package quickdb

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"sync"
)

type QuickDB struct {
	path string
	dir  *os.DirEntry
}

func (instance *QuickDB) Path() string {
	return instance.path
}

func (instance *QuickDB) Initialize() error {
	os.MkdirAll(instance.path, fs.ModeDir)
	return nil
}

func (instance *QuickDB) Open(dbName string) (*Database, error) {
	var db = &Database{}

	db.path = path.Join(instance.path, fmt.Sprintf("%s_%s.db", dbNamePrefix, dbName))
	var err error
	if db.file, err = os.OpenFile(db.path, os.O_RDWR|os.O_CREATE, fs.ModeExclusive); err != nil {
		_ = db.close()
		return nil, err
	}

	if info, err := db.file.Stat(); err != nil {
		return nil, err
	} else if info.Size() == 0 {
		// Initialize new files with meta pages.
		if err := db.init(); err != nil {
			return nil, err
		}
	} else {
		// Read the first meta page to determine the page size.
		var buf [0x1000]byte
		if _, err := db.file.ReadAt(buf[:], 0); err == nil {
			m := db.pageInBuffer(buf[:], 0).meta()
			if err := m.validate(); err != nil {
				// If we can't read the page size, we can assume it's the same
				// as the OS -- since that's how the page size was chosen in the
				// first place.
				//
				// If the first page is invalid and this OS uses a different
				// page size than what the database was created with then we
				// are out of luck and cannot access the database.
				db.pageSize = os.Getpagesize()
			} else {
				db.pageSize = int(m.pageSize)
			}
		}
	}

	// Initialize page pool.
	db.pagePool = sync.Pool{
		New: func() interface{} {
			return make([]byte, db.pageSize)
		},
	}

	// Memory map the data file.
	if err := db.mmap(options.InitialMmapSize); err != nil {
		_ = db.close()
		return nil, err
	}

	// Read in the freelist.
	db.freelist = newFreelist()
	db.freelist.read(db.page(db.meta().freelist))

	// Mark the database as opened and return.
	return db, nil
}
