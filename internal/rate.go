package internal

// ComputeRate computes rate based on the given input.
func ComputeRate(actual uint64, last uint64, interval float64) float64 {
	diff := int64(actual - last)
	if diff > 0 {
		return float64(diff) / interval
	} else {
		return float64(-diff) / interval
	}
}

// RateComputer is clojure around ComputeRate.
func RateComputer(interval float64) func(uint64, uint64) float64 {
	return func(actual, last uint64) float64 {
		return ComputeRate(actual, last, interval)
	}
}
