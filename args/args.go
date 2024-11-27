package args

import (
	"errors"
	"fmt"
	"gofuzz/utils"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// struct used by given commands
type Fuzzy struct {
	Method       string
	URL          string
	RequestDump  string
	CustomHeader []string
	Body         string
	WordlistPath string
	StatusCodes  []string
	SearchString string
	NoSearch     bool
	Proxies      string
	Timeout      int
	Threads      int
}

// function to parse the flags needed for the given commands
// -> add/adjust commands as necessary
func ParseFlags(cmd *cobra.Command, args *Fuzzy) {
	cmd.Flags().StringVarP(&args.Method, "method", "m", "GET", "specify the HTTP method to use")
	cmd.Flags().StringVarP(&args.URL, "url", "u", "", "specify the target URL to fuzz")
	cmd.Flags().StringSliceVarP(&args.CustomHeader, "header", "H", nil, "specify headers (comma-separated")
	cmd.Flags().StringVarP(&args.RequestDump, "request", "R", "", "specify path to custom headers file (line-separated) [accepts burpsuite request dumps]")
	cmd.Flags().StringVarP(&args.Body, "body", "b", "", "specify the request body to fuzz")
	cmd.Flags().StringVarP(&args.WordlistPath, "wordlist", "w", "", "specify path to a wordlist used for fuzzing")
	cmd.Flags().StringSliceVarP(&args.StatusCodes, "status-codes", "c", nil, "specify status codes to filter by (comma-separated)")
	cmd.Flags().StringVarP(&args.SearchString, "search", "s", "", "specify the string to search/filter in response body")
	cmd.Flags().BoolVarP(&args.NoSearch, "no-search", "N", false, "output responses that do NOT contain the search string")
	cmd.Flags().StringVarP(&args.Proxies, "proxies", "p", "", "specify proxy(ies) (line-separated) [accepts file]")
	cmd.Flags().IntVarP(&args.Timeout, "timeout", "T", 10, "specify a timeout per request")
	cmd.Flags().IntVarP(&args.Threads, "threads", "t", 5, "specify the number of threads to use")

	fmt.Println(args.WordlistPath)

	if args.Method != "GET" {
		if _, err := os.Stat(args.Body); !errors.Is(err, os.ErrNotExist) {
			body, err := utils.ReadFile(args.Body)
			if err != nil {
				log.Fatalf("[-] Error opening file containing the request body %s\n -> Error: %v\n", args.Body, err)
			}
			args.Body = strings.Join(body, "\n")
		}
	}
}
