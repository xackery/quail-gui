package component

import (
	"github.com/xackery/quail/common"
	"github.com/xackery/wlk/walk"
)

// TreeNode represents an element within a tree view
type TreeNode struct {
	name     string
	parent   *TreeNode
	icon     *walk.Bitmap
	ref      interface{}
	rootRef  interface{}
	children []*TreeNode
	isEdited bool
}

// NewTreeNode creates a new tree node
func NewTreeNode(parent *TreeNode, icon *walk.Bitmap, name string, ref interface{}) *TreeNode {
	tn := new(TreeNode)
	tn.name = name
	tn.parent = parent
	tn.icon = icon
	tn.ref = ref
	return tn
}

// Text returns the text of a tree node
func (tn *TreeNode) Text() string {
	text := tn.name
	if tn.isEdited {
		text += "*"
	}
	return text
}

// Parent returns the parent of a tree node
func (tn *TreeNode) Parent() walk.TreeItem {
	if tn.parent == nil {
		return nil
	}

	return tn.parent
}

// ChildCount returns the number of children a tree node has
func (tn *TreeNode) ChildCount() int {
	if tn.children == nil {
		return 0
	}
	return len(tn.children)
}

// ChildAt returns a child at a given index
func (tn *TreeNode) ChildAt(index int) walk.TreeItem {
	return tn.children[index]
}

// Image returns the icon of a tree node
func (tn *TreeNode) Image() interface{} {
	return tn.icon
}

// ResetChildren resets the children of a tree node
func (tn *TreeNode) ResetChildren() error {
	tn.children = nil

	return nil
}

// ChildAdd adds a child to a tree node
func (tn *TreeNode) ChildAdd(icon *walk.Bitmap, name string, rootRef interface{}, ref interface{}) *TreeNode {
	child := new(TreeNode)
	child.parent = tn
	child.name = name
	child.icon = icon
	child.rootRef = rootRef
	child.ref = ref
	tn.children = append(tn.children, child)
	return child
}

// DuplicateRef duplicates the ref of a tree node
func (tn TreeNode) DuplicateRef() interface{} {
	switch tn.ref.(type) {
	case *common.Layer:
		lay := tn.ref.(*common.Layer)
		layInst := &common.Layer{}
		layInst.Material = lay.Material
		layInst.Diffuse = lay.Diffuse
		layInst.Normal = lay.Normal
		return layInst
	}

	return tn.ref
}

// Ref returns the ref of a tree node
func (tn *TreeNode) Ref() interface{} {
	return tn.ref
}

// SetRef sets the ref of a tree node
func (tn *TreeNode) SetRef(ref interface{}) {
	tn.ref = ref
}

func (tn *TreeNode) RootRef() interface{} {
	return tn.rootRef
}

func (tn *TreeNode) SetRootRef(ref interface{}) {
	tn.rootRef = ref
}

// RemoveChild removes a child from a tree node
func (tn *TreeNode) RemoveChild(child walk.TreeItem) {
	childNode := child.(*TreeNode)

	for i, c := range tn.children {
		if c == childNode {
			tn.children = append(tn.children[:i], tn.children[i+1:]...)
			return
		}
	}
}

// IsEdited returns if a tree node is edited
func (tn *TreeNode) IsEdited() bool {
	return tn.isEdited
}

// Name returns the name of a tree node
func (tn *TreeNode) Name() string {
	return tn.name
}

// SetName sets the name of a tree node
func (tn *TreeNode) SetName(name string) {
	tn.name = name
}

// Icon returns the icon of a tree node
func (tn *TreeNode) Icon() *walk.Bitmap {
	return tn.icon
}

// SetIcon sets the icon of a tree node
func (tn *TreeNode) SetIcon(icon *walk.Bitmap) {
	tn.icon = icon
}
