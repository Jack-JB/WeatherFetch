// Author: https://github.com/Jack-JB
// Description: A simple command line tool to get the weather for a given location.
// Copyright (c) 2023, Jack Burrows.
// All rights reserved.
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type WeatherData struct {
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Sys struct {
		Sunrise int64 `json:"sunrise"`
		Sunset  int64 `json:"sunset"`
	} `json:"sys"`
}

type Config struct {
	DefaultLocation string `json:"default_location"`
}

func main() {
	location := flag.String("location", "", "Location to get weather information for")
	setDefaultLocation := flag.String("set-default-location", "", "Set the defualt location in config.json")
	flag.Parse()

	configData, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error: Could not read config.json", err)
		os.Exit(1)
	}

	var config Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		fmt.Println("Error: Could not parse config.json", err)
		os.Exit(1)
	}

	// Set the default location in the configuration file if the -set-default-location flag is specified
	if *setDefaultLocation != "" {
		config.DefaultLocation = *setDefaultLocation

		// Marshal the Config struct into JSON data
		configData, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling configuration data:", err)
			os.Exit(1)
		}

		// Write the JSON data to the configuration file
		err = os.WriteFile("config.json", configData, 0644)
		if err != nil {
			fmt.Println("Error writing configuration file:", err)
			os.Exit(1)
		}

		fmt.Printf("Default location set to %s\n", *setDefaultLocation)
		os.Exit(0)
	}

	if *location == "" {
		*location = config.DefaultLocation
	}

	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	if apiKey == "" {
		fmt.Println("ERROR: OPENWEATHERMAP_API_KEY environment variable not set")
		os.Exit(1)
	}

	done := make(chan bool)
	go func() {
		url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", *location, apiKey)

		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		resp, err := client.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var weatherData WeatherData
		err = json.Unmarshal(data, &weatherData)
		if err != nil {
			panic(err)
		}

		tempCelsius := weatherData.Main.Temp - 273.15

		sunriseTime := time.Unix(weatherData.Sys.Sunrise, 0).Format("15:04:05")
		sunsetTime := time.Unix(weatherData.Sys.Sunset, 0).Format("15:04:05")

		asciiSun := `
          \   /
           .-.
        ‒ (   ) ‒
           .-.
          /   \
      `
		asciiCloud := `
            .--.
         .-(    ).
        (___.__)__)
      `
		acsiiRainCloud := `
          .--.
       .-(    ).
      (___.__)__)
       ʻ‚ʻ‚ʻ‚ʻ‚ʻ
      `
		asciiSnow := `
          .--.
       .-(    ).
      (___.__)__)
       * * * * *
      `

		asciiThunder := `
         .--.
      .-(    ).
     (___.__)__)
          /_
           /
      `

		asciiArt := map[string]string{
			"sun":     asciiSun,
			"rain":    acsiiRainCloud,
			"snow":    asciiSnow,
			"thunder": asciiThunder,
		}

		if art, ok := asciiArt[strings.ToLower(weatherData.Weather[0].Description)]; ok {
			fmt.Println(art)
		} else {
			fmt.Println(asciiCloud)
		}

		fmt.Printf("Location: %s\n", *location)
		fmt.Printf("Weather: %s\n", weatherData.Weather[0].Description)
		fmt.Printf("Temperature: %.2f°C\n", tempCelsius)
		fmt.Printf("Sunrise: %s\n", sunriseTime)
		fmt.Printf("Sunset: %s\n", sunsetTime)

		done <- true
	}()

	select {
	case <-done:
		// Goroutine completed successfully
	case <-time.After(15 * time.Second):
		// Goroutine timed out
		fmt.Println("Error: Request timed out")
		os.Exit(1)
	}
}
