package allocator

// Section represents an allocated section of bytes
type Section struct {
	Pair
	Bytes []byte
}

func (s *Section) getOnGrow(a Allocator) (fn func()) {
	fn = func() {
		if s.IsEmpty() {
			return
		}

		s.Bytes = a.Get(s.Offset, s.Size)
	}

	return
}

func (s *Section) destroy() {
	s.Offset = 0
	s.Size = 0
	s.Bytes = nil
}
