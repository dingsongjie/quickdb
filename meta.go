package quickdb

import (
	"hash/fnv"
	"unsafe"
)

type meta struct {
	version  uint32
	pageSize uint32
	flags    uint32
	pgid     uint
	checksum uint64
}

func (m *meta) validate() error {
	if m.version != version {
		return ErrVersionMismatch
	} else if m.checksum != 0 && m.checksum != m.sum64() {
		return ErrVersionMismatch
	}
	return nil
}

func (m *meta) sum64() uint64 {
	var h = fnv.New64a()
	_, _ = h.Write((*[unsafe.Offsetof(meta{}.checksum)]byte)(unsafe.Pointer(m))[:])
	return h.Sum64()
}
