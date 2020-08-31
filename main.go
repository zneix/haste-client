package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	version = "1.0"
)

var (
	printVersion = flag.Bool("v", false, fmt.Sprintf("Shows program version"))
	hasteURL     = flag.String("d", "https://haste.zneix.eu", fmt.Sprintf("Hastebin server's URL to which stdin should be uploaded."))

	apiRoute   = "/documents"
	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

func main() {
	//handle CLI arguments
	flag.Parse()
	if *printVersion {
		fmt.Println(fmt.Sprintf("Haste Client %s", version))
		return
	}

	//figure out prefix from hasteURL
	if *hasteURL != "" {
		u, err := url.Parse(*hasteURL)
		if err != nil {
			log.Fatal(err)
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			log.Fatal("Scheme must be included in haste url")
		}
	}

	//dirty hack from https://stackoverflow.com/a/41999124
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeNamedPipe) == 0 {
		log.Fatal("Haste client needs to be called from the pipe!")
		return
	}

	stdinBuffer, _ := ioutil.ReadAll(os.Stdin)
	content := string(stdinBuffer)
	uploadToHaste(*hasteURL, content)
}

func uploadToHaste(url string, data string) {

	type HasteResponseData struct {
		Key string `json:"key,omitempty"`
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", *hasteURL, apiRoute), bytes.NewBuffer([]byte(data)))
	if err != nil {
		log.Fatal(err)
		return
	}
	req.Header.Set("User-Agent", fmt.Sprintf("haste-client/%s", version))

	//send the request to hastebin server
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusMultipleChoices {
		log.Fatal(fmt.Sprintf("Error while uploading data to haste server: %d", resp.StatusCode))
		return
	}

	//error out if the invite isn't found or something else went wrong with the request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error while reading resp from haste server: %s", err.Error()))
		return
	}

	var jsonResponse HasteResponseData
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		log.Fatal(fmt.Sprintf("Error while unmarshaling JSON resp from haste server: %s", err.Error()))
		return
	}

	fmt.Println(fmt.Sprintf("%s/%s", url, jsonResponse.Key))
	return
}
