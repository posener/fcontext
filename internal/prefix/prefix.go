package prefix

import "bytes"

const size = 64

type Prefix struct {
	content [size]byte
	i       int
	rank    int
	prev    *Prefix
}

func New() *Prefix {
	return nil
}

func (p *Prefix) Append(b byte) *Prefix {
	if p == nil || p.i == size {
		rank := 0
		if p != nil {
			rank = p.rank + 1
		}
		ret := Prefix{prev: p, rank: rank}
		ret.content[ret.i] = b
		ret.i++
		return &ret
	}
	ret := *p
	ret.content[ret.i] = b
	ret.i++
	return &ret
}

func (p *Prefix) HasPrefix(q *Prefix) bool {
	if q == nil {
		return true
	}
	if p == nil {
		return false
	}
	if p.rank < q.rank || (p.rank == q.rank && p.i < q.i) {
		return false
	}
	for p.rank > q.rank {
		p = p.prev
	}
	// Here p.rank == q.rank.
	// Iterate over all the bytes blocks and check for prefix match.
	for q != nil {
		if !bytes.HasPrefix(p.content[:q.i], q.content[:q.i]) {
			return false
		}
		p, q = p.prev, q.prev
	}
	return true
}

func (p *Prefix) len() int {
	return p.rank*size + p.i
}
