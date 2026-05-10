package transit 


type Prediction interface {
	GetMinutes() string	
	GetDirectionText() string
	GetRouteID() string
	GetPrediction() string
}

func PredictionsToString(predictions []Prediction) string {
	out := ""
	for _, p := range predictions {
		out += p.GetPrediction()
	}
	return out
}

