package backend

func next32(val int64) int64 {
	rem := val % 32
	if rem == 0 {
		return val
	}

	return val + (32 - rem)
}

func nextCap(cap, sz int64) int64 {
	if sz <= cap {
		return -1
	} else if cap == 0 {
		return next32(sz)
	}

	for cap < sz {
		cap *= 2
	}

	return cap
}
