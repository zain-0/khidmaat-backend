package utils

var ClassLabels = []string{
	"Normal",
	"Supraventricular prem",
	"Premature ventr.",
	"Fusion of ve. & no.",
	"Unclassifiable",
}

type PredictionRequest struct {
	Segment []float64 `json:"segment"`
}

type PredictionResponse struct {
	Prediction [][]float64 `json:"prediction"` // Change this to an array of arrays
}
