package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

func main() {
	location := "Lancaster,UK"
	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENWEATHERMAP_API_KEY environment variable not set")
		os.Exit(1)
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", location, apiKey)

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

	fmt.Println("       .--.      ")
	fmt.Println("    .-(    ).    ")
	fmt.Println("   (___.__)__)   ")
	fmt.Println("                 ")
	fmt.Printf("Location: %s\n", location)
	fmt.Printf("Weather: %s\n", weatherData.Weather[0].Description)
	fmt.Printf("Temperature: %.2f°C\n", tempCelsius)
	fmt.Printf("Sunrise: %s\n", sunriseTime)
	fmt.Printf("Sunset: %s\n", sunsetTime)
}
