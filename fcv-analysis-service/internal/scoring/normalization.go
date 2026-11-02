package scoring

import "math"

// Normalize maps a non-negative value to 0..100 using log scaling so one prolific
// contributor cannot dominate solely through raw volume.
func Normalize(value, maxValue float64) float64 {
	if value <= 0 || maxValue <= 0 {
		return 0
	}
	v := 100 * math.Log1p(value) / math.Log1p(maxValue)
	if v > 100 {
		return 100
	}
	return math.Round(v*100) / 100
}
