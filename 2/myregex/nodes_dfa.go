package myregex

func ensureData(m *map[node]*nodeData, n node) *nodeData {
	if val, found := (*m)[n]; found {
		return val
	}
	(*m)[n] = &nodeData{}
	return (*m)[n]
}

func (n *nodeEpsilon) fillNullable(m *map[node]*nodeData) {
	ensureData(m, n).nullable = true
}
func (n *nodeLiteral) fillNullable(m *map[node]*nodeData) {
	ensureData(m, n).nullable = false
}
func (n *nodeSet) fillNullable(m *map[node]*nodeData) {
	ensureData(m, n).nullable = false
}
func (n *nodeEnd) fillNullable(m *map[node]*nodeData) {
	ensureData(m, n).nullable = false
}
func (n *nodeOr) fillNullable(m *map[node]*nodeData) {
	n.left.fillNullable(m)
	n.right.fillNullable(m)
	ensureData(m, n).nullable = ensureData(m, n.left).nullable || ensureData(m, n.right).nullable
}
func (n *nodeAnd) fillNullable(m *map[node]*nodeData) {
	n.left.fillNullable(m)
	n.right.fillNullable(m)
	ensureData(m, n).nullable = ensureData(m, n.left).nullable && ensureData(m, n.right).nullable
}
func (n *nodeKleene) fillNullable(m *map[node]*nodeData) {
	ensureData(m, n).nullable = true
	n.child.fillNullable(m)
}
func (n *nodeRepeat) fillNullable(m *map[node]*nodeData) {
	n.child.fillNullable(m)
	if n.number == 0 {
		ensureData(m, n).nullable = true
	} else {
		ensureData(m, n).nullable = ensureData(m, n.child).nullable
	}
}
func (n *nodeGroup) fillNullable(m *map[node]*nodeData) {
	n.child.fillNullable(m)
	ensureData(m, n).nullable = ensureData(m, n.child).nullable
}
func (n *nodeGroupRef) fillNullable(m *map[node]*nodeData) { panic("unsupported") }

func (n *nodeEpsilon) fillFirst(m *map[node]*nodeData) {
	ensureData(m, n).first = map[int]struct{}{}
}
func (n *nodeLiteral) fillFirst(m *map[node]*nodeData) {
	pos := ensureData(m, n).pos
	ensureData(m, n).first = map[int]struct{}{pos: {}}
}
func (n *nodeSet) fillFirst(m *map[node]*nodeData) {
	pos := ensureData(m, n).pos
	ensureData(m, n).first = map[int]struct{}{pos: {}}
}
func (n *nodeEnd) fillFirst(m *map[node]*nodeData) {
	pos := ensureData(m, n).pos
	ensureData(m, n).first = map[int]struct{}{pos: {}}
}
func (n *nodeOr) fillFirst(m *map[node]*nodeData) {
	n.left.fillFirst(m)
	n.right.fillFirst(m)
	ensureData(m, n).first = unionMaps((*m)[n.left].first, (*m)[n.right].first)
}
func (n *nodeAnd) fillFirst(m *map[node]*nodeData) {
	n.left.fillFirst(m)
	n.right.fillFirst(m)
	first := copyMap((*m)[n.left].first)
	if (*m)[n.left].nullable {
		for k := range (*m)[n.right].first {
			first[k] = struct{}{}
		}
	}
	ensureData(m, n).first = first
}
func (n *nodeKleene) fillFirst(m *map[node]*nodeData) {
	n.child.fillFirst(m)
	ensureData(m, n).first = copyMap((*m)[n.child].first)
}
func (n *nodeRepeat) fillFirst(m *map[node]*nodeData) {
	n.child.fillFirst(m)
	if n.number == 0 {
		ensureData(m, n).first = map[int]struct{}{}
	} else {
		ensureData(m, n).first = copyMap((*m)[n.child].first)
	}
}
func (n *nodeGroup) fillFirst(m *map[node]*nodeData) {
	n.child.fillFirst(m)
	ensureData(m, n).first = copyMap((*m)[n.child].first)
}
func (n *nodeGroupRef) fillFirst(m *map[node]*nodeData) {
	panic("unsupported")
}

func (n *nodeEpsilon) fillLast(m *map[node]*nodeData) {
	ensureData(m, n).last = map[int]struct{}{}
}
func (n *nodeLiteral) fillLast(m *map[node]*nodeData) {
	pos := ensureData(m, n).pos
	ensureData(m, n).last = map[int]struct{}{pos: {}}
}
func (n *nodeSet) fillLast(m *map[node]*nodeData) {
	pos := ensureData(m, n).pos
	ensureData(m, n).last = map[int]struct{}{pos: {}}
}
func (n *nodeEnd) fillLast(m *map[node]*nodeData) {
	pos := ensureData(m, n).pos
	ensureData(m, n).last = map[int]struct{}{pos: {}}
}
func (n *nodeOr) fillLast(m *map[node]*nodeData) {
	n.left.fillLast(m)
	n.right.fillLast(m)
	ensureData(m, n).last = unionMaps((*m)[n.left].last, (*m)[n.right].last)
}
func (n *nodeAnd) fillLast(m *map[node]*nodeData) {
	n.left.fillLast(m)
	n.right.fillLast(m)
	last := copyMap((*m)[n.right].last)
	if (*m)[n.right].nullable {
		for k := range (*m)[n.left].last {
			last[k] = struct{}{}
		}
	}
	ensureData(m, n).last = last
}
func (n *nodeKleene) fillLast(m *map[node]*nodeData) {
	n.child.fillLast(m)
	ensureData(m, n).last = copyMap((*m)[n.child].last)
}
func (n *nodeRepeat) fillLast(m *map[node]*nodeData) {
	n.child.fillLast(m)
	if n.number == 0 {
		ensureData(m, n).last = map[int]struct{}{}
	} else {
		ensureData(m, n).last = copyMap((*m)[n.child].last)
	}
}
func (n *nodeGroup) fillLast(m *map[node]*nodeData) {
	n.child.fillLast(m)
	ensureData(m, n).last = copyMap((*m)[n.child].last)
}

func (n *nodeGroupRef) fillLast(m *map[node]*nodeData) {
	panic("unsupported")
}

func (n *nodeLiteral) buildIndex(i *[]node) {
	*i = append(*i, n)
}
func (n *nodeSet) buildIndex(i *[]node) {
	*i = append(*i, n)
}
func (n *nodeEpsilon) buildIndex(i *[]node) {
	*i = append(*i, n)
}
func (n *nodeEnd) buildIndex(i *[]node) {
	*i = append(*i, n)
}
func (n *nodeOr) buildIndex(i *[]node) {
	n.left.buildIndex(i)
	n.right.buildIndex(i)
}
func (n *nodeAnd) buildIndex(i *[]node) {
	n.left.buildIndex(i)
	n.right.buildIndex(i)
}
func (n *nodeKleene) buildIndex(i *[]node) {
	n.child.buildIndex(i)
}
func (n *nodeRepeat) buildIndex(i *[]node) {
	n.child.buildIndex(i)
}
func (n *nodeGroup) buildIndex(i *[]node) {
	panic("unsupported")
}
func (n *nodeGroupRef) buildIndex(i *[]node) {
	panic("unsupported")
}
