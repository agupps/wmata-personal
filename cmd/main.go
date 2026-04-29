package main

import (
	"net/http"
	"fmt"
	"io"
	"encoding/json"
)

func main() {
	stopNumbers := []int{1001694, 1001212}
	url := "https://api.wmata.com/NextBusService.svc/json/jPredictions?StopID=%d"

	for _, stopNumber := range stopNumbers {
	request, _ := http.NewRequest("GET", fmt.Sprintf(url, stopNumber), nil)
	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("api_key", "a85258d235654c9abc4a37e47f7b0c20")
	
	client := &http.Client{}
	response, _ := client.Do(request)

	defer response.Body.Close()

	// Read and print the response body
	body, _ := io.ReadAll(response.Body)

	var nbs NextBusServiceResponse
	err := json.Unmarshal(body, &nbs)
	if err != nil {
		fmt.Println("%v",err)
		return
	}

	for _, prediction := range nbs.Predictions {
		if prediction.RouteID == "D72" {
			fmt.Println(prediction)
		}
	}
}
}

type NextBusServiceResponse struct {
	StopName string `json:"StopName"`
	Predictions []Prediction `json:"Predictions"`
}

type Prediction struct {
	RouteID string `json:"RouteID"`
	DirectionText string `json:"DirectionText"`
	Minutes int `json:"Minutes"`
}


