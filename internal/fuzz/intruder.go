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

func GoIntruder(fuzzy args.Fuzzy, wordlist []string, headers []string, proxies []string) {
	timeout := time.Duration(fuzzy.Timeout) * time.Second
	session := goclient.NewSession(timeout, goclient.NewProxyRotator(proxies))
	semaphore := make(chan struct{}, fuzzy.Threads)
	var wg sync.WaitGroup

	for _, word := range wordlist {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(word string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			encodedWord := url.QueryEscape(word)
			modifiedURL := strings.Replace(fuzzy.URL, "FUZZ", encodedWord, -1)
			modifiedBody := strings.Replace(strings.TrimSpace(fuzzy.Body), "FUZZ", word, -1)

			req, err := goclient.GoRequest(fuzzy, modifiedURL, headers, modifiedBody, word)

			resp, err := session.Client.Do(req)
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
			} else if resp.StatusCode != 404 { // ignore 404 responses
				utils.RequestOutput(pathAndQuery, fuzzy.SearchString, fmt.Sprintf("%d", resp.StatusCode), string(responseBody), modifiedBody, fuzzy.NoSearch)
			}
		}(word)
		if fuzzy.Delay > 0 {
			time.Sleep(time.Duration(fuzzy.Delay) * time.Second)
		}
	}
	wg.Wait()
}
