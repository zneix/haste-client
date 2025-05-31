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

const version = "1.2"

var (
	uploadAPIRoute = "/documents"
	httpClient     = &http.Client{
		Timeout: 10 * time.Second,
	}
)

// readStdin wrapper for reading data from standard input and then uploading it
func readStdin(hasteURL string) {
	stdinBuffer, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("Error while processing stdin:", err)
	}
	uploadToHaste(hasteURL, string(stdinBuffer))
}

func uploadToHaste(uploadURL, data string) {
	// Parse the request URL and append the API route for uploading text
	destination, err := url.ParseRequestURI(uploadURL)
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

var returnRaw = flag.Bool("r", false, "Makes returned link point to raw content")

func main() {
	// Handle CLI arguments
	hasteURL := flag.String("d", "https://haste.zneix.eu", "Hastebin server's URL to which data will be uploaded")
	printVersion := flag.Bool("v", false, "Shows program version")
	flag.Parse()

	// Print version and quit
	if *printVersion {
		fmt.Printf("Haste Client %s\n", version)
		return
	}

	// Upload from stdin if there's no file name provided
	if len(flag.Args()) < 1 {
		readStdin(*hasteURL)
		return
	}

	// Otherwise, if arguments are provided use them as file names to upload
	for _, file := range flag.Args() {
		if file == "-" {
			readStdin(*hasteURL)
			continue
		}

		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("%s: Failed reading data from file: %s\n", os.Args[0], err)
		}
		uploadToHaste(*hasteURL, string(data))
	}
}
