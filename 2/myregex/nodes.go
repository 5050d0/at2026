package myregex

type node interface {
	children() (node, node)
	reverse() node
	fillNullable(*map[node]*nodeData)
	fillFirst(*map[node]*nodeData)
	fillLast(*map[node]*nodeData)
	buildIndex(i *[]node)
}

type nodeLiteral struct {
	value rune
}

func (n *nodeLiteral) children() (node, node) {
	return nil, nil
}

func (n *nodeLiteral) reverse() node {
	return n
}

type nodeOr struct {
	left, right node
}

func (n *nodeOr) children() (node, node) {
	return n.left, n.right
}

func (n *nodeOr) reverse() node {
	return &nodeOr{n.right.reverse(), n.left.reverse()}
}

type nodeAnd struct {
	left, right node
}

func (n *nodeAnd) children() (node, node) {
	return n.left, n.right
}

func (n *nodeAnd) reverse() node {
	return &nodeAnd{n.right.reverse(), n.left.reverse()}
}

type nodeKleene struct {
	child node
}

func (n *nodeKleene) children() (node, node) {
	return n.child, nil
}

func (n *nodeKleene) reverse() node {
	return &nodeKleene{n.child.reverse()}
}

type nodeSet struct {
	values []rune
}

func (n *nodeSet) children() (node, node) {
	return nil, nil
}

func (n *nodeSet) reverse() node {
	return n
}

type nodeRepeat struct {
	child  node
	number int
}

func (n *nodeRepeat) children() (node, node) {
	return n.child, nil
}

func (n *nodeRepeat) reverse() node {
	return &nodeRepeat{n.child.reverse(), n.number}
}

type nodeGroup struct {
	index int
	child node
}

func (n *nodeGroup) children() (node, node) {
	return n.child, nil
}

func (n *nodeGroup) reverse() node {
	return &nodeGroup{child: n.child.reverse(), index: n.index}
}

type nodeGroupRef struct {
	index int
}

func (n *nodeGroupRef) children() (node, node) {
	return nil, nil
}

func (n *nodeGroupRef) reverse() node {
	return n
}

type nodeEpsilon struct {
}

func (n *nodeEpsilon) children() (node, node) { return nil, nil }
func (n *nodeEpsilon) reverse() node          { return &nodeEpsilon{} }

func hasGroups(n node) bool {
	if n == nil {
		return false
	}

	switch n.(type) {
	case *nodeGroup, *nodeGroupRef:
		return true
	}

	left, right := n.children()
	return hasGroups(left) || hasGroups(right)

}

type nodeEnd struct {
}

func (n *nodeEnd) children() (node, node) {
	return nil, nil
}

func (n *nodeEnd) reverse() node {
	return n
}
