/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"gofuzz/args"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type Fuzzie struct {
	Method       string
	URL          string
	Body         string
	Wordlist     []string
	StatusCodes  []string
	SearchString string
	NoSearch     bool
	Timeout      time.Duration
	Threads      int
}

var fuzzy args.Fuzzy

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "GoFuzz",
	Short: "Web application security testing by fuzzing paths, requests, etc.",
	Long: `GoFuzz is a simple yet powerful fuzzing tool written in Go,
  designed for web application security testing. With its concurrent execution model,
  GoFuzz can fuzz endpoints, parameters, and headers efficiently,
  making it a perfect companion for penetration testers and bug hunters.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(cmd.Short)
		cmd.Println(cmd.Long)
	},
}

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	rootCmd.AddCommand(fuzzCmd)
	rootCmd.AddCommand(openredirectCmd)
	rootCmd.AddCommand(pathtraversalCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}
