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
	if !strings.HasSuffix(rawRequest, "\r\n\r\n") {
		rawRequest += "\r\n\r\n"
	}

	reader := bufio.NewReader(strings.NewReader(rawRequest))
	req, err := http.ReadRequest(reader)
	if err != nil {
		return "", "", nil, "", err
	}

	if req.Host == "" || req.URL.RequestURI() == "" {
		return "", "", nil, "", err
	}

	req.URL.Scheme = "https"

	var bodyBytes []byte
	bodyBytes, err = io.ReadAll(reader)
	if err != nil && err != io.EOF {
		return "", "", nil, "", err
	}

	targetURL := fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.Host, req.URL.RequestURI())

	for name, values := range req.Header {
		for _, value := range values {
			headers = append(headers, fmt.Sprintf("%s: %s", name, value))
		}
	}

	return targetURL, req.Method, headers, string(bodyBytes), nil
}

func RequestOutput(path string, search string, statuscode string, responsebody string, requestBody string, noSearch bool) {
	if noSearch {
		if !strings.Contains(string(responsebody), search) {
			fmt.Printf("Path: %-40s [%s] Length: %-10d\n", path, statuscode, len(responsebody))
			fmt.Printf("Request Body: %-40s\n\n", requestBody)
		}
	} else if search == "" || strings.Contains(string(responsebody), search) { // get paths that contain search string
		fmt.Printf("Path: %-40s [%s] Length: %-10d\n", path, statuscode, len(responsebody))
		fmt.Printf("Request Body: %-40s\n\n", requestBody)
	}
}
