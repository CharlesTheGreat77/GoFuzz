package fuzz

import (
	"fmt"
	"gofuzz/args"
	"gofuzz/internal/goclient"
	"gofuzz/utils"
	"io"
	"log"
	"strings"
	"sync"
	"time"
)

func GoTraverse(fuzzy args.Fuzzy, headers []string, proxies []string, wordlist []string) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, fuzzy.Threads)

	timeout := time.Duration(fuzzy.Timeout) * time.Second

	fmt.Println(fuzzy.URL)
	for _, word := range wordlist {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(word string) {
			defer wg.Wait()
			defer func() { <-semaphore }()
			session := goclient.NewSession(timeout, goclient.NewProxyRotator(proxies))

			for i := 0; i < 20; i++ { // max of 20 '../word' increments
				traversalPath := strings.Repeat("../", i) + word
				url := strings.Replace(fuzzy.URL, "FUZZ", traversalPath, -1)

				req, err := goclient.GoRequest(fuzzy, url, headers, fuzzy.Body, word)
				if err != nil {
					log.Printf("[-] Error creating a request for %s\n -> Error: %v\n", url, err)
				}

				//goclient.DebugRequest(req)
				resp, err := session.Client.Do(req)
				if err != nil {
					log.Printf("[-] Error sending request to %s\n -> Error: %v\n", url, err)
				}

				responseBody, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("[-] Error reading the response body %s\n -> Error: %v", resp.Request.URL, err)
				}

				path := req.URL.RequestURI()

				if len(fuzzy.StatusCodes) > 0 {
					sc := fmt.Sprintf("%d", resp.StatusCode)
					for _, code := range fuzzy.StatusCodes {
						if string(code) == sc {
							utils.RequestOutput(path, fuzzy.SearchString, sc, fuzzy.Body, "", fuzzy.NoSearch)
						}
					}
				} else if resp.StatusCode != 404 || resp.StatusCode != 504 || resp.StatusCode != 400 { // ignore 404, 400, 504 by default
					utils.RequestOutput(path, fuzzy.SearchString, fmt.Sprintf("%d", resp.StatusCode), string(responseBody), fuzzy.Body, fuzzy.NoSearch)
				}
			}
		}(word)
	}
	wg.Wait()
}
