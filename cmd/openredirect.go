package cmd

import (
	"errors"
	"fmt"
	"gofuzz/args"
	"gofuzz/internal/fuzz"
	"gofuzz/utils"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// openredirectCmd represents the openredirect command
var openredirectCmd = &cobra.Command{
	Use:   "openredirect",
	Short: "openredirect is used to fuzz parameters in URL(s)/Headers for open redirects",
	Long: `the openredirect command is used to FUZZ parameters in URL(s)/Headers.
  openredirect can be used by specifing:

  method              -> specify a POST or GET method [Default: GET]
  u [url]             -> specify the target url
  H [headers]         -> specify path to custom headers file (seperate by line)
  body                -> specify the body to FUZZ for a given POST request
  w [wordlist]        -> specify a wordlist used to FUZZ the given parameters
  sc [status-code(s)] -> specify status code(s) [seperated by comma: 200,403]
  s [search]          -> specify a string to search/filter for in the response body
  N [NoSearch]        -> enable to output responses that do NOT contain the search string
  timeout             -> specify the time for timeout for each request
  t [threads]         -> specify the number of threads


  example(s):
  gofuzz openredirect -u https://example.com/path?redirect=FUZZ -w urls.txt
  gofuzz openredirect -burp requestDump.txt -w urls.txt
  gofuzz openredirect -u https://example.com/path?return=FUZZ -w urls.txt -sc 302,200


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
		fmt.Println(fuzzy)

		fuzz.GOpenRedirect(fuzzy, headers, proxies, wordlist)
	},
}

func init() {
	args.ParseFlags(openredirectCmd, &fuzzy)
}
