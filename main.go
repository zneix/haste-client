package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	version = "1.2"
)

var (
	printVersion = flag.Bool("v", false, "Shows program version")
	returnRaw    = flag.Bool("r", false, "Makes returned link point to raw content")
	hasteURL     = flag.String("d", "https://haste.zneix.eu", "Hastebin server's URL to which data will be uploaded")

	uploadAPIRoute = "/documents"
	httpClient     = &http.Client{
		Timeout: 10 * time.Second,
	}
)

func readStdin() {
	stdinBuffer, _ := io.ReadAll(os.Stdin)
	content := string(stdinBuffer)
	uploadToHaste(content)
}

func uploadToHaste(data string) {
	type HasteResponseData struct {
		Key string `json:"key,omitempty"`
	}

	// Parse the request URL and append the API route for uploading text
	destination, err := url.ParseRequestURI(*hasteURL)
	if err != nil {
		log.Fatal("Error while parsing destination URL (full URL with a protocol scheme is required):", err)
		return
	}
	uploadEndpoint := destination.JoinPath(uploadAPIRoute)

	// Create & Send the request
	req, err := http.NewRequest("POST", uploadEndpoint.String(), bytes.NewBufferString(data))
	if err != nil {
		log.Fatal("Error while creating HTTP request:", err)
		return
	}
	req.Header.Set("User-Agent", fmt.Sprintf("haste-client/%s", version))

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal("Error while performing the request:", err)
		return
	}
	defer resp.Body.Close()

	// Error out if request wasn't handled correctly by the remote server (which is indicated by the response status)
	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusMultipleChoices {
		log.Fatalln("Failed to upload data, server responded with", resp.StatusCode)
		return
	}

	// Read & process the received response data
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error while reading response:", err)
		return
	}

	jsonResponse := new(HasteResponseData)
	if err := json.Unmarshal(body, jsonResponse); err != nil {
		log.Fatalln("Error while unmarshaling JSON response:", err)
		return
	}

	// Handle returning raw text if it is desired so and respond with uploaded text's URL
	if *returnRaw {
		destination = destination.JoinPath("/raw")
	}
	fmt.Println(destination.JoinPath(jsonResponse.Key).String())
}

func main() {
	// Handle CLI arguments
	flag.Parse()

	// Print version and quit
	if *printVersion {
		fmt.Printf("Haste Client %s\n", version)
		return
	}

	// Upload from stdin if there's no file name provided
	if len(flag.Args()) < 1 {
		readStdin()
		return
	}

	// Otherwise, if arguments are provided use them as file names to upload
	for _, file := range flag.Args() {
		if file == "-" {
			readStdin()
			continue
		}

		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("%s: Failed reading data from file: %s\n", os.Args[0], err)
		}
		uploadToHaste(string(data))
	}
}
