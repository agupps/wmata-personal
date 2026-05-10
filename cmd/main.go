package main

import (
	"net/http"
	"fmt"
	"log"
	"github.com/spf13/viper"
	"wmata/internal/buses"
	"wmata/internal/metro"
)

func initConfig() error {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")

    viper.AddConfigPath(".")
		viper.AddConfigPath("..")          

    // Read the config file
    if err := viper.ReadInConfig(); err != nil {
        // Handle the case where config file is not found
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            log.Println("No config file found, using defaults")
            return nil
        }
        return fmt.Errorf("error reading config file: %w", err)
    }

    log.Printf("Using config file: %s", viper.ConfigFileUsed())
    return nil
}

func busStopsHandler(w http.ResponseWriter, req *http.Request) {


	log.Printf("recieving bus request")
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

func metroStopsHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("recieving metro request")
	metroStations := map[string]string{
		"N04": "Spring Hill",
		"C03": "Farragut West",
	}

	printVal := ""
	for metroCode, metroStation := range metroStations {
		printVal += fmt.Sprintf("%s\n", metroStation)
		printVal += metro.MetroStops(metroCode)
		printVal += "\n"
	}
	fmt.Fprintf(w, printVal)
}

func main() {

	if err := initConfig(); err != nil {
  	log.Fatalf("Failed to initialize config: %v", err)
	}
	// Register a handler for the root path "/"
	http.HandleFunc("/busStops", busStopsHandler)
  http.HandleFunc("/metroStops", metroStopsHandler)
	// Start the server on port 8080
	fmt.Println("Server starting on :8080...")
	http.ListenAndServe(":8080", nil)
}
