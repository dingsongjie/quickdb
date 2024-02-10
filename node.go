package quickdb

type nodes []node
type node struct {
	isLeaf       bool
	unbalanced   bool
	spilled      bool
	key          []byte
	pgid         uint
	parent       *node
	children     nodes
	nodeElements nodeElements
}

func (n *node) read(p *page) {
	n.pgid = p.id
	n.isLeaf = ((p.flags & leafPageFlag) != 0)
	n.nodeElements = make(nodeElements, int(p.count))

	for i := 0; i < int(p.count); i++ {
		nodeElement := &n.nodeElements[i]
		if n.isLeaf {
			elem := p.leafPageElement(uint16(i))
			nodeElement.flags = elem.flags
			nodeElement.key = elem.key()
			nodeElement.value = elem.value()
		} else {
			elem := p.branchPageElement(uint16(i))
			nodeElement.pgid = elem.pgid
			nodeElement.key = elem.key()
		}
	}

	if len(n.nodeElements) > 0 {
		n.key = n.nodeElements[0].key
	} else {
		n.key = nil
	}
}
