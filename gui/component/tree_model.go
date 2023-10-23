package component

import "github.com/xackery/wlk/walk"

type TreeModel struct {
	walk.TreeModelBase
	ref   interface{}
	roots []*TreeNode
}

// NewTreeModel creates a new tree model
func NewTreeModel() *TreeModel {
	tm := new(TreeModel)
	return tm
}

// LazyPopulation returns true if the tree model is lazy
func (tm *TreeModel) LazyPopulation() bool {
	return true
}

// RootCount returns the number of roots
func (tm *TreeModel) RootCount() int {
	return len(tm.roots)
}

// RootAt returns a root at a given index
func (tm *TreeModel) RootAt(index int) walk.TreeItem {
	return tm.roots[index]
}

// RootAdd adds a root to the tree model
func (tm *TreeModel) RootAdd(icon *walk.Bitmap, name string, ref interface{}) *TreeNode {
	root := new(TreeNode)
	root.name = name
	root.icon = icon
	root.ref = ref
	tm.roots = append(tm.roots, root)
	return root
}

// SetRef sets the reference of the tree model
func (tm *TreeModel) SetRef(ref interface{}) {
	tm.ref = ref
}
