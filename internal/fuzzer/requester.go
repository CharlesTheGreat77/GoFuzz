package fuzzer

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
)

// function to fuzz the parameters in a given wordlist based on position in url or body
func GoRequest(method string, targetURL string, customHeaders []string, body string, wordlist []string, maxConcurrentRequests int, timeout time.Duration, statusCodes []string) {
	// configuration for each request
	client := &http.Client{
		Transport: &http.Transport{
			TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
			DialContext: (&net.Dialer{
				Timeout: timeout, // max time for TCP handshake
			}).DialContext,
			IdleConnTimeout:       timeout, // max time for idle conn. to remain open
			TLSHandshakeTimeout:   timeout, // max time for TLS handshake
			ExpectContinueTimeout: timeout, // max time to wait for server response
		},
		Timeout: timeout, // max time overall for the request
	}

	semaphore := make(chan struct{}, maxConcurrentRequests)
	var wg sync.WaitGroup

	for _, word := range wordlist {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(word string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			// replace FUZZ for each line
			encodedWord := url.QueryEscape(word) // URL encode
			modifiedURL := strings.Replace(targetURL, "FUZZ", encodedWord, -1)
			modifiedBody := strings.Replace(body, "FUZZ", word, -1)

			var req *http.Request
			var err error
			if method == "POST" { // if POST, set body, otherwise just set the url
				req, err = http.NewRequest(method, targetURL, bytes.NewReader([]byte(modifiedBody)))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				req.ContentLength = int64(len(modifiedBody))
			} else { // GET method... obviously..
				req, err = http.NewRequest(method, modifiedURL, nil)
			}
			if err != nil {
				log.Printf("Error creating request for %s: %v\n", targetURL, err)
				return
			}

			if len(customHeaders) > 0 { // if custom headers were set, replace FUZZ in corresponding section
				updatedHeaders := strings.Join(customHeaders, "\n")
				updatedHeaders = strings.Replace(updatedHeaders, "FUZZ", word, -1) // globally replace FUZZ words in headers
				headers := strings.Split(updatedHeaders, "\n")
				for _, line := range headers {
					header := strings.TrimSpace(line)
					splitHeader := strings.SplitN(header, ":", 2)
					if len(splitHeader) == 2 {
						key := strings.TrimSpace(splitHeader[0])
						value := strings.TrimSpace(splitHeader[1])
						req.Header.Set(key, value)
					} else {
						log.Fatalf("Invalid header format: %s\n", line)
					}
				}
			}

			// debugging request shii for each request sent
			// requestDump, err := httputil.DumpRequestOut(req, true)
			// if err != nil {
			//	fmt.Printf("Error dumping request: %v\n", err)
			//	return
			// }
			// fmt.Printf("Full Request:\n%s\n", string(requestDump))

			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error sending request to %s: %v\n", targetURL, err)
				return
			}

			defer func() {
				if err := resp.Body.Close(); err != nil {
					log.Printf("Error closing response body for %s: %v\n", targetURL, err)
				}
			}()

			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading response from %s: %v\n", targetURL, err)
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
						fmt.Printf("Path: %-40s [%s] Length: %-10d\n", pathAndQuery, sc, len(responseBody))
						fmt.Printf("Request Body: %-40s\n\n", modifiedBody)
					}
				}
			} else if resp.StatusCode != 404 { // ignore 404 responses
				fmt.Printf("Path: %-40s [%d] Length: %-10d\n", pathAndQuery, resp.StatusCode, len(responseBody))
				fmt.Printf("Request Body: %-40s\n\n", modifiedBody)
			}
		}(word)
	}
	wg.Wait()
}
