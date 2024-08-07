package cmd

import (
	"flag"
	"gofuzz/internal/fuzzer"
	"gofuzz/utils"
	"log"
	"time"
)

func Execute() {
	targetURL := flag.String("url", "", "specify the host url")
	wordlist := flag.String("wordlist", "", "specify a wordlist used to fuzz")
	method := flag.String("method", "GET", "specify the request method [POST, GET]")
	body := flag.String("body", "", "specify POST request body")
	headers := flag.String("custom-headers", "", "specify the file that contains headers [separated by line]")
	threads := flag.Int("threads", 3, "specify thread count [default: 3]")
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
		log.Fatal(err)
	}

	var customHeaders []string
	if *headers != "" {
		customHeaders, err = utils.ReadFile(*headers)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = fuzzer.GoRequest(
		*method,
		*targetURL,
		customHeaders,
		*body,
		fuzzList,
		*threads,
		timeoutDuration)
	if err != nil {
		log.Printf("Error Occurred Sending Request\n -> Error: %v\n", err)
	}
}
