package cmd

import (
	"flag"
	"gofuzz/internal/fuzzer"
	"gofuzz/utils"
	"log"
	"strings"
	"time"
)

func Execute() {
	link := flag.String("url", "", "specify the host url")
	burpsuite := flag.String("burp", "", "specify path to burp request")
	wordlist := flag.String("wordlist", "", "specify a wordlist used to fuzz")
	method := flag.String("method", "GET", "specify the request method [POST, GET]")
	requestBody := flag.String("body", "", "specify POST request body (or file containing the body)")
	headers := flag.String("custom-headers", "", "specify the file that contains headers [separated by line]")
	threads := flag.Int("threads", 3, "specify thread count [default: 3]")
	statusCode := flag.String("statuscode", "", "specify a status code(s) to output")
	timeout := flag.Int("timeout", 5, "specify timeout in seconds [default 5]")
	help := flag.Bool("h", false, "show usage")
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	timeoutDuration := time.Duration(*timeout) * time.Second

	fuzzList, err := utils.ReadFile(*wordlist)
	if err != nil {
		log.Fatalf("[-] Error Occurred reading wordlist file\n -> Error: %v\n", err)
	}

	targetURL := *link
	var customHeaders []string
	var body string
	requestMethod := *method

	if *burpsuite != "" {
		request, err := utils.ReadFile(*burpsuite)
		if err != nil {
			log.Printf("[-] Error Occurred reading burp request file\n -> Error: %v\n", err)
		}
		rawRequest := strings.Join(request, "\n")
		targetURL, requestMethod, customHeaders, body, err = utils.ParseBurpRequest(rawRequest)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *headers != "" {
		customHeaders, err = utils.ReadFile(*headers)
		if err != nil {
			log.Fatalf("[-] Error Occurred reading headers file\n -> Error: %v\n", err)
		}
	}

	var statuscodeSplit []string
	if *statusCode != "" {
		statuscodeSplit = strings.Split(*statusCode, ",")
	}

	if *requestBody != "" {
		contents, err := utils.ReadFile(*requestBody)
		if err == nil { // get contents of file for body
			body = strings.Join(contents, "\n")
		}
	}

	fuzzer.GoRequest(
		requestMethod,
		targetURL,
		customHeaders,
		body,
		fuzzList,
		*threads,
		timeoutDuration,
		statuscodeSplit)
}
