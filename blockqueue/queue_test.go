package blockqueue

import (
	"testing"
)

func TestPermutations(t *testing.T) {
	var N []int
	for i := 0; i < 9; i++ {
		N = append(N, i)
	}

	for n := range N {
		fixed := make([]int, n)
		tryPermutation(t, fixed, n, 0)
	}
}

func tryPermutation(t *testing.T, fixed []int, n, j int) {
	if j == n {
		queue := New()
		for k := 0; k < n; k++ {
			queue.Push(&OrderedBlock{Position: fixed[k]})
		}

		min := 0
		for k := 0; k < n; k++ {
			p := queue.Pop().Position
			if p < min {
				t.Errorf("Inserted %v and got %d as %dth element", fixed, p, k)
			}
			min = p
		}

		if queue.Len() != 0 {
			t.Errorf("Inserted %v and after removing %d elements there are still %d left", fixed, n, queue.Len())
		}
		return
	}

	for i := 0; i < n; i++ {
		fixed[j] = i
		tryPermutation(t, fixed, n, j+1)
	}
}
