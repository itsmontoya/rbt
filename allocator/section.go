package allocator

// Section represents an allocated section of bytes
type Section struct {
	Offset int64
	Size   int64
}

func (s *Section) destroy() {
	s.Offset = 0
	s.Size = 0
}

// IsEmpty will return if a pair is empty
func (s *Section) IsEmpty() bool {
	if s.Offset != 0 {
		return false
	}

	if s.Size != 0 {
		return false
	}

	return true
}
