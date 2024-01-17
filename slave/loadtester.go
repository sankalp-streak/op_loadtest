package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type statusmap struct {
	mu            sync.Mutex
	statusCodeMap map[string]int32
}

var stmp = statusmap{statusCodeMap: map[string]int32{}}

func statusCode(output string) int {
	// Find and extract the HTTP status code from the response headers
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "HTTP/") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				statusCode, err := strconv.Atoi(parts[1])
				if err != nil {
					statusCode = -1
				}
				return statusCode
				break
			}
		}
	}
	return 0
}

func required(s string) bool {

	nots := []string{
		"track.streak.ninja",
		"os-analytics.streak.tech",
		"/static/js",
		"open-v2.streak.ninja/_next/static",
		"curl 'https://open-v2.streak.ninja/home'",
		"https://www.googletagmanager.com/'",
		"fonts.googleapis.com",
		"facebook",
		"os-analytics",
		"www.google.co.in",
		"www.googletagmanager.com",
		"data:image/gif",
		"streak.ninja/logo3.ico",
		"/manifest.json",
		"data:image/",
		"_next/static",
		"blog.streak.tech",
		"treak-public-assets",
		"chrome-extension",
		"streak_192.png",
		"analytics.google",
		"streak.ninja/home/strategies",
		"refapi.streak.tech",
		"wss://ss-op.st",
		"nt-op.streak.tech",
	}

	for _, n := range nots {
		if strings.Contains(s, n) {
			return false
		}
	}
	return true

}

func removeEmptyStrings(slice []string) []string {
	var result []string
	for _, str := range slice {
		trimmedStr := strings.TrimSpace(str)
		if trimmedStr != "" {
			result = append(result, trimmedStr)
		}
	}
	return result
}

func cleanCurls(curls []string) []string {
	var filteredCurls []string
	for _, c := range curls {
		if required(c) {
			c = strings.ReplaceAll(c, "\r", "")
			filteredCurls = append(filteredCurls, strings.TrimSpace(c))
		}
	}
	filteredCurls = removeEmptyStrings(filteredCurls)
	return filteredCurls

}

func getCurls(s string) []string {
	// Read the content of the file
	fileContent, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return []string{}
	}

	fileContent = []byte(strings.ReplaceAll(string(fileContent), "--compressed ;", "--compressed ; @@Sankalp@@"))
	curls := strings.Split(string(fileContent), " @@Sankalp@@")
	curls = cleanCurls(curls)
	return curls

}

var mp = map[string]int{}

func hitTarget(curlRequest string) ress {
	var r ress
	cmd := exec.Command("sh", "-c", curlRequest) // Create a new Cmd struct to run the cURL command
	op, err := cmd.Output()                      // Execute the cURL command
	if err != nil {
		fmt.Println("Error executing cURL:", err)
		return ress{}
	}
	st := statusCode(string(op))
	r.mu.Lock()
	if st >= 200 && st <= 299 {
		r.suc++
		r.Total++
	} else if st > 299 {
		r.fail++
		r.Total++
	}
	r.mu.Unlock()
	return ress{}

}
