package kmgMath

import "math"

func FloorToInt(x float64) int {
	return int(math.Floor(x))
}

func CeilToInt(x float64) int {
	return int(math.Ceil(x))
}
