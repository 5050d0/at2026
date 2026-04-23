package myregex

func (n nodeEpsilon) fillNullable(m *map[node]*nodeData) {
	(*m)[n].nullable = false
}

func (n nodeEpsilon) fillFirst(m *map[node]*nodeData) {
	//TODO implement me
	panic("implement me")
}

func (n nodeEpsilon) fillLast(m *map[node]*nodeData) {
	//TODO implement me
	panic("implement me")
}

func (n nodeGroupRef) fillNullable(m *map[node]*nodeData) {
	panic("unsupported")
}

func (n nodeGroupRef) fillFirst(m *map[node]*nodeData) {
	//TODO implement me
	panic("implement me")
}

func (n nodeGroupRef) fillLast(m *map[node]*nodeData) {
	//TODO implement me
	panic("implement me")
}

func (n nodeGroup) fillNullable(m *map[node]*nodeData) {
	n.child.fillNullable(m)
	(*m)[n].nullable = (*m)[n.child].nullable
}

func (n nodeGroup) fillFirst(m *map[node]*nodeData) {
	//TODO implement me
	panic("implement me")
}

func (n nodeGroup) fillLast(m *map[node]*nodeData) {
	//TODO implement me
	panic("implement me")
}

func (n nodeRepeat) fillNullable(m *map[node]*nodeData) {
	n.child.fillNullable(m)
	if n.number == 0 {
		(*m)[n].nullable = true
	} else {
		(*m)[n].nullable = (*m)[n.child].nullable
	}
}

func (n nodeRepeat) fillFirst(m *map[node]*nodeData) {
	//TODO implement me
	panic("implement me")
}

func (n nodeRepeat) fillLast(m *map[node]*nodeData) {
	//TODO implement me
	panic("implement me")
}
func (n nodeSet) fillNullable(m *map[node]*nodeData) {
	(*m)[n].nullable = false
}

func (n nodeSet) fillFirst(m *map[node]*nodeData) {
	//TODO implement me
	panic("implement me")
}

func (n nodeSet) fillLast(m *map[node]*nodeData) {
	//TODO implement me
	panic("implement me")
}
func (n nodeKleene) fillNullable(m *map[node]*nodeData) {
	(*m)[n].nullable = true
	n.child.fillNullable(m)
}

func (n nodeKleene) fillFirst(m *map[node]*nodeData) {
	//TODO implement me
	panic("implement me")
}

func (n nodeKleene) fillLast(m *map[node]*nodeData) {
	//TODO implement me
	panic("implement me")
}
func (n nodeAnd) fillNullable(m *map[node]*nodeData) {
	n.left.fillNullable(m)
	n.right.fillNullable(m)
	(*m)[n].nullable = (*m)[n.left].nullable && (*m)[n.right].nullable
}

func (n nodeAnd) fillFirst(m *map[node]*nodeData) {
	//TODO implement me
	panic("implement me")
}

func (n nodeAnd) fillLast(m *map[node]*nodeData) {
	//TODO implement me
	panic("implement me")
}
func (n nodeOr) fillNullable(m *map[node]*nodeData) {
	n.left.fillNullable(m)
	n.right.fillNullable(m)

	(*m)[n].nullable = (*m)[n.left].nullable || (*m)[n.right].nullable
}

func (n nodeOr) fillFirst(m *map[node]*nodeData) {
	//TODO implement me
	panic("implement me")
}

func (n nodeOr) fillLast(m *map[node]*nodeData) {
	//TODO implement me
	panic("implement me")
}
func (n nodeLiteral) fillNullable(m *map[node]*nodeData) {
	(*m)[n].nullable = false
}

func (n nodeLiteral) fillFirst(m *map[node]*nodeData) {
	//TODO implement me
	panic("implement me")
}

func (n nodeLiteral) fillLast(m *map[node]*nodeData) {
	//TODO implement me
	panic("implement me")
}
func (n nodeEnd) fillNullable(m *map[node]*nodeData) {
	(*m)[n].nullable = false
}

func (n nodeEnd) fillFirst(m *map[node]*nodeData) {
	//TODO implement me
	panic("implement me")
}

func (n nodeEnd) fillLast(m *map[node]*nodeData) {
	//TODO implement me
	panic("implement me")
}
