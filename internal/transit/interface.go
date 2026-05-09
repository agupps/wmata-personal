package transit 


type Prediction interface {
	GetMinutes() string
	GetRouteID() string
	GetDirectionText() string
	GetPrediction() string
}

func PredictionsToString(predictions []Prediction) string {
	out := ""
	for _, p := range predictions {
		out += p.GetPrediction()
	}
	return out
}

