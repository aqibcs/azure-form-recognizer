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

var (
	endpoint        string
	modelID         string
	apiVersion      string
	subscriptionKey string
	documentURL     string
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// Load environment variables from .env file
	endpoint = os.Getenv("ENDPOINT")
	modelID = os.Getenv("MODEL_ID")
	apiVersion = os.Getenv("API_VERSION")
	subscriptionKey = os.Getenv("SUBSCRIPTION_KEY")
	documentURL = os.Getenv("DOCUMENT_URL")

	fmt.Println("Sending POST request to analyze document...")
	resultID, err := analyzeDocument()
	if err != nil {
		fmt.Println("Error analyzing document:", err)
		return
	}

	fmt.Println("Result ID:", resultID)

	fmt.Println("Checking for results...")

	var status string
	var detailedResults map[string]interface{}

	for status != "succeeded" {
		detailedResults, err = getAnalyzeResult(resultID)
		if err != nil {
			fmt.Println("Error getting analyze result:", err)
			return
		}

		status, _ = detailedResults["status"].(string)

		fmt.Println("Status:", status)

		if status != "succeeded" {
			time.Sleep(5 * time.Second)
		}
	}

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

func analyzeDocument() (string, error) {
	url := fmt.Sprintf("%s/formrecognizer/documentModels/%s:analyze?api-version=%s", endpoint, modelID, apiVersion)

	data := map[string]interface{}{
		"urlSource": documentURL,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "image/jpeg")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", subscriptionKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		responseBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code: %d, response body: %s", resp.StatusCode, responseBody)
	}

	resultID := resp.Header.Get("Operation-Location")
	return resultID, nil
}

func getAnalyzeResult(resultID string) (map[string]interface{}, error) {
	url := resultID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", subscriptionKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
