package treeop

import (
	"fmt"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/slog"
)

// Delete a node
func Delete(node *component.TreeNode) error {
	if node == nil {
		return fmt.Errorf("node is nil")
	}
	if node.Parent() == nil {
		return fmt.Errorf("node parent is nil")
	}
	node.Parent().RemoveChild(node)
	//gui.treeModel.PublishItemsReset(node.Parent())
	slog.Printf("Deleted %+v\n", node)
	return nil
}
