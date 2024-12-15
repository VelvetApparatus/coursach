package config

import (
	"encoding/json"
	"log"
	"os"
)

var c *Config

type Config struct {
	ExperimentsCount int `json:"experiments_count"`

	InitialDriverACount int `json:"initial_driver_a_count"`
	InitialDriverBCount int `json:"initial_driver_b_count"`

	InitialBusCount int `json:"initial_bus_count"`

	InitialBusStationsCount int `json:"initial_bus_stations_count"`
	DistinctPathCount       int `json:"distinct_path_count"`
	TimeSeriesPathsCount    int `json:"time_series_paths_count"`
}

func C() *Config { return c }

func init() {
	c = new(Config)
	f, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	err = decoder.Decode(c)
	if err != nil {
		log.Fatal(err)
	}
}
