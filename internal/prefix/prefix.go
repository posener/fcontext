package prefix

import "bytes"

const (
	chunkSize = 64
	MaxByte   = 32
)

type Prefix struct {
	content [chunkSize]byte
	i       int
	// bi is index in current byte.
	bi   uint
	rank int
	prev *Prefix
}

func New() *Prefix {
	return nil
}

func (p *Prefix) Append(b byte) *Prefix {
	if b >= MaxByte {
		panic("prefix: exceeded MaxByte")
	}
	p = p.prefixToAppend()
	p.setCurrentByte(b)
	return p
}

func (p *Prefix) prefixToAppend() *Prefix {
	if p == nil || p.i == chunkSize {
		rank := 0
		if p != nil {
			rank = p.rank + 1
		}
		return &Prefix{prev: p, rank: rank}
	}
	// Copy current prefix
	ret := *p
	return &ret
}
func (p *Prefix) setCurrentByte(b byte) {
	p.content[p.i] |= b << p.bi
	p.bi++
	if p.bi == 8 {
		p.bi = 0
		p.i++
	}
}

func (p *Prefix) HasPrefix(q *Prefix) bool {
	if q == nil {
		return true
	}
	if p == nil {
		return false
	}
	if p.rank < q.rank || (p.rank == q.rank && (p.i < q.i || (p.i == q.i && p.bi < q.bi))) {
		return false
	}
	for p.rank > q.rank {
		p = p.prev
	}
	// Here p.rank == q.rank.
	// Iterate over all the bytes blocks and check for prefix match.
	for q != nil {
		if q.i > 0 {
			if !bytes.HasPrefix(p.content[:q.i-1], q.content[:q.i-1]) {
				return false
			}
		}
		// Compare last byte.
		mask := byte((1 << q.bi) - 1)
		if (p.content[q.i] & mask) != (q.content[q.i] & mask) {
			return false
		}
		p, q = p.prev, q.prev
	}
	return true
}
