package engine

import (
	"math/rand"
	"time"
)

type Engine struct {
	Array   []int
	I, J    int
	Swapped bool
	Done    bool
	Start   time.Time
	Elapsed time.Duration
}

func NewEngine(width, height int) *Engine {
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

	return &Engine{
		Array: arr,
		Start: time.Now(),
	}
}

func (e *Engine) Tick() {
	if e.Done {
		return
	}
	e.Elapsed = time.Since(e.Start)
	
	// Process multiple steps per tick for smoother animation
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
