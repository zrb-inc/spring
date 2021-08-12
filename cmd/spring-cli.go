package main

import (
	"github.com/spf13/cobra"
)

var (
	path string
	dry  bool

	rootCmd = &cobra.Command{
		Use:   "spring-cli",
		Short: "spring-cli generate compile ioc code.",
	}
)

func main() {
	Execute()
}

func Execute() bool {
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", ".", "-p or --path")
	rootCmd.PersistentFlags().BoolVarP(&dry, "path", "p", false, "-p or --path")
	rootCmd.Execute()
	return true
}
