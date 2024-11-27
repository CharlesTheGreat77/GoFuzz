package goclient

import (
	"bytes"
	"crypto/tls"
	"gofuzz/args"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
)

type ProxyRotator struct {
	proxies []string
	mu      sync.Mutex
	index   int
}

// init a New ProxyRotator with a list of proxies
func NewProxyRotator(proxies []string) *ProxyRotator {
	return &ProxyRotator{proxies: proxies}
}

// GetNextProxy returns the next proxy in the rotation
func (pr *ProxyRotator) GetNextProxy() (*url.URL, error) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	if len(pr.proxies) == 0 {
		return nil, nil // no proxies available
	}

	proxy := pr.proxies[pr.index]
	pr.index = (pr.index + 1) % len(pr.proxies) // rotate index

	return url.Parse(proxy)
}

// function to return a client with given http request configuration(s)
// -> timeout value(s) and proxy rotation(s)
func GoClient(timeout time.Duration, proxyRotator *ProxyRotator) *http.Client {
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
			Proxy: func(req *http.Request) (*url.URL, error) {
				if proxyRotator != nil {
					return proxyRotator.GetNextProxy() // attempt to rotate proxies in slice
				}
				return nil, nil // no proxy if rotator is nil
			},
		},
		Timeout: timeout, // max time overall for the request
	}

	return client
}

// function to set the HTTP request
func GoRequest(fuzzy args.Fuzzy, url string, headers []string, body string, word string) (*http.Request, error) {
	var req *http.Request
	var err error

	if fuzzy.Method == "POST" { // if POST, set body, otherwise just set the url
		req, err = http.NewRequest(fuzzy.Method, fuzzy.URL, bytes.NewReader([]byte(body)))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		// default User-Agent
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.6 Safari/605.1.1")
		req.ContentLength = int64(len(body))
	} else { // GET method... obviously..
		req, err = http.NewRequest(fuzzy.Method, url, nil)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		log.Printf("Error creating request for %s: %v\n", fuzzy.URL, err)
		return nil, err
	}

	// set custom headers (ie. User-Agent etc.)
	if fuzzy.CustomHeader != nil {
		for _, header := range fuzzy.CustomHeader {
			splitHeader := strings.SplitN(header, ":", 2)
			if len(splitHeader) == 2 {
				req.Header.Set(strings.TrimSpace(splitHeader[0]), strings.TrimSpace(splitHeader[1]))
			}
		}
	}

	if len(headers) > 0 { // add custom headers to request
		updatedHeaders := strings.Join(headers, "\n")
		updatedHeaders = strings.Replace(updatedHeaders, "FUZZ", word, -1) // replace occurence of FUZZ with word in list
		headers := strings.Split(updatedHeaders, "\n")
		for _, line := range headers {
			header := strings.TrimSpace(line)
			splitHeader := strings.SplitN(header, ":", 2)
			if len(splitHeader) == 2 {
				key := strings.TrimSpace(splitHeader[0])
				value := strings.TrimSpace(splitHeader[1])
				req.Header.Set(key, value)
			}
		}
	}

	return req, nil
}

// function for debugging the request(s) sent
func DebugRequest(req *http.Request) {
	requestDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Printf("Error dumping request: %v\n", err)
		return
	}
	log.Printf("Full Request:\n%s\n", string(requestDump))
}
