package skiplist

import (
	"cmp"
	"fmt"
	"math/rand"
	"strings"
)

const p = 0.5

// SkipList is a skip list data structure.
// It is a probabilistic data structure that allows for fast search, insert and delete operations.
// It is a sorted data structure, so it is possible to iterate over the elements in order.
// The zero value is ready to use.
type SkipList[T cmp.Ordered] struct {
	head *Node[T]
	Rand *rand.Rand
}

// Get returns the value if it is found in the skip list.
// The second return value is true if the value was found, false otherwise.
func (l *SkipList[T]) Get(searchedVal T) (value T, found bool) {
	var resultNode *Node[T]
	l.walk(l.head, searchedVal, func(node *Node[T]) {
		if found {
			return
		}

		resultNode = node
		found = true
	})

	resultNode = resultNode.Right
	if resultNode == nil || resultNode.Value != searchedVal {
		return value, false
	}
	return resultNode.Value, true
}

// Insert inserts a value into the skip list.
func (l *SkipList[T]) Insert(val T) {
	var (
		previousNode *Node[T]
		shouldInsert = true
	)

	l.walk(l.head, val, func(node *Node[T]) {
		if shouldInsert {
			newNode := Node[T]{Value: val, Right: node.Right, Down: previousNode}
			node.Right = &newNode
			shouldInsert = l.coinFlip()
			previousNode = &newNode
		}
	})

	for shouldInsert {
		newNode := Node[T]{Value: val, Down: previousNode}
		newHead := Node[T]{Right: &newNode, Down: l.head, isHead: true}
		l.head = &newHead

		shouldInsert = l.coinFlip()
		previousNode = &newNode
	}
}

// Remove removes a value from the skip list.
func (l *SkipList[T]) Remove(val T) (removed bool) {
	l.walk(l.head, val, func(current *Node[T]) {
		if current.Right == nil || current.Right.Value != val {
			return
		}

		current.Right = current.Right.Right
	})
	return removed
}

// IterateFrom returns an iterator that starts from the given value.
// The iterator will iterate over all the values in the skip list that are greater or equal to the given value.
func (l *SkipList[T]) IterateFrom(value T) *Iterator[T] {
	var startNode *Node[T]
	l.walk(l.head, value, func(node *Node[T]) {
		if startNode == nil {
			startNode = node
		}
	})

	return &Iterator[T]{skipList: l, current: startNode.Right}
}

// Iterate returns an iterator that starts from the smallest value in the skip list.
func (l *SkipList[T]) Iterate() *Iterator[T] {
	head := l.head
	for head.Down != nil {
		head = head.Down
	}
	return &Iterator[T]{skipList: l, current: head.Right}
}

func (l *SkipList[T]) String() string {
	var result []string
	head := l.head
	for head != nil {
		var row []string
		node := head
		for node != nil {
			s := "HEAD"
			if !node.isHead {
				s = fmt.Sprintf("%v", node.Value)
			}
			row = append(row, s)
			node = node.Right
		}
		result = append(result, strings.Join(row, " -> "))
		head = head.Down
	}

	return strings.Join(result, "\n")
}

// walk is a helper function that walks through the skip list and calls the callback function.
// The callback function is called for each node on the LEFT of the searchedVal, for every level, starting with the deepest level.
// This function is recursive, could theoretically cause stack overflow if the skip list is too big, but that would
// require a skip list that is probably bigger than any computer could handle.
func (l *SkipList[T]) walk(start *Node[T], searchedVal T, callback func(*Node[T])) {
	if start == nil {
		return
	}

	current := start
	for current.Right != nil {
		if current.Right.Value >= searchedVal {
			break
		}
		current = current.Right
	}

	if current.Down == nil {
		callback(current)
		return
	}

	l.walk(current.Down, searchedVal, callback)
	callback(current)
}

func (l *SkipList[T]) coinFlip() bool {
	if l.Rand == nil {
		return rand.Float64() < p
	} else {
		return l.Rand.Float64() < p
	}
}
