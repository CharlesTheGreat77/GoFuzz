package fuzzer

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// function to fuzz the parameters in a given wordlist based on position in url or body
func GoRequest(method string, targetURL string, customHeaders []string, body string, wordlist []string, maxConcurrentRequests int, timeout time.Duration, statusCodes []string) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
			DialContext: (&net.Dialer{
				Timeout: timeout,
			}).DialContext,
			IdleConnTimeout:       timeout,
			TLSHandshakeTimeout:   timeout,
			ExpectContinueTimeout: timeout,
		},
		Timeout: timeout,
	}

	semaphore := make(chan struct{}, maxConcurrentRequests)
	var wg sync.WaitGroup

	for _, word := range wordlist {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(word string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			encodedWord := url.QueryEscape(word)
			modifiedURL := strings.Replace(targetURL, "FUZZ", encodedWord, -1)
			modifiedBody := strings.Replace(body, "FUZZ", encodedWord, -1)

			var req *http.Request
			var err error
			if method == "POST" {
				req, err = http.NewRequest(method, targetURL, bytes.NewReader([]byte(modifiedBody)))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				req.ContentLength = int64(len(modifiedBody))
			} else {
				req, err = http.NewRequest(method, modifiedURL, nil)
			}
			if err != nil {
				log.Printf("Error creating request for %s: %v\n", targetURL, err)
				return
			}
			if req == nil {
				log.Printf("Error: Request is nil for %s\n", targetURL)
				return
			}

			if len(customHeaders) > 0 {
				for _, line := range customHeaders {
					header := strings.TrimSpace(line)
					splitHeader := strings.SplitN(header, ":", 2)
					if len(splitHeader) == 2 {
						req.Header.Set(splitHeader[0], splitHeader[1])
					} else {
						fmt.Printf("Invalid header format: %s\n", line)
					}
				}
			}

			// debugging request shii
			// requestDump, err := httputil.DumpRequestOut(req, true)
			// if err != nil {
			// 	fmt.Printf("Error dumping request: %v\n", err)
			// 	return
			// }
			// fmt.Printf("Full Request:\n%s\n", string(requestDump))

			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Error sending request to %s: %v\n", targetURL, err)
				return
			}

			defer func() {
				if err := resp.Body.Close(); err != nil {
					fmt.Printf("Error closing response body for %s: %v\n", targetURL, err)
				}
			}()

			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error reading response from %s: %v\n", targetURL, err)
				return
			}

			parsedURL, err := url.Parse(modifiedURL)
			if err != nil {
				log.Printf("Error parsing URL %s: %v\n", modifiedURL, err)
				return
			}
			pathAndQuery := parsedURL.Path + parsedURL.RawQuery

			if len(statusCodes) != 0 {
				sc := fmt.Sprintf("%d", resp.StatusCode)
				for _, code := range statusCodes {
					if string(code) == string(sc) {
						fmt.Printf("Path: %s\tResponse Code: %s\nResponse Length: %d\nRequest Body: %s\n\n", pathAndQuery, sc, len(responseBody), modifiedBody)
					}
				}
			} else {
				fmt.Printf("Path: %s\tResponse Code: %d\nResponse Length: %d\nRequest Body: %s\n\n", pathAndQuery, resp.StatusCode, len(responseBody), modifiedBody)
			}
		}(word)
	}
	wg.Wait()
}
