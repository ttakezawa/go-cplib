package main

import "container/heap"

type (
	MinCostFlow struct {
		n   int
		pos [][2]int
		g   [][]_MCFEdge
	}
	_MCFEdge struct {
		to   int
		rev  int
		capa int
		cost int
	}
	MCFEdge struct {
		From int
		To   int
		Capa int
		Flow int
		Cost int
	}
)

func NewMinCostFlow(n int) *MinCostFlow {
	return &MinCostFlow{n: n, g: make([][]_MCFEdge, n)}
}

func (mcf *MinCostFlow) AddEdge(from, to, capa, cost int) int {
	m := len(mcf.pos)
	mcf.pos = append(mcf.pos, [2]int{from, len(mcf.g[from])})
	mcf.g[from] = append(mcf.g[from], _MCFEdge{to, len(mcf.g[to]), capa, cost})
	mcf.g[to] = append(mcf.g[to], _MCFEdge{from, len(mcf.g[from]) - 1, 0, -cost})
	return m
}

func (mcf *MinCostFlow) GetEdge(i int) MCFEdge {
	e := mcf.g[mcf.pos[i][0]][mcf.pos[i][1]]
	re := mcf.g[e.to][e.rev]
	return MCFEdge{mcf.pos[i][0], e.to, e.capa + re.capa, re.capa, e.cost}
}

func (mcf *MinCostFlow) Edges() []MCFEdge {
	m := len(mcf.pos)
	res := make([]MCFEdge, m)
	for i := 0; i < m; i++ {
		res[i] = mcf.GetEdge(i)
	}
	return res
}

func (mcf *MinCostFlow) Flow(s, t int) (flow, cost int) {
	res := mcf.Slope(s, t)
	last := res[len(res)-1]
	return last[0], last[1]
}

func (mcf *MinCostFlow) FlowLimit(s, t, flowLim int) (flow, cost int) {
	res := mcf.SlopeLimit(s, t, flowLim)
	last := res[len(res)-1]
	return last[0], last[1]
}

func (mcf *MinCostFlow) Slope(s, t int) [][2]int {
	return mcf.SlopeLimit(s, t, int(1e+18))
}

func (mcf *MinCostFlow) SlopeLimit(s, t, flowLim int) [][2]int {
	dual, dist := make([]int, mcf.n), make([]int, mcf.n)
	pv, pe := make([]int, mcf.n), make([]int, mcf.n)
	vis := make([]bool, mcf.n)
	dualRef := func() bool {
		for i := 0; i < mcf.n; i++ {
			dist[i], pv[i], pe[i] = int(1e+18), -1, -1
			vis[i] = false
		}
		pq := MCFHeapq{}
		dist[s] = 0
		pq.push(&MCFItem{value: s, priority: 0})
		for pq.Len() != 0 {
			v := heap.Pop(&pq).(*MCFItem).value
			if vis[v] {
				continue
			}
			vis[v] = true
			if v == t {
				break
			}
			for i := 0; i < len(mcf.g[v]); i++ {
				e := mcf.g[v][i]
				if vis[e.to] || e.capa == 0 {
					continue
				}
				cost := e.cost - dual[e.to] + dual[v]
				if dist[e.to]-dist[v] > cost {
					dist[e.to] = dist[v] + cost
					pv[e.to] = v
					pe[e.to] = i
					pq.push(&MCFItem{value: e.to, priority: dist[e.to]})
				}
			}
		}
		if !vis[t] {
			return false
		}
		for v := 0; v < mcf.n; v++ {
			if !vis[v] {
				continue
			}
			dual[v] -= dist[t] - dist[v]
		}
		return true
	}
	flow, cost, prevCost := 0, 0, -1
	res := make([][2]int, 0, mcf.n)
	res = append(res, [2]int{flow, cost})
	for flow < flowLim {
		if !dualRef() {
			break
		}
		c := flowLim - flow
		for v := t; v != s; v = pv[v] {
			c = mcf._min(c, mcf.g[pv[v]][pe[v]].capa)
		}
		for v := t; v != s; v = pv[v] {
			mcf.g[pv[v]][pe[v]].capa -= c
			mcf.g[v][mcf.g[pv[v]][pe[v]].rev].capa += c
		}
		d := -dual[s]
		flow += c
		cost += c * d
		if prevCost == d {
			res = res[:len(res)-1]
		}
		res = append(res, [2]int{flow, cost})
		prevCost = cost
	}
	return res
}

func (mcf *MinCostFlow) _min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type (
	MCFItem struct {
		value    int
		priority int
		index    int
	}
	MCFHeapq []*MCFItem
)

func (q MCFHeapq) Less(i, j int) bool    { return q[i].priority < q[j].priority }
func (q MCFHeapq) Len() int              { return len(q) }
func (q MCFHeapq) Swap(i, j int)         { q[i], q[j] = q[j], q[i] }
func (q *MCFHeapq) Push(x interface{})   { *q = append(*q, x.(*MCFItem)) }
func (q *MCFHeapq) Pop() (x interface{}) { *q, x = (*q)[:len(*q)-1], (*q)[len(*q)-1]; return }
func (q *MCFHeapq) push(v *MCFItem)      { heap.Push(q, v) }
func (q *MCFHeapq) pop() *MCFItem        { return heap.Pop(q).(*MCFItem) }
