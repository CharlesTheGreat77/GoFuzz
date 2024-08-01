package fuzzer

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
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

			req, err := http.NewRequest(method, modifiedURL, nil)
			if err != nil {
				fmt.Printf("Error creating request for %s: %v\n", modifiedURL, err)
				return
			}

			if method == "POST" {
				req.Body = io.NopCloser(bytes.NewReader([]byte(body)))
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
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Error sending request to %s: %v\n", modifiedURL, err)
				return
			}
			defer resp.Body.Close()

			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error reading response from %s: %v\n", modifiedURL, err)
				return
			}
			fmt.Printf("URL: %s\nResponse Code: %d\nResponse Length: %d\n\n", resp.Request.URL, resp.StatusCode, len(responseBody))
		}(word)
	}
	wg.Wait()

	return nil
}
