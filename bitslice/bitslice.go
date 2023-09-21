package bitslice

type Bitslice struct {
	bits []uint64
}

// New creates a new Bitslice of a certain size
func New(size int) *Bitslice {
	return &Bitslice{
		bits: make([]uint64, (size+63)/64),
	}
}

func (b *Bitslice) Get(idx int) bool {
	pos := idx / 64
	j := uint(idx % 64)
	return (b.bits[pos] & (uint64(1) << j)) != 0
}

func (b *Bitslice) Set(idx int, val bool) {
	pos := idx / 64
	j := uint(idx % 64)
	if val {
		b.bits[pos] |= uint64(1) << j
	} else {
		b.bits[pos] &= ^(uint64(1) << j)
	}
}

func (b *Bitslice) Len() int {
	return len(b.bits) * 64
}
