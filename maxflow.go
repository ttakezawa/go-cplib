package main

func NewMaxFlow(n int) *MaxFlow {
	return &MaxFlow{
		n: n,
		g: make([][]_MFEdge, n),
	}
}

type MFEdge struct {
	From, To  int
	Cap, Flow int
}
type MaxFlow struct {
	n   int
	pos [][2]int
	g   [][]_MFEdge
}
type _MFEdge struct{ to, rev, cap int }

func (mf *MaxFlow) _smaller(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (mf *MaxFlow) AddEdge(from, to, cap int) int {
	m := len(mf.pos)
	mf.pos = append(mf.pos, [2]int{from, len(mf.g[from])})
	mf.g[from] = append(mf.g[from], _MFEdge{to, len(mf.g[to]), cap})
	mf.g[to] = append(mf.g[to], _MFEdge{from, len(mf.g[from]) - 1, 0})
	return m
}

func (mf *MaxFlow) GetEdge(i int) MFEdge {
	_e := mf.g[mf.pos[i][0]][mf.pos[i][1]]
	_re := mf.g[_e.to][_e.rev]
	return MFEdge{mf.pos[i][0], _e.to, _e.cap + _re.cap, _re.cap}
}

func (mf *MaxFlow) Edges() []MFEdge {
	m := len(mf.pos)
	result := make([]MFEdge, 0, m)
	for i := 0; i < m; i++ {
		result = append(result, mf.GetEdge(i))
	}
	return result
}

func (mf *MaxFlow) ChangeEdge(i, newCap, newFlow int) {
	_e := &mf.g[mf.pos[i][0]][mf.pos[i][1]]
	_re := &mf.g[_e.to][_e.rev]
	_e.cap = newCap - newFlow
	_re.cap = newFlow
}

func (mf *MaxFlow) Flow(s, t int) int {
	return mf.FlowLimit(s, t, int(1e+18))
}

func (mf *MaxFlow) FlowLimit(s, t, flowLim int) int {
	level := make([]int, mf.n)
	iter := make([]int, mf.n)
	bfs := func() {
		for i := range level {
			level[i] = -1
		}
		level[s] = 0
		q := make([]int, 0, mf.n)
		q = append(q, s)
		for len(q) != 0 {
			v := q[0]
			q = q[1:]
			for _, e := range mf.g[v] {
				if e.cap == 0 || level[e.to] >= 0 {
					continue
				}
				level[e.to] = level[v] + 1
				if e.to == t {
					return
				}
				q = append(q, e.to)
			}
		}
	}
	var dfs func(v, up int) int
	dfs = func(v, up int) int {
		if v == s {
			return up
		}
		res := 0
		lv := level[v]
		for ; iter[v] < len(mf.g[v]); iter[v]++ {
			e := &mf.g[v][iter[v]]
			if lv <= level[e.to] || mf.g[e.to][e.rev].cap == 0 {
				continue
			}
			d := dfs(e.to, mf._smaller(up-res, mf.g[e.to][e.rev].cap))
			if d <= 0 {
				continue
			}
			mf.g[v][iter[v]].cap += d
			mf.g[e.to][e.rev].cap -= d
			res += d
			if res == up {
				break
			}
		}
		return res
	}
	flow := 0
	for flow < flowLim {
		bfs()
		if level[t] == -1 {
			break
		}
		for i := range iter {
			iter[i] = 0
		}
		for flow < flowLim {
			f := dfs(t, flowLim-flow)
			if f == 0 {
				break
			}
			flow += f
		}
	}
	return flow
}

func (mf *MaxFlow) MinCut(s int) []bool {
	visited := make([]bool, mf.n)
	q := make([]int, 0, mf.n)
	q = append(q, s)
	for len(q) != 0 {
		p := q[0]
		q = q[1:]
		visited[p] = true
		for _, e := range mf.g[p] {
			if e.cap > 0 && !visited[e.to] {
				visited[e.to] = true
				q = append(q, e.to)
			}
		}
	}
	return visited
}
