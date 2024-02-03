package cplib

type SCCGraph struct {
	n     int
	edges [][2]int
}

// 強連結成分分解
func NewSCCGraph(n int) *SCCGraph          { return &SCCGraph{n: n} }
func (scc *SCCGraph) AddEdge(from, to int) { scc.edges = append(scc.edges, [2]int{from, to}) }
func (scc *SCCGraph) Groups() [][]int {
	groupNum, ids := scc._sccIds()
	groups := make([][]int, groupNum)
	for i := 0; i < scc.n; i++ {
		groups[ids[i]] = append(groups[ids[i]], i)
	}
	return groups
}

func (scc *SCCGraph) _sccIds() (int, []int) {
	g := NewCSR(scc.n, scc.edges)
	nowOrd, groupNum := 0, 0
	visited, low := make([]int, 0, scc.n), make([]int, scc.n)
	ord, ids := make([]int, scc.n), make([]int, scc.n)
	for i := 0; i < scc.n; i++ {
		ord[i] = -1
	}
	min := func(x, y int) int {
		if x < y {
			return x
		}
		return y
	}
	var dfs func(v int)
	dfs = func(v int) {
		low[v], ord[v] = nowOrd, nowOrd
		nowOrd++
		visited = append(visited, v)
		for i := g.start[v]; i < g.start[v+1]; i++ {
			to := g.elist[i]
			if ord[to] == -1 {
				dfs(to)
				low[v] = min(low[v], low[to])
			} else {
				low[v] = min(low[v], ord[to])
			}
		}
		if low[v] == ord[v] {
			for {
				u := visited[len(visited)-1]
				visited = visited[:len(visited)-1]
				ord[u], ids[u] = scc.n, groupNum
				if u == v {
					break
				}
			}
			groupNum++
		}
	}
	for i := 0; i < scc.n; i++ {
		if ord[i] == -1 {
			dfs(i)
		}
	}
	for i := 0; i < len(ids); i++ {
		ids[i] = groupNum - 1 - ids[i]
	}
	return groupNum, ids
}

type CSR struct{ start, elist []int }

func NewCSR(n int, edges [][2]int) *CSR {
	c := &CSR{start: make([]int, n+1), elist: make([]int, len(edges))}
	for _, e := range edges {
		c.start[e[0]+1]++
	}
	for i := 1; i <= n; i++ {
		c.start[i] += c.start[i-1]
	}
	counter := make([]int, n+1)
	copy(counter, c.start)
	for _, e := range edges {
		c.elist[counter[e[0]]] = e[1]
		counter[e[0]]++
	}
	return c
}
