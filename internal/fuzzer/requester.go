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

func GoRequest(method string, targetURL string, headers []string, body string, wordlist []string, maxConcurrentRequests int, timeout time.Duration) error {
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
			} else {
				req, err = http.NewRequest(method, modifiedURL, nil)
			}
			if err != nil {
				fmt.Printf("Error creating request for %s: %v\n", targetURL, err)
				return
			}

			if len(headers) > 0 {
				for _, line := range headers {
					header := strings.TrimSpace(line)
					splitHeader := strings.SplitN(header, ":", 2)
					if len(splitHeader) == 2 {
						req.Header.Set(splitHeader[0], splitHeader[1])
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
			defer resp.Body.Close()

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

			fmt.Printf("Path: %s\tResponse Code: %d\nResponse Length: %d\nRequest Body: %s\n\n", pathAndQuery, resp.StatusCode, len(responseBody), modifiedBody)
		}(word)
	}
	wg.Wait()

	return nil
}
