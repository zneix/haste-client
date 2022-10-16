package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	version = "1.1"
)

var (
	printVersion = flag.Bool("v", false, "Shows program version")
	returnRaw    = flag.Bool("r", false, "Returns link to raw content")
	hasteURL     = flag.String("d", "https://haste.zneix.eu", "Hastebin server's URL to which data will be uploaded")

	apiRoute   = "/documents"
	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

func main() {
	// Handle CLI arguments
	flag.Parse()

	if *printVersion {
		fmt.Printf("Haste Client %s\n", version)
		return
	}

	if len(os.Args) == 1 {
		readStdin()
	} else {
		for _, file := range os.Args[1:] {
			if file == "-" {
				readStdin()
			} else {
				data, err := ioutil.ReadFile(file)
				if err != nil {
					log.Fatalf("%s: Failed reading data from file: %s\n", os.Args[0], err)
				}
				uploadToHaste(*hasteURL, string(data))
			}
		}
	}
}

func readStdin() {
	stdinBuffer, _ := ioutil.ReadAll(os.Stdin)
	content := string(stdinBuffer)
	uploadToHaste(*hasteURL, content)
}

func uploadToHaste(url, data string) {
	type HasteResponseData struct {
		Key string `json:"key,omitempty"`
	}

	req, err := http.NewRequest("POST", *hasteURL+apiRoute, bytes.NewBuffer([]byte(data)))
	if err != nil {
		log.Fatal("Error while creating HTTP request:", err)
		return
	}
	req.Header.Set("User-Agent", fmt.Sprintf("haste-client/%s", version))

	// Send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal("Error while performing the request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusMultipleChoices {
		log.Fatalln("Failed to upload data, server responded with", resp.StatusCode)
		return
	}

	// Error out if the invite isn't found or something else went wrong with the request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error while reading response:", err)
		return
	}

	var jsonResponse HasteResponseData
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		log.Fatalln("Error while unmarshaling JSON response:", err)
		return
	}

	var finalURL = url
	if *returnRaw {
		finalURL += "/raw"
	}
	finalURL += "/" + jsonResponse.Key

	fmt.Println(finalURL)
	return
}
