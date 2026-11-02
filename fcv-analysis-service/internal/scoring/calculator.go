package scoring

import "math"

type Components struct{ Code, Collaboration, Breadth, Consistency float64 }
type Result struct {
	Score      float64    `json:"score"`
	Confidence float64    `json:"confidence"`
	Components Components `json:"components"`
	Version    string     `json:"version"`
}

func Calculate(c Components, confidence float64) Result {
	clamp := func(v float64) float64 {
		if v < 0 {
			return 0
		}
		if v > 100 {
			return 100
		}
		return v
	}
	c.Code, c.Collaboration, c.Breadth, c.Consistency = clamp(c.Code), clamp(c.Collaboration), clamp(c.Breadth), clamp(c.Consistency)
	if confidence < 0 {
		confidence = 0
	}
	if confidence > 1 {
		confidence = 1
	}
	score := c.Code*WeightCode + c.Collaboration*WeightCollaboration + c.Breadth*WeightBreadth + c.Consistency*WeightConsistency
	return Result{Score: math.Round(score*100) / 100, Confidence: math.Round(confidence*100) / 100, Components: c, Version: Version}
}
