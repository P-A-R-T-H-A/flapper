package controllers

import (
	"context"
	"encoding/json"
	agents "flapper/Agents"
	"fmt"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/prathyushnallamothu/swarmgo/llm"
)

type TestController struct {
	beego.Controller
}

func (c *TestController) Get() {
	// Get the query string parameter
	location := c.GetString("voice")
	if location == "" {
		c.Data["json"] = map[string]string{"error": "Missing required query parameter: voice"}
		c.ServeJSON()
		return
	}

	// Create a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Step 1: Configure and execute weather fetcher agent
	weatherAgent := &agents.Agents{
		Name:         "weather-fetcher",
		Model:        "gemini-1.5-flash",
		Instructions: fmt.Sprintf("Return ONLY the raw JSON response from this API: http://api.weatherapi.com/v1/current.json?key=0227c8e628654700b92100850251901&q=%s&aqi=no", location),
		Provider:     llm.Gemini,
		Stream:       false,
		Debug:        true,
		MaxTurns:     5,
		ExecuteTools: true,
	}
	weatherAgent.LoadAgent()

	weatherMessages := []llm.Message{
		{
			Role:    "system",
			Content: "You are a weather data fetching agent. Return only raw JSON responses without any formatting or explanation.",
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("Get weather data for %s", location),
		},
	}

	weatherContext := map[string]interface{}{
		"location": location,
	}

	weatherResponse, err := weatherAgent.Execute(ctx, weatherMessages, weatherContext)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Weather agent execution failed", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Extract weather data from response
	var weatherResult string
	for _, msg := range weatherResponse.Messages {
		if msg.Role == "assistant" {
			weatherResult = msg.Content
			break
		}
	}

	// Clean up the weather result - remove any backticks or code block markers
	weatherResult = strings.TrimSpace(weatherResult)
	weatherResult = strings.TrimPrefix(weatherResult, "```json")
	weatherResult = strings.TrimPrefix(weatherResult, "```")
	weatherResult = strings.TrimSuffix(weatherResult, "```")
	weatherResult = strings.TrimSpace(weatherResult)

	if weatherResult == "" {
		c.Data["json"] = map[string]string{"error": "No weather data returned"}
		c.ServeJSON()
		return
	}

	// Validate weather JSON before proceeding
	var weatherJSON map[string]interface{}
	if err := json.Unmarshal([]byte(weatherResult), &weatherJSON); err != nil {
		c.Data["json"] = map[string]string{"error": "Invalid weather data JSON", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Step 2: Configure and execute JSON processor agent
	processingAgent := &agents.Agents{
		Name:         "json-processor",
		Model:        "gemini-1.5-flash",
		Instructions: "Extract weather information and return it as a JSON object. Do not include any explanation or formatting. Return only the JSON object.",
		Provider:     llm.Gemini,
		Stream:       false,
		Debug:        true,
		MaxTurns:     5,
		ExecuteTools: false,
	}
	processingAgent.LoadAgent()

	processingMessages := []llm.Message{
		{
			Role:    "system",
			Content: "You are a JSON processing agent. Format the data exactly like this, with no additional text or formatting: {\"location\":{\"name\":\"\",\"region\":\"\",\"country\":\"\"},\"current\":{\"temp_c\":0,\"temp_f\":0,\"condition\":{\"text\":\"\",\"icon\":\"\"},\"wind_kph\":0,\"wind_mph\":0,\"humidity\":0,\"cloud\":0,\"last_updated\":\"\"}}",
		},
		{
			Role:    "user",
			Content: weatherResult,
		},
	}

	processingContext := map[string]interface{}{}

	processedResponse, err := processingAgent.Execute(ctx, processingMessages, processingContext)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Processing agent execution failed", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Extract processed data from response
	var processedResult string
	for _, msg := range processedResponse.Messages {
		if msg.Role == "assistant" {
			processedResult = msg.Content
			break
		}
	}

	// Clean up the processed result
	processedResult = strings.TrimSpace(processedResult)
	processedResult = strings.TrimPrefix(processedResult, "```json")
	processedResult = strings.TrimPrefix(processedResult, "```")
	processedResult = strings.TrimSuffix(processedResult, "```")
	processedResult = strings.TrimSpace(processedResult)

	if processedResult == "" {
		c.Data["json"] = map[string]string{"error": "No processed data returned"}
		c.ServeJSON()
		return
	}

	// Parse the processed JSON
	var finalWeatherJSON map[string]interface{}
	err = json.Unmarshal([]byte(processedResult), &finalWeatherJSON)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to parse processed JSON", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Step 3: Respond with the processed JSON
	c.Data["json"] = map[string]interface{}{
		"voice":    location,
		"response": finalWeatherJSON,
	}
	c.ServeJSON()
}
