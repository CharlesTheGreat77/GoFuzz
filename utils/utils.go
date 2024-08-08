package utils

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// function to read file contents
func ReadFile(path string) ([]string, error) {
	var contents []string
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		contents = append(contents, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return contents, nil
}

// function to parse burpsuite request dump
func ParseBurpRequest(rawRequest string) (URL string, method string, headers []string, body string, err error) {
	fmt.Println("reading request..")
	reader := bufio.NewReader(strings.NewReader(rawRequest))
	req, err := http.ReadRequest(reader)
	if err != nil {
		return "", "", nil, "", err
	}
	req.URL.Scheme = "https"

	rawBody, err := io.ReadAll(reader)
	if err != nil {
		return "", "", nil, "", err
	}

	targetURL := fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.Host, req.URL.RequestURI())
	for name, values := range req.Header {
		for _, value := range values {
			headers = append(headers, fmt.Sprintf("%s: %s", name, value))
		}
	}

	return targetURL, req.Method, headers, string(rawBody), nil
}
