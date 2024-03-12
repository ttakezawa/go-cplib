package cplib

/*
type EdgeData = int
type S = int
g := NewRerootingDP[EdgeData](N)
for i := 0; i < N-1; i++ {
	a, b := ReadInt()-1, ReadInt()-1
	g.AddEdge(a, b, 0)
}
dp := Calculate[EdgeData, S](
	g,
	// addEdge(x, edge) : 子の部分木にedgeを付ける計算をする
	func(x S, edge RDPEdge[EdgeData]) S {
		return (x + 1) % M
	},
	// merge(x, y) : edge付き部分木であるxとyを合計する
	func(x, y S) S {
		return x * y % M
	},
	// addChildren(x, idx) : 子孫をmergeして得られたxを親ノードidxと合算する
	func(x S, idx int) S {
		return x
	},
	// E: 単位元 merge(x, E) = merge(E, x) = x
	1,
)
*/

type RerootingDP[EdgeData any] struct{ AdjL [][]RDPEdge[EdgeData] }

type RDPEdge[EdgeData any] struct {
	From, To int
	Data     EdgeData
}

func NewRerootingDP[EdgeData any](n int) *RerootingDP[EdgeData] {
	return &RerootingDP[EdgeData]{make([][]RDPEdge[EdgeData], n)}
}

func (g *RerootingDP[EdgeData]) AddEdge(u, v int, data EdgeData) {
	g.AddDirectedEdge(u, v, data)
	g.AddDirectedEdge(v, u, data)
}

func (g *RerootingDP[EdgeData]) AddDirectedEdge(u, v int, data EdgeData) {
	g.AdjL[u] = append(g.AdjL[u], RDPEdge[EdgeData]{From: u, To: v, Data: data})
}

func Calculate[EdgeData, S any](g *RerootingDP[EdgeData], addEdge func(x S, edge RDPEdge[EdgeData]) S, merge func(x, y S) S, addChildren func(x S, idx int) S, E S) []S {
	dp := make([][]S, len(g.AdjL))
	for i := 0; i < len(g.AdjL); i++ {
		dp[i] = make([]S, len(g.AdjL[i]))
	}

	var dfs1 func(p, v int) S
	dfs1 = func(p, v int) S {
		res := E
		for i, edge := range g.AdjL[v] {
			if edge.To == p {
				continue
			}
			dp[v][i] = dfs1(v, edge.To)
			res = merge(res, addEdge(dp[v][i], edge))
		}
		return addChildren(res, v)
	}

	var dfs2 func(p, v int, val S)
	dfs2 = func(p, v int, val S) {
		for i, edge := range g.AdjL[v] {
			if edge.To == p {
				dp[v][i] = val
				break
			}
		}
		pR := make([]S, len(g.AdjL[v])+1)
		pR[len(g.AdjL[v])] = E
		for i := len(g.AdjL[v]); i > 0; i-- {
			pR[i-1] = merge(pR[i], addEdge(dp[v][i-1], g.AdjL[v][i-1]))
		}
		pL := E
		for i, edge := range g.AdjL[v] {
			if edge.To != p {
				dfs2(v, edge.To, addChildren(merge(pL, pR[i+1]), v))
			}
			pL = merge(pL, addEdge(dp[v][i], edge))
		}
	}

	dfs1(-1, 0)
	dfs2(-1, 0, E)
	res := make([]S, len(g.AdjL))
	for v := 0; v < len(g.AdjL); v++ {
		res[v] = E
		for i, edge := range g.AdjL[v] {
			res[v] = merge(res[v], addEdge(dp[v][i], edge))
		}
		res[v] = addChildren(res[v], v)
	}
	return res
}
