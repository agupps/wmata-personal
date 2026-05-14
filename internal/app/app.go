package app

import (
	"net/http"
	"fmt"
	"log/slog"
	"encoding/json"
	"os"
	"github.com/spf13/viper"
	"wmata/internal/buses"
	"wmata/internal/metro"
	"wmata/internal/transit"
)

type App struct {}

func New() *App {
	return &App{}
}

func (a *App) Run() int {  
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := initConfig(); err != nil {
  	logger.Error("Failed to initialize config", "error", err)
		return 1
	}
	logger.Info("config file successful")
	// Register a handler for the root path "/"
	http.HandleFunc("/busStops", busStopsHandler)
  http.HandleFunc("/metroStops", metroStopsHandler)
	// Start the server on port 8080
	logger.Info("server starting on :8080...")
	http.ListenAndServe(":8080", nil)

	return 0
}

func metroStopsHandler(w http.ResponseWriter, req *http.Request) {
	var metroStations []string
	if req.URL.Query().Has("stop") {
		metroStations = []string{req.URL.Query().Get("stop")}
	} else {	
		metroStations = []string{
			"N04", // Spring Hill
			"C03", // Farragut West
		}
	}

	metroResponse := []transit.Prediction{}
	for _, metroCode := range metroStations {
		metroResponse = append(metroResponse, metro.MetroStops(metroCode)...)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(metroResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func busStopsHandler(w http.ResponseWriter, req *http.Request) {
	stopNumbers := map[int]string{
		1001694: "my home", // my house south 
		1001212: "farragut", // farragut to north
  }

	printVal := ""

	for stopNum, stopName:= range stopNumbers {
		printVal += fmt.Sprintf("%s\n", stopName)
		printVal += buses.BusStops(stopNum)
		printVal += "\n"
	}
	
	fmt.Fprintf(w, printVal)
}

func initConfig() error {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")

    viper.AddConfigPath(".")
		viper.AddConfigPath("..")          

    // Read the config file
    if err := viper.ReadInConfig(); err != nil {
        return fmt.Errorf("error reading config file: %w", err)
    }

    return nil
}

