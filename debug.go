package rbt

// DebugBlock is a simple debug block
type DebugBlock struct {
	Color     color     `json:"color"`
	ChildType childType `json:"childType"`

	Parent   string         `json:"parent"`
	Children [2]*DebugBlock `json:"children"`

	Key string `json:"key"`
}

// GetDebug will get the debug block tree
func GetDebug(t *Tree) (b *DebugBlock) {
	return getDebugBlock(t, t.t.root)
}

func getDebugBlock(t *Tree, index int64) *DebugBlock {
	blk := t.getBlock(index)
	if blk == nil {
		return nil
	}

	var b DebugBlock
	b.Key = string(t.getKey(blk))
	b.ChildType = blk.ct
	b.Color = blk.c
	b.Children[0] = getDebugBlock(t, blk.children[0])
	b.Children[1] = getDebugBlock(t, blk.children[1])
	if parent := t.getBlock(blk.parent); parent != nil {
		b.Parent = string(t.getKey(parent))
	}

	return &b
}
