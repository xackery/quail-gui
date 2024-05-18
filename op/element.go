package op

import "github.com/xackery/quail/raw"

var (
	root  *Node
	focus *Node
)

type Node struct {
	name     string
	parent   *Node
	isEdited bool
	value    raw.ReadWriter
}

func NewNode(name string, value raw.ReadWriter) *Node {
	return &Node{
		name:  name,
		value: value,
	}
}

func (n *Node) Name() string {
	return n.name
}

func (n *Node) Value() raw.ReadWriter {
	return n.value
}

func (n *Node) Parent() *Node {
	return n.parent
}

func (n *Node) SetParent(parent *Node) {
	n.parent = parent
}

func (n *Node) SetIsEdited(isEdited bool) {
	n.isEdited = isEdited
}

func (n *Node) IsEdited() bool {
	return n.isEdited
}

func Clear() {
	root = nil
	focus = nil
}

func SetRoot(node *Node) {
	root = node
}

func Root() *Node {
	return root
}

func SetFocus(node *Node) {
	focus = node
}

func Focus() *Node {
	return focus
}

func Breadcrumb() string {
	value := ""
	if root != nil {
		value = root.Name()
	}
	if focus != nil {
		if value != "" {
			value += " > "
		}

		if focus.parent != nil {
			value += focus.parent.Name() + " > "
		}

		value = focus.Name()
	}
	return value
}
