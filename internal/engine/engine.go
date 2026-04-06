package engine

import (
	"math/rand"
	"time"
)

type Engine struct {
	AlgoID  int
	Array   []int
	I, J    int
	Swapped bool
	Done    bool
	Start   time.Time
	Elapsed time.Duration

	// Quick Sort State
	Stack [][2]int // Stack of ranges [low, high]
	Pivot int
	PIdx  int

	// Merge Sort State
	Width int
}

func NewEngine(width, height, algoID int) *Engine {
	rand.Seed(time.Now().UnixNano())
	visWidth := int(float64(width) * 0.7)
	if visWidth > 100 {
		visWidth = 100
	}
	if visWidth < 20 {
		visWidth = 20
	}

	maxHeight := (height - 20) * 8
	if maxHeight < 8 {
		maxHeight = 8
	}

	arr := make([]int, visWidth)
	for i := range arr {
		arr[i] = rand.Intn(maxHeight-2) + 1
	}

	e := &Engine{
		AlgoID: algoID,
		Array:  arr,
		Start:  time.Now(),
	}

	// Initialize Algo Specifics
	switch algoID {
	case 1: // Quick Sort
		e.Stack = append(e.Stack, [2]int{0, len(arr) - 1})
	case 2: // Merge Sort
		e.Width = 1
	}

	return e
}

func (e *Engine) Tick() {
	if e.Done {
		return
	}
	e.Elapsed = time.Since(e.Start)

	switch e.AlgoID {
	case 0: // Bubble Sort
		e.tickBubble()
	case 1: // Quick Sort
		e.tickQuick()
	case 2: // Merge Sort
		e.tickMerge()
	}
}

func (e *Engine) tickBubble() {
	for k := 0; k < 6; k++ {
		if e.I < len(e.Array)-1 {
			if e.J < len(e.Array)-e.I-1 {
				if e.Array[e.J] > e.Array[e.J+1] {
					e.Array[e.J], e.Array[e.J+1] = e.Array[e.J+1], e.Array[e.J]
					e.Swapped = true
				}
				e.J++
			} else {
				if !e.Swapped {
					e.Done = true
					break
				}
				e.J, e.I, e.Swapped = 0, e.I+1, false
			}
		} else {
			e.Done = true
			break
		}
	}
}

func (e *Engine) tickQuick() {
	// Iterative Quick Sort Step
	if len(e.Stack) == 0 {
		e.Done = true
		return
	}

	// Partitioning visualization
	for k := 0; k < 4; k++ {
		if len(e.Stack) == 0 {
			e.Done = true
			return
		}
		top := e.Stack[len(e.Stack)-1]
		low, high := top[0], top[1]

		if low < high {
			// Select pivot (last element)
			if e.PIdx == 0 && e.J == 0 {
				e.Pivot = e.Array[high]
				e.PIdx = low
				e.J = low
			}

			if e.J < high {
				if e.Array[e.J] < e.Pivot {
					e.Array[e.I], e.Array[e.J] = e.Array[e.J], e.Array[e.I]
					e.I++
				}
				e.J++
			} else {
				// Swap pivot into place
				e.Array[e.I], e.Array[high] = e.Array[high], e.Array[e.I]
				p := e.I

				// Pop and push sub-ranges
				e.Stack = e.Stack[:len(e.Stack)-1]
				if p+1 < high {
					e.Stack = append(e.Stack, [2]int{p + 1, high})
				}
				if low < p-1 {
					e.Stack = append(e.Stack, [2]int{low, p - 1})
				}
				// Reset indices
				e.I, e.J, e.PIdx = 0, 0, 0
			}
		} else {
			e.Stack = e.Stack[:len(e.Stack)-1]
		}
	}
}

func (e *Engine) tickMerge() {
	// Bottom-up Iterative Merge Sort
	if e.Width >= len(e.Array) {
		e.Done = true
		return
	}

	for k := 0; k < 2; k++ { // Process fewer steps for Merge as it's more complex per step
		if e.I < len(e.Array) {
			left := e.I
			mid := e.I + e.Width - 1
			right := e.I + 2*e.Width - 1

			if mid >= len(e.Array) {
				e.I = len(e.Array) // Nothing to merge
				break
			}
			if right >= len(e.Array) {
				right = len(e.Array) - 1
			}

			e.merge(left, mid, right)
			e.I += 2 * e.Width
		} else {
			e.I = 0
			e.Width *= 2
		}
	}
}

func (e *Engine) merge(l, m, r int) {
	n1 := m - l + 1
	n2 := r - m

	L := make([]int, n1)
	R := make([]int, n2)

	for i := 0; i < n1; i++ {
		L[i] = e.Array[l+i]
	}
	for j := 0; j < n2; j++ {
		R[j] = e.Array[m+1+j]
	}

	i, j, k := 0, 0, l
	for i < n1 && j < n2 {
		if L[i] <= R[j] {
			e.Array[k] = L[i]
			i++
		} else {
			e.Array[k] = R[j]
			j++
		}
		k++
	}

	for i < n1 {
		e.Array[k] = L[i]
		i++
		k++
	}
	for j < n2 {
		e.Array[k] = R[j]
		j++
		k++
	}
	// Visual marker for active merge
	e.J = l
}
