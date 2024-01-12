# Golang Form Recognizer Document Analysis

This Golang application demonstrates document analysis using the Azure Form Recognizer API. It takes an image URL of a document, sends a POST request to the Form Recognizer API for analysis, and retrieves detailed results.

## Prerequisites

Before running the application, make sure you have the following:

- [Golang](https://golang.org/) installed
- Azure Form Recognizer subscription key and endpoint
- Set up a `.env` file with the required environment variables.

## Installation

1. Clone the repository:

    ```bash
    git clone git@github.com:aqibcs/azure-form-recognizer.git
    ```

2. Install dependencies:

    ```bash
    go get -u github.com/joho/godotenv
    ```

3. Create a `.env` file and populate it with your Azure Form Recognizer credentials and document URL.

4. Run the application:

    ```bash
    go run main.go
    ```

## Usage

Follow the on-screen instructions to analyze a document. The application will provide the result ID and continuously check the status until the analysis is complete. The detailed results will be displayed in JSON format.
