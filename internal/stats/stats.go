package stats

// Min returns the minimum value in a float64 slice.
// The slice must not be empty.
func Min(array []float64) float64 {
	min := array[0]
	for _, value := range array {
		if value < min {
			min = value
		}
	}
	return min
}

// Max returns the maximum value in a float64 slice.
// The slice must not be empty.
func Max(array []float64) float64 {
	max := array[0]
	for _, value := range array {
		if value > max {
			max = value
		}
	}
	return max
}

// Avg returns the arithmetic mean of a float64 slice.
// The slice must not be empty.
func Avg(array []float64) float64 {
	sum := 0.0
	for _, value := range array {
		sum += value
	}
	return sum / float64(len(array))
}
