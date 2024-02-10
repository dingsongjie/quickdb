package quickdb

import "unsafe"

type branchPageElement struct {
	pos   uint32
	ksize uint32
	pgid  uint
}

func (n *branchPageElement) key() []byte {
	buf := (*[maxAllocSize]byte)(unsafe.Pointer(n))
	return (*[maxAllocSize]byte)(unsafe.Pointer(&buf[n.pos]))[:n.ksize]
}
