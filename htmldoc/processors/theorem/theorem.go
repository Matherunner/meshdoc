package theorem

import (
	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshforce/tree"
)

func NewTheoremProcessor() meshdoc.Postprocessor {
	return &theoremProcessor{}
}

type theoremProcessor struct {
}

func (p *theoremProcessor) Process(ctx *meshdoc.Context, r meshdoc.ParsedReader) (meshdoc.ParsedReader, error) {
	// Unwrap the first paragraph of every theorem blocks.

	for _, t := range r.Files() {
		it := tree.NewIterator(t.Root())
		for it.Next(tree.InstructionEnterChild) {
			node := it.Value()
			if it.Exit() {
				continue
			}

			block, ok := node.Value.(*tree.BlockNode)
			if !ok {
				continue
			}
			if block.Name() != "THEOREM" && block.Name() != "PROOF" {
				continue
			}
			if node.Child == nil {
				continue
			}

			child, ok := node.Child.Value.(*tree.BlockNode)
			if !ok {
				continue
			}
			if child.Name() != "P" {
				continue
			}

			if node.Child.Child == nil {
				// No grandchild, so just unlink this P child
				node.Child = node.Child.Right
			} else {
				right := node.Child.Right
				// Find the right most child and link that child to the next node
				cur := node.Child.Child
				for cur.Right != nil {
					cur = cur.Right
				}
				cur.Right = right
				node.Child = node.Child.Child
			}
		}
	}

	return r, nil
}
