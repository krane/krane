package job

type jobHeap []Job

func (h jobHeap) Len() int {
	return len(h)
}
func (h jobHeap) Less(i, j int) bool {
	return h[i].StartTime < h[j].StartTime
}
func (h jobHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *jobHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(Job))
}

func (h *jobHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *jobHeap) Peek() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	return x
}
