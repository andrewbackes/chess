package position

func popcount(b uint64) uint {
	var count uint
	for i := uint(0); i < 64; i++ {
		if (b & (1 << i)) != 0 {
			count++
		}
	}
	return count
}

func bitscan(b uint64) uint {
	for i := uint(0); i < 64; i++ {
		if (b & (1 << i)) != 0 {
			return i
		}
	}
	return 64
}

func bsf(b uint64) uint {
	for i := uint(0); i < 64; i++ {
		if (b & (1 << i)) != 0 {
			return i
		}
	}
	return 64
}

func bsr(b uint64) uint {
	for i := uint(63); i > 0; i-- {
		if (b & (1 << i)) != 0 {
			return i
		}
	}
	if b&1 != 0 {
		return 0
	}
	return 64
}
