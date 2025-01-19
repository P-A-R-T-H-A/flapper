package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	swarmgo "github.com/prathyushnallamothu/swarmgo"
	"github.com/prathyushnallamothu/swarmgo/llm"
)

// WebsiteSource represents a travel trend data source
type WebsiteSource struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

// WeatherRequest represents the parameters for the getWeather function
type WeatherRequest struct {
	Location string `json:"location"`
}

// fetchWebsiteContent retrieves content from a given URL
func fetchWebsiteContent(url string) (string, error) {
	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Send GET request
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching URL %s: %v", url, err)
	}
	defer resp.Body.Close()

	// Check HTTP response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(body), nil
}

func main() {
	// Initialize Gemini client with API key
	apiKey := "AIzaSyCd5duoB6JPX2Gz5CqbmeIcDefvXGkpVLM"
	if apiKey == "" {
		log.Fatal("Google API key is required")
	}

	swarm := swarmgo.NewSwarm(apiKey, llm.Gemini)
	ctx := context.Background()

	// Define travel trend data sources
	sources := []WebsiteSource{
		// {
		// 	URL:   "https://www.forbes.com/travel/trends/",
		// 	Title: "Forbes Travel Trends",
		// },
		{
			URL:   "https://www.forbes.com/sites/forbesdaily/",
			Title: "Fobes",
		},
		// {
		// 	URL:   "https://destinationinsights.withgoogle.com/intl/en_ALL/",
		// 	Title: "Google Destination Insights",
		// },
	}

	// Collect and process data from multiple sources
	var allTrendData []string
	for _, source := range sources {
		content, err := fetchWebsiteContent(source.URL)
		if err != nil {
			log.Printf("Error fetching %s: %v", source.Title, err)
			continue
		}

		// Create an agent to analyze the web content
		analyzeAgent := &swarmgo.Agent{
			Name: "WebTrendAnalyzer",
			// Instructions: "Analyze the provided web content and extract key travel trends. set Country-specific travel demand section data Origin country = worldWide, Destination country= worldwide, Trip type = internation,Category= accommodation. Format the output as a JSON object with categories: TOP DEMAND BY ORIGIN COUNTRY",
			Instructions: "Analyze the provided web content and extract key of headings and make a json output by each sections and then concatenate all the sections and make a single json output",
			Model:        "gemini-1.5-flash",
		}

		// Prepare messages for analysis
		analyzeMessages := []llm.Message{
			{
				Role:    llm.RoleUser,
				Content: fmt.Sprintf("Analyze travel trends from %s. Web content: %s", source.Title, content),
			},
		}

		// Run analysis
		analysisResponse, err := swarm.Run(ctx, analyzeAgent, analyzeMessages, nil, "", true, false, 5, true)
		if err != nil {
			log.Printf("Error analyzing %s: %v", source.Title, err)
			continue
		}

		// Store the JSON trend data
		allTrendData = append(allTrendData, analysisResponse.Messages[len(analysisResponse.Messages)-1].Content)
	}
	fmt.Println("All trend data:", allTrendData)
	// // Combine and validate trend data
	// finalJSON := "["
	// for _, trendJSON := range allTrendData {
	// 	// Validate each trend JSON
	// 	var trendData json.RawMessage
	// 	err := json.Unmarshal([]byte(trendJSON), &trendData)
	// 	if err != nil {
	// 		log.Printf("Error parsing trend JSON: %v", err)
	// 		continue
	// 	}

	// 	// Add to combined JSON
	// 	finalJSON += trendJSON + ","
	// }

	// // Remove trailing comma and close JSON array
	// if len(finalJSON) > 1 {
	// 	finalJSON = finalJSON[:len(finalJSON)-1]
	// }
	// finalJSON += "]"

	// // Pretty print the final JSON
	// var prettyJSON bytes.Buffer
	// err := json.Indent(&prettyJSON, []byte(finalJSON), "", "  ")
	// if err != nil {
	// 	log.Printf("Error formatting JSON: %v", err)
	// 	return
	// }

	// fmt.Println(prettyJSON.String())
}