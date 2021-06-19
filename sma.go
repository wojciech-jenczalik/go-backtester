package main

func Sma(numbers ...float64) float64 {
	var sum float64
	for _, price := range numbers {
		sum += price
	}

	return sum / float64(len(numbers))
}
