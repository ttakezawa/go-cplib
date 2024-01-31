package main

type DynamicDSU[S comparable] struct {
	_parent map[S]S
	_size   map[S]int
}

func NewDynamicDSU[S comparable](n int) *DynamicDSU[S] {
	d := &DynamicDSU[S]{}
	d._parent = map[S]S{}
	d._size = map[S]int{}
	return d
}

func (d *DynamicDSU[S]) Leader(a S) S {
	if _, ok := d._parent[a]; !ok {
		d._parent[a] = a
		d._size[a] = 1
		return a
	}
	if d._parent[a] == a {
		return a
	}
	d._parent[a] = d.Leader(d._parent[a])
	return d._parent[a]
}

func (d *DynamicDSU[S]) Merge(a, b S) S {
	x, y := d.Leader(a), d.Leader(b)
	if x == y {
		return x
	}
	if d._size[x] < d._size[y] {
		x, y = y, x
	}
	d._parent[y] = x
	d._size[x] += d._size[y]
	d._size[y] = d._size[x]
	return x
}

func (d *DynamicDSU[S]) Same(a, b S) bool {
	return d.Leader(a) == d.Leader(b)
}

func (d *DynamicDSU[S]) Size(a S) int {
	return d._size[d.Leader(a)]
}

func (d *DynamicDSU[S]) Groups() [][]S {
	leaderToGroup := map[S][]S{}
	for k := range d._parent {
		l := d.Leader(k)
		leaderToGroup[l] = append(leaderToGroup[l], k)
	}
	groups := make([][]S, 0, len(leaderToGroup))
	for _, group := range leaderToGroup {
		groups = append(groups, group)
	}
	return groups
}
