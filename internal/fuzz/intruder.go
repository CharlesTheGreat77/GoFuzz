package fuzz

import (
	"fmt"
	"gofuzz/args"
	"gofuzz/internal/goclient"
	"gofuzz/utils"
	"io"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"
)

// function to fuzz the parameters in a given wordlist based on position in url or body
func GoIntruder(fuzzy args.Fuzzy, wordlist []string, headers []string, proxies []string) {
	timeout := time.Duration(fuzzy.Timeout) * time.Second
	client := goclient.GoClient(timeout, goclient.NewProxyRotator(proxies))

	semaphore := make(chan struct{}, fuzzy.Threads)
	var wg sync.WaitGroup

	for _, word := range wordlist {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(word string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			// replace FUZZ for each line
			encodedWord := url.QueryEscape(word) // URL encode
			modifiedURL := strings.Replace(fuzzy.URL, "FUZZ", encodedWord, -1)
			modifiedBody := strings.Replace(fuzzy.Body, "FUZZ", word, -1)

			req, err := goclient.GoRequest(fuzzy, modifiedURL, headers, modifiedBody, word)

			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error sending request to %s: %v\n", fuzzy.URL, err)
				return
			}

			defer func() {
				if err := resp.Body.Close(); err != nil {
					log.Printf("Error closing response body for %s: %v\n", fuzzy.URL, err)
				}
			}()

			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading response from %s: %v\n", fuzzy.URL, err)
				return
			}

			parsedURL, err := url.Parse(modifiedURL)
			if err != nil {
				log.Printf("Error parsing URL %s: %v\n", modifiedURL, err)
				return
			}
			pathAndQuery := parsedURL.Path + parsedURL.RawQuery

			if len(fuzzy.StatusCodes) != 0 {
				sc := fmt.Sprintf("%d", resp.StatusCode)
				for _, code := range fuzzy.StatusCodes {
					if string(code) == sc {
						utils.RequestOutput(pathAndQuery, fuzzy.SearchString, sc, string(responseBody), modifiedBody, fuzzy.NoSearch)
					}
				}
			} else if resp.StatusCode != 404 { // Ignore 404 responses
				utils.RequestOutput(pathAndQuery, fuzzy.SearchString, fmt.Sprintf("%d", resp.StatusCode), string(responseBody), modifiedBody, fuzzy.NoSearch)
			}
		}(word)
	}
	wg.Wait()
}
