package diskv

func BlockTransform(blockSize, maxSlice int, suffix bool) func(string) []string {
	return func(s string) []string {
		sliceSize := func() int {
			if sliceSize := len(s) / blockSize; sliceSize <= maxSlice {
				return sliceSize
			}
			return maxSlice
		}()

		if suffix {
			s = s[len(s)-sliceSize*blockSize:]
		}

		pathSlice := make([]string, sliceSize)
		for i := 0; i < sliceSize; i++ {
			from, to := i*blockSize, (i*blockSize)+blockSize
			pathSlice[i] = s[from:to]
		}

		return pathSlice
	}
}
