package whiskey

// DebugBlock is a simple debug block
type DebugBlock struct {
	Color     color     `json:"color"`
	ChildType childType `json:"childType"`

	Parent   string         `json:"parent"`
	Children [2]*DebugBlock `json:"children"`

	Key string `json:"key"`
}

// GetDebug will get the debug block tree
func GetDebug(w *Whiskey) (b *DebugBlock) {
	return getDebugBlock(w, w.l.root)
}

func getDebugBlock(w *Whiskey, index int64) *DebugBlock {
	blk := w.getBlock(index)
	if blk == nil {
		return nil
	}

	var b DebugBlock
	b.Key = string(w.getKey(blk))
	b.ChildType = blk.ct
	b.Color = blk.c
	b.Children[0] = getDebugBlock(w, blk.children[0])
	b.Children[1] = getDebugBlock(w, blk.children[1])
	if parent := w.getBlock(blk.parent); parent != nil {
		b.Parent = string(w.getKey(parent))
	}

	return &b
}
