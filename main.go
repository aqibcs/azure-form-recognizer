package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Define global variables to store configuration and document information
var (
	endpoint        string
	modelID         string
	apiVersion      string
	subscriptionKey string
	documentURL     string
)

// The main function where the execution begins
func main() {
	// Load environment variables from the .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// Retrieve configuration values from environment variables
	endpoint = os.Getenv("ENDPOINT")
	modelID = os.Getenv("MODEL_ID")
	apiVersion = os.Getenv("API_VERSION")
	subscriptionKey = os.Getenv("SUBSCRIPTION_KEY")
	documentURL = os.Getenv("DOCUMENT_URL")

	fmt.Println("Sending POST request to analyze document...")

	// Perform document analysis and obtain the result ID
	resultID, err := analyzeDocument()
	if err != nil {
		fmt.Println("Error analyzing document:", err)
		return
	}

	// Display the result ID obtained from the analysis request
	fmt.Println("Result ID:", resultID)

	// Poll for the analysis result status until it succeeds
	fmt.Println("Checking for results...")
	var status string
	var detailedResults map[string]interface{}

	for status != "succeeded" {
		// Retrieve detailed results and status from the analysis result
		detailedResults, err = getAnalyzeResult(resultID)
		if err != nil {
			fmt.Println("Error getting analyze result:", err)
			return
		}

		// Extract the status from the detailed results
		status, _ = detailedResults["status"].(string)

		// Display the current status
		fmt.Println("Status:", status)

		// If the status is not yet succeeded, wait for 5 seconds before polling again
		if status != "succeeded" {
			time.Sleep(5 * time.Second)
		}
	}

	// Display a success message and output detailed results in JSON format
	fmt.Println("Analysis completed successfully. Detailed Results:")

	// Convert detailed results to JSON
	jsonOutput, err := json.MarshalIndent(detailedResults, "", "  ")
	if err != nil {
		fmt.Println("Error encoding detailed results to JSON:", err)
		return
	}

	// Output the JSON
	fmt.Println(string(jsonOutput))
}

// Function to initiate document analysis and obtain the result ID
func analyzeDocument() (string, error) {
	// Construct the URL for document analysis based on the endpoint, model ID, and API version
	url := fmt.Sprintf("%s/formrecognizer/documentModels/%s:analyze?api-version=%s", endpoint, modelID, apiVersion)

	// Prepare the payload data for the analysis request (in this case, specifying the document URL)
	data := map[string]interface{}{
		"urlSource": documentURL,
	}

	// Convert payload data to JSON format
	payload, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// Create a new HTTP POST request with the constructed URL and payload
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	// Set headers for the request, including content type and subscription key
	req.Header.Set("Content-Type", "image/jpeg")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", subscriptionKey)

	// Send the request and handle the response
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the response status code indicates success (StatusAccepted)
	if resp.StatusCode != http.StatusAccepted {
		responseBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code: %d, response body: %s", resp.StatusCode, responseBody)
	}

	// Extract the operation location (result ID) from the response headers
	resultID := resp.Header.Get("Operation-Location")
	return resultID, nil
}

// Function to retrieve the detailed results of the document analysis
func getAnalyzeResult(resultID string) (map[string]interface{}, error) {
	// Construct the URL for obtaining the analysis result based on the result ID
	url := resultID

	// Create a new HTTP GET request with the constructed URL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers for the request, including the subscription key
	req.Header.Set("Ocp-Apim-Subscription-Key", subscriptionKey)

	// Send the request and handle the response
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode the JSON response body into a map
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	// Return the obtained detailed results
	return result, nil
}
