package rbt

func isBlack(b *Block) bool {
	if b == nil {
		return true
	}

	return b.c == colorBlack
}

func isRed(b *Block) bool {
	if b == nil {
		return false
	}

	return b.c == colorRed
}
