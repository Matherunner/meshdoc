package counter

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Matherunner/meshdoc"
	"github.com/Matherunner/meshdoc/htmldoc/processors/toc"
	"github.com/Matherunner/meshdoc/utils"
	"github.com/Matherunner/meshforce/tree"
)

var (
	ErrDuplicateElement = errors.New("duplicate element counter")
	ErrInvalidParent    = errors.New("invalid parent element counter")
)

type rootNodeType struct{}
type fileNodeType struct{}

var (
	RootKey = rootNodeType{}
	FileKey = fileNodeType{}
)

type valueNode struct {
	Parent   *valueNode
	Children []*valueNode
	Key      interface{}
	Value    int
}

type valueHierarchy struct {
	hierarchy *Hierarchy
	nodeByKey map[interface{}]*valueNode
	tree      *valueNode
	cur       *valueNode
}

func newValueHierarchy(hierarchy *Hierarchy) *valueHierarchy {
	root := &valueNode{
		Key:   RootKey,
		Value: 1,
	}
	return &valueHierarchy{
		hierarchy: hierarchy,
		nodeByKey: map[interface{}]*valueNode{
			RootKey: root,
		},
		tree: root,
		cur:  root,
	}
}

func (h *valueHierarchy) CurDisplay() (display string) {
	cur := h.cur
	var nums []string
	for cur.Key != RootKey {
		nums = append(nums, strconv.Itoa(cur.Value))
		cur = cur.Parent
	}
	utils.ReverseSlice(nums)
	display = strings.Join(nums, ".")
	return display
}

func (h *valueHierarchy) Increment(key interface{}) (incremented bool) {
	_, ok := h.hierarchy.nodeByKey[key]
	if !ok {
		// If not a numbered element, do nothing.
		return false
	}

	_, ok = h.nodeByKey[key]
	if ok {
		// Found in the map, which means this key is either an ancestor of cur, or cur itself
		// Traverse up the ancestry path, and remove nodes along the way
		for h.cur.Key != key {
			delete(h.nodeByKey, h.cur.Key)
			h.cur = h.cur.Parent
			h.cur.Children = nil
		}
		h.cur.Value++
	} else {
		// Not found in the map, so this is some new child with a common ancestor
		newNode := &valueNode{
			Key:   key,
			Value: 1,
		}
		h.nodeByKey[key] = newNode
		h.cur = newNode

		for {
			parentKey := h.hierarchy.nodeByKey[newNode.Key].Parent.Key
			parentNode, ok := h.nodeByKey[parentKey]
			if ok {
				parentNode.Children = append(parentNode.Children, newNode)
				newNode.Parent = parentNode
				break
			}
			parentNode = &valueNode{
				Children: []*valueNode{newNode},
				Key:      parentKey,
				Value:    1,
			}
			h.nodeByKey[parentKey] = parentNode
			newNode.Parent = parentNode
			newNode = parentNode
		}
	}

	return true
}

type hierarchyNode struct {
	Parent   *hierarchyNode
	Children []*hierarchyNode
	Key      interface{}
}

type Hierarchy struct {
	nodeByKey map[interface{}]*hierarchyNode
	tree      *hierarchyNode
}

func NewHierarchy() *Hierarchy {
	root := &hierarchyNode{
		Key: RootKey,
	}
	return &Hierarchy{
		tree: root,
		nodeByKey: map[interface{}]*hierarchyNode{
			RootKey: root,
		},
	}
}

func (h *Hierarchy) Add(key, parent interface{}) error {
	if _, ok := h.nodeByKey[key]; ok {
		return ErrDuplicateElement
	}

	node, ok := h.nodeByKey[parent]
	if !ok {
		return ErrInvalidParent
	}

	childNode := &hierarchyNode{
		Parent: node,
		Key:    key,
	}
	node.Children = append(node.Children, childNode)
	h.nodeByKey[key] = childNode
	return nil
}

type Options struct {
	Hierarchy *Hierarchy
}

type Counter struct {
	options *Options
}

func NewCounter(options *Options) meshdoc.Postprocessor {
	return &Counter{options: options}
}

func (c *Counter) Process(ctx *meshdoc.Context, r meshdoc.ParsedReader) (meshdoc.ParsedReader, error) {
	toc := toc.FromContext(ctx)

	ctxValue := &ContextValue{
		fileByKey:    map[meshdoc.GenericPath]string{},
		elementByKey: map[*tree.Node]string{},
	}
	valueHierarchy := newValueHierarchy(c.options.Hierarchy)

	for _, entry := range toc {
		incremented := valueHierarchy.Increment(FileKey)
		if incremented {
			display := valueHierarchy.CurDisplay()
			ctxValue.fileByKey[entry] = display
		}

		t := r.Files()[entry]
		it := tree.NewIterator(t)
		for it.Next(tree.InstructionEnterChild) {
			node := it.Value()
			if !it.Exit() {
				switch node := node.Value.(type) {
				case *tree.BlockNode:
					incremented = valueHierarchy.Increment(node.Name())
				case *tree.InlineNode:
					incremented = valueHierarchy.Increment(node.Name())
				default:
					incremented = false
				}
				if incremented {
					display := valueHierarchy.CurDisplay()
					ctxValue.elementByKey[node] = display
				}
			}
		}
	}

	setToContext(ctx, ctxValue)

	return r, nil
}
