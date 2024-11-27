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

// pathtraversalCmd represents the pathtraversal command
var pathtraversalCmd = &cobra.Command{
	Use:   "pathtraversal",
	Short: "Fuzzes URLs to detect path traversal vulnerabilities using a wordlist.",
	Long: `The "pathtraversal" command is used to test endpoints for potential path traversal vulnerabilities by injecting and incrementing "../" sequences into the URL. 

Path traversal attacks attempt to access files and directories stored outside the web root directory by manipulating variables referencing file paths. This command automates the process of crafting such payloads and sending requests to identify misconfigurations or security flaws in the target application.

Key features include:
  - Automatic generation and injection of path traversal payloads (e.g., "../../", "../../../").
  - Support for appending or replacing parts of the URL to craft dynamic attacks.

Example usage:

  gofuzz pathtraversal -u https://example.com/path/to/file?parameter=FUZZ -w paths.txt
  gofuzz pathtraversal -u https://example.com/api/v1/files/FUZZ -w sensitive_files.txt -s "root:x:"
  gofuzz pathtraversal -u https://example.com/#parameter=FUZZ -sc 200 -timeout 5 -t 10
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

		fuzz.GoTraverse(fuzzy, headers, proxies, wordlist)
	},
}

func init() {
	args.ParseFlags(pathtraversalCmd, &fuzzy)
}
