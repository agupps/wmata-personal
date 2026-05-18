package app

import (
	"net/http"
  "strconv"
	"log/slog"
	"encoding/json"
	"os"
	"github.com/rs/cors"
	"wmata/internal/buses"
	"wmata/internal/metro"
	"wmata/internal/transit"
	"wmata/internal/config"
)

type App struct {
	config *config.Config
}

func New() *App {
	c := &config.Config{}
	if err := c.Parse(); err != nil {
		panic(err)
	}
	return &App{
		config: c,
	}
}

func (a *App) Run() int {  
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("config load success")
	// Register a handler for the root path "/"
	http.HandleFunc("/busStops", a.busStopsHandler)
  http.HandleFunc("/metroStops", a.metroStopsHandler)
	// Start the server on port 8080
	
	logger.Info("server starting on :8080...")
	
	handler := cors.Default().Handler(http.DefaultServeMux)
	http.ListenAndServe(":8080", handler)

	return 0
}

func (a *App) metroStopsHandler(w http.ResponseWriter, req *http.Request) {
	var metroStations []string
	if req.URL.Query().Has("stop") {
		metroStations = []string{req.URL.Query().Get("stop")}
	} else {	
		metroStations = []string{
			"N04", // Spring Hill
			"C03", // Farragut West
		}
	}
	var lines []string
	if req.URL.Query().Has("line") {
		lines = []string{req.URL.Query().Get("line")}
	} else {
		
		lines = []string{
			"SV",
		}
	}

	metroResponse := []transit.Prediction{}
	for _, metroCode := range metroStations {
		metroResponse = append(metroResponse, metro.MetroStops(metroCode, lines, a.config)...)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(metroResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *App) busStopsHandler(w http.ResponseWriter, req *http.Request) {
	var busStops []int
	if req.URL.Query().Has("stop") {
		stopNum, err := strconv.Atoi(req.URL.Query().Get("stop"))
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(http.StatusBadRequest)
      json.NewEncoder(w).Encode(map[string]string{"error": "Invalid input"})
			return
		}
		busStops = []int{stopNum}
	} else {	
		busStops = []int{
			1001694, // home 
			1001212, // farragut sq 
		}
	}
	var lines []string
	if req.URL.Query().Has("line") {
		lines = []string{req.URL.Query().Get("line")}
	} else {
		lines = []string{
			"D72",
		}
	}

	busResponse := []transit.Prediction{}
	for _, busStopNum := range busStops {
		busResponse = append(busResponse, buses.BusStops(busStopNum, a.config, lines)...)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(busResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

