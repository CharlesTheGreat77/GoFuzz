package fuzz

import (
	"fmt"
	"gofuzz/args"
	"gofuzz/internal/goclient"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"
)

// function to do a simple check for open redirect vulns. in a given param
func GOpenRedirect(fuzzy args.Fuzzy, headers []string, proxies []string, wordlist []string) {
	timeout := time.Duration(fuzzy.Timeout) * time.Second
	semaphore := make(chan struct{}, fuzzy.Threads)
	var wg sync.WaitGroup

	for _, word := range wordlist {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(word string) {
			defer wg.Done()
			defer func() { <-semaphore }()
			session := goclient.NewSession(timeout, goclient.NewProxyRotator(proxies))

			encodedWord := url.QueryEscape(word) // URL encode the occurence??
			modifiedURL := strings.Replace(fuzzy.URL, "FUZZ", encodedWord, -1)

			req, err := goclient.GoRequest(fuzzy, modifiedURL, headers, "", word)

			resp, err := session.Client.Do(req)
			if err != nil {
				log.Printf("[-] Error sending request to %s\n -> Error: %v\n", modifiedURL, err)
				return
			}

			defer resp.Body.Close()

			if resp.StatusCode >= 300 && resp.StatusCode < 400 { // looking for 3XX response
				fmt.Printf("[*] Potential Target -> URL: %-20s -> Redirect: %-15s", fuzzy.URL, modifiedURL)
				location := resp.Header.Get("Location")
				if location != "" && strings.Contains(location, word) {
					fmt.Printf("(Valid) URL: %-20s -> Redirect: %-15s [%d]", fuzzy.URL, modifiedURL, resp.StatusCode) // potential open redirects
				}
			}
		}(word)
	}
	wg.Wait()
}
