package cplib

import "container/heap"

// type T = int

func NewPriorityQueue[T any](less func(a, b T) bool) *PriorityQueue[T] {
	return &PriorityQueue[T]{pqHeap[T]{less: less}}
}

type (
	PriorityQueue[T any] struct{ _h pqHeap[T] }
	pqHeap[T any]        struct {
		less func(a, b T) bool
		v    []T
	}
)

func (q PriorityQueue[T]) Len() int       { return len(q._h.v) }
func (q *PriorityQueue[T]) Peek() T       { return q._h.v[0] }
func (q *PriorityQueue[T]) Push(v T)      { heap.Push(&q._h, v) }
func (q *PriorityQueue[T]) Pop() T        { return heap.Pop(&q._h).(T) }
func (h pqHeap[T]) Less(i, j int) bool    { return h.less(h.v[i], h.v[j]) }
func (h pqHeap[T]) Len() int              { return len(h.v) }
func (h pqHeap[T]) Swap(i, j int)         { h.v[i], h.v[j] = h.v[j], h.v[i] }
func (h *pqHeap[T]) Push(v interface{})   { h.v = append(h.v, v.(T)) }
func (h *pqHeap[T]) Pop() (r interface{}) { h.v, r = h.v[:len(h.v)-1], h.v[len(h.v)-1]; return }
