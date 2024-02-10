package quickdb

type nodeElement struct {
	flags uint32
	pgid  uint
	key   []byte
	value []byte
}

type nodeElements []nodeElement
