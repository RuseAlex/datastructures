package btree

// Comparator is a function that returns -1 if a < b, 0 if a == b and +1 if a > b
type Comparator func(a, b any) int

type Ordered interface {
	Float | Integer | Character
}

// GenericComparator is a comparator function for the Ordered types (string, int, floats, etc.)
func GenericComparator[T Ordered](a, b T) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint32 | ~uint64 | byte
}

type Character interface {
	~string | ~rune
}

type Float interface {
	~float32 | ~float64
}
