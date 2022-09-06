package slice

// Partition creates partitions of a standard maximum size.
func Partition[T any](s []T, partitionSize int) [][]T {
	if len(s) == 0 || partitionSize <= 0 {
		return [][]T{}
	}

	partitions := make([][]T, 0, (len(s)+partitionSize-1)/partitionSize)

	for {
		left := len(partitions) * partitionSize
		if left >= len(s) {
			break
		}

		right := Min(left+partitionSize, len(s))

		part := s[left:right]
		partition := make([]T, len(part))
		copy(partition, part)

		partitions = append(partitions, partition)
	}

	return partitions
}
