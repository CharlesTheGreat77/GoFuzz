package cmd

import (
	"errors"
	"gofuzz/args"
	"gofuzz/internal/fuzz"
	"gofuzz/utils"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// fuzzCmd represents the fuzz command
var fuzzCmd = &cobra.Command{
	Use:   "fuzz",
	Short: "fuzz URL parameters, HTTP request bodies, and headers",
	Long: `the fuzz command is used to FUZZ parameters or points in an HTTP request.
  Fuzz can be used by specifing:

  method              -> specify a POST or GET method [Default: GET]
  u [url]             -> specify the target url
  H [headers]         -> specify path to custom headers file (seperate by line) 
                         [accepts burpsuite request dumps]
  body                -> specify the body to FUZZ for a given POST request
  w [wordlist]        -> specify a wordlist used to FUZZ the given parameters
  c [status-codes] -> specify status code(s) [seperated by comma: 200,403]
  s [search]          -> specify a string to search/filter for in the response body
  N [NoSearch]        -> enable to output responses that do NOT contain the search string
  timeout             -> specify the time for timeout for each request
  t [threads]         -> specify the number of threads

  example(s):
  
  gofuzz fuzz -w common.txt -u https://example.com/FUZZ -sc 200,403
  gofuzz fuzz -w commont.txt -u https://ex.com/#id=FUZZ -sc 200
  gofuzz fuzz -w passw.txt -method POST -u https://example.com/ -s "login failed" -N -body body.json
  gofuzz fuzz -w list.txt -u https://example.com/FUZZ -s "passw" -sc 200 -timeout 5 -t 10

  `,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			headers []string
			proxies []string
		)
		wordlist, err := utils.ReadFile(fuzzy.WordlistPath)
		if err != nil {
			log.Fatalf("[-] Error opening wordlist file %s\n -> Error: %v\n", fuzzy.WordlistPath, err)
		}

		if _, err := os.Stat(fuzzy.Body); !errors.Is(err, os.ErrNotExist) { // if body is a file
			contents, err := utils.ReadFile(fuzzy.Body)
			if err != nil {
				log.Fatalf("[-] Error reading file for body of request(s) %s\n -> Error: %v\n", fuzzy.Body, err)
			}

			fuzzy.Body = strings.Join(contents, "\n")
		}

		if fuzzy.RequestDump != "" {
			contents, err := utils.ReadFile(fuzzy.RequestDump)
			if err != nil {
				log.Fatalf("[-] Error opening file contain the HTTP request dump %s\n -> Error: %v\n", fuzzy.RequestDump, err)
			}

			headerDump := strings.Join(contents, "\n")
			fuzzy.URL, fuzzy.Method, headers, fuzzy.Body, err = utils.ParseBurpRequest(headerDump)
			if err != nil {
				log.Fatalf("[-] Error parsing HTTP request dump %s\n -> Error: %v\n", fuzzy.RequestDump, err)
			}
		}

		if fuzzy.Proxies != "" {
			proxies, err = utils.ReadFile(fuzzy.Proxies)
			if err != nil {
				log.Fatalf("[-] Error opening the file containing proxy(ies)%s\n -> Error: %v\n", fuzzy.Proxies, err)
			}
		}

		fuzz.GoIntruder(fuzzy, wordlist, headers, proxies)
	},
}

func init() {
	args.ParseFlags(fuzzCmd, &fuzzy)
}
