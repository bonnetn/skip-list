package skiplist

import "cmp"

// Iterator is an iterator for a SkipList.
type Iterator[T cmp.Ordered] struct {
	skipList *SkipList[T]
	current  *Node[T]
}

func (i *Iterator[T]) Next() (val T, done bool) {
	if i.current == nil {
		return val, true
	}

	val = i.current.Value
	i.current = i.current.Right
	return val, false
}
