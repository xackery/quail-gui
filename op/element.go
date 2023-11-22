package op

var (
	root  *Node
	focus *Node
)

type Node struct {
	name   string
	parent *Node
	value  interface{}
}

func NewNode(name string, value interface{}) *Node {
	return &Node{
		name:  name,
		value: value,
	}
}

func (n *Node) Name() string {
	return n.name
}

func (n *Node) Value() interface{} {
	return n.value
}

func (n *Node) Parent() *Node {
	return n.parent
}

func (n *Node) SetParent(parent *Node) {
	n.parent = parent
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
