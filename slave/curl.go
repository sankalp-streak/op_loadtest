package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	curl "optest/curlgo"
	"strings"
)

func parseCurlRequest(curlString string) (*http.Request, error) {
	// Split the curl command into individual tokens
	tokens := strings.Fields(curlString)

	// Initialize a new http.Request
	req := http.Request{}

	// Set the request method (e.g., GET, POST)
	req.Method = tokens[1]

	// Parse the URL
	parsedURL, err := url.Parse(tokens[0])
	if err != nil {
		return nil, err
	}
	req.URL = parsedURL

	// Parse headers
	for i := 2; i < len(tokens)-1; i += 2 {
		headerName := strings.Trim(tokens[i], " -H'")
		headerValue := strings.Trim(tokens[i+1], "'")
		req.Header.Add(headerName, headerValue)
	}

	// Parse request body
	req.Body = http.NoBody
	for i := 0; i < len(tokens); i++ {
		if tokens[i] == "--data-raw" && i+1 < len(tokens) {
			req.Body = ioutil.NopCloser(strings.NewReader(strings.Join(tokens[i+1:], " ")))
			break
		}
	}

	return &req, nil
}

func RunCurl(curlCommand string) int {
	command, err := curl.Parse(curlCommand)
	if err != nil {
		return -1
	}
	req, err := command.ToRequest()
	if err != nil {
		return -1
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return -1
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Println("DEKHO", resp.StatusCode, bodyString)
		fmt.Println(curlCommand)
	}

	return resp.StatusCode

}