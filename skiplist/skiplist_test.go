package skiplist

import (
	"math/rand"
	"sort"
	"testing"
)

// BenchmarkSkipList benchmarks the SkipList Get function.
func BenchmarkSkipList_Get(b *testing.B) {
	s := SkipList[int]{Rand: random()}
	numbers := testNumbers(b)
	for _, n := range numbers {
		s.Insert(n)
	}

	b.ResetTimer()

	for j := 0; j < b.N; j++ {
		_, found := s.Get(numbers[0])
		if !found {
			b.Fatal("not found")
		}
	}
}

// BenchmarkBisect benchmarks the sort.Search function.
func BenchmarkSort_Search(b *testing.B) {
	numbers := testNumbers(b)
	valueToFind := numbers[0]
	sort.Ints(numbers)

	b.ResetTimer()

	for j := 0; j < b.N; j++ {
		v := sort.SearchInts(numbers, valueToFind)
		if numbers[v] != valueToFind {
			b.Fatal("not found")
		}
	}
}

func BenchmarkSkipList_Insert(b *testing.B) {
	s := SkipList[int]{Rand: random()}
	numbers := testNumbers(b)

	b.ResetTimer()

	for j := 0; j < b.N; j++ {
		s.Insert(numbers[j%len(numbers)])
	}
}

func BenchmarkSortedInsert(b *testing.B) {
	numbers := testNumbers(b)
	var arr []int

	b.ResetTimer()

	for j := 0; j < b.N; j++ {
		numToInsert := numbers[j%len(numbers)]
		insertionPos := sort.SearchInts(arr, numToInsert)
		arr = append(arr, 0)
		copy(arr[insertionPos+1:], arr[insertionPos:])
		arr[insertionPos] = numToInsert
	}
}

func testNumbers(b *testing.B) []int {
	b.Helper()

	const SIZE = 1_000_000
	r := random()

	var numbers []int
	for i := 0; i < SIZE; i++ {
		numbers = append(numbers, r.Int())
	}

	return numbers
}

func random() *rand.Rand {
	return rand.New(rand.NewSource(42))
}
