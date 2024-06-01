package skiplist

import "cmp"

type Node[T cmp.Ordered] struct {
	isHead bool

	Value T
	Right *Node[T]
	Down  *Node[T]
}
