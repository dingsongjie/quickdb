package quickdb

import "fmt"

type cursor struct {
	table *table
	stack []elemRef
}

func (c *cursor) Table() *table {
	return c.table
}

func (c *cursor) seek(seek []byte) (key []byte, value []byte, flags uint32) {

	c.stack = c.stack[:0]
	c.search(seek, c.table.root)
	ref := &c.stack[len(c.stack)-1]

	// If the cursor is pointing to the end of page/node then return nil.
	if ref.index >= ref.count() {
		return nil, nil, 0
	}

	// If this is a bucket then return a nil value.
	return c.keyValue()
}
func (c *cursor) search(key []byte, pgid uint) {
	p, n := c.table.pageNode(pgid)
	if p != nil && (p.flags&(branchPageFlag|leafPageFlag)) == 0 {
		panic(fmt.Sprintf("invalid page type: %d: %x", p.id, p.flags))
	}
	e := elemRef{page: p, node: n}
	c.stack = append(c.stack, e)

	// If we're on a leaf page/node then find the specific node.
	if e.isLeaf() {
		c.nsearch(key)
		return
	}

	if n != nil {
		c.searchNode(key, n)
		return
	}
	c.searchPage(key, p)
}

type elemRef struct {
	page  *page
	node  *node
	index int
}
