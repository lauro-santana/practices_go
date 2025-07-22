package iterators_test

import (
	"iterators"
	"slices"
	"testing"
)

func TestFilter(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8}
	filter := func(v int) bool { return v%2 == 0 }

	testCases := []struct {
		desc     string
		exp      []int
		testCase func(s []int) func(t *testing.T)
	}{
		{
			desc: "FilterYield",
			exp:  []int{2, 4, 6, 8},
			testCase: func(s []int) func(t *testing.T) {
				return func(t *testing.T) {
					for v := range iterators.FilterYield(slice, filter) {
						if !slices.Contains(s, v) {
							t.Fatal("Expected slice don't have", v)
						}
					}
				}
			},
		},
		{
			desc: "FilterSlice",
			exp:  []int{2, 4, 6, 8},
			testCase: func(s []int) func(t *testing.T) {
				return func(t *testing.T) {
					for _, v := range iterators.FilterSlice(slice, filter) {
						if !slices.Contains(s, v) {
							t.Fatal("Expected slice don't have", v)
						}
					}
				}
			},
		},
		{
			desc: "FilterIterator",
			exp:  []int{2, 4, 6, 8},
			testCase: func(s []int) func(t *testing.T) {
				return func(t *testing.T) {
					iterator := iterators.Iterator[int]{Slice: slice}
					for iterator.Next() {
						v := iterator.FilterIterator(filter)
						if v != nil {
							if !slices.Contains(s, *v) {
								t.Fatal("Expected slice don't have", v)
							}
						}
					}
				}
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, tC.testCase(tC.exp))
	}
}

func Range(size int) []int {
	s := make([]int, size)
	for i := 0; i < size; i++ {
		s[i] = i
	}
	return s
}

func Operation(v int) {}

func SetupSlice() []int {
	return Range(100)
}

func BenchmarkFilterIterator(b *testing.B) {
	filter := func(v int) bool { return v%2 == 0 }

	b.ResetTimer()
	b.StartTimer()
	iterator := iterators.Iterator[int]{Slice: SetupSlice()}

	for iterator.Next() {
		v := iterator.FilterIterator(filter)
		if v != nil {
			Operation(*v)
		}
	}
}

func BenchmarkFilterSlice(b *testing.B) {
	filter := func(v int) bool { return v%2 == 0 }

	b.ResetTimer()
	b.StartTimer()
	for _, v := range iterators.FilterSlice(SetupSlice(), filter) {
		Operation(v)
	}
}

func BenchmarkFilterYield(b *testing.B) {
	filter := func(v int) bool { return v%2 == 0 }

	b.ResetTimer()
	b.StartTimer()
	for v := range iterators.FilterYield(SetupSlice(), filter) {
		Operation(v)
	}
}
