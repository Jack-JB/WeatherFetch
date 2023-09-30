# Weather App
This is a simple command-line weather app that uses the OpenWeatherMap API to display the current weather in a given location.

## Prerequisites
Before you can run this app, you'll need to sign up for an API key from OpenWeatherMap. You can do this by following these steps:

1. Go to the OpenWeatherMap website and click "Sign Up" in the top right corner.

2. Follow the prompts to create an account and sign in.

3. Once you're signed in, go to the API keys page and click "Generate" to create a new API key.

4. Copy the API key to your clipboard.

## Installation

To install this app, follow these steps:

1. Clone this repository to your local machine.

2. Install the required dependencies by running go get in the project directory:
> go get github.com/joho/godotenv

3. Create a .env file in the project directory and add your OpenWeatherMap API key to it:
> OPENWEATHERMAP_API_KEY=your_api_key_here

Replace your_api_key_here with the API key you generated earlier.

## Usage
To use this app, run the following command in the project directory:
> go run main.go

This will display the current weather in London by default. To get the weather for a different location, modify the location variable in main.go to the desired location.

## License
This project is licensed under the terms of the MIT license. See the [LICENSE](LICENSE) file for details.