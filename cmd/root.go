package cmd

import (
	"flag"
	"fmt"
	"gofuzz/internal/fuzzer"
	"gofuzz/utils"
	"log"
	"os"
	"strings"
	"time"
)

func Execute() {
	targetURL := flag.String("url", "", "specify the host url")
	wordlist := flag.String("wordlist", "", "specify a wordlist used to fuzz")
	method := flag.String("method", "GET", "specify the request method [POST, GET]")
	body := flag.String("body", "", "specify POST request body")
	headers := flag.String("custom-headers", "", "specify the file that contains headers [seperated by line]")
	threads := flag.Int("threads", 3, "specify thread count [default: 3]")
	timeout := flag.Int("timeout", 5, "specify timeout in seconds [default 5]")
	help := flag.Bool("h", false, "show usage")
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	time := time.Duration(*timeout) * time.Second

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fuzzList, err := utils.ReadFile(fmt.Sprintf("%s/%s", cwd, *wordlist))
	if err != nil {
		log.Fatal(err)
	}

	var customHeaders []string
	if *headers != "" {
		customHeaders, err = utils.ReadFile(fmt.Sprintf("%s/%s", cwd, *headers))
		if err != nil {
			log.Fatal(err)
		}
	}

	var links []string
	for _, word := range fuzzList {
		links = append(links, strings.Replace(*targetURL, "FUZZ", word, -1))
	}

	err = fuzzer.GoRequest(
		*method,
		*targetURL,
		customHeaders,
		*body,
		links,
		*threads,
		time)
	if err != nil {
		log.Printf("Error Occurred Sending Request\n -> Error: %v\n", err)
	}
}
