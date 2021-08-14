package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/zrb-inc/spring"
)

var (
	path string
	dry  bool

	rootCmd = &cobra.Command{
		Use:   "spring-cli",
		Short: "spring-cli generate compile ioc code.",
	}
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func main() {
	log.SetPrefix(fmt.Sprintf("%s%s %s", Green, "[spring]", Reset))
	Execute()
}

func Execute() bool {
	defer func() {
		c := recover()
		log.Print("error : ", c)
	}()
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", ".", "-p or --path")
	rootCmd.PersistentFlags().BoolVarP(&dry, "dry", "d", false, "-d or --dry")
	spring.RegisterCmds(rootCmd)
	rootCmd.Execute()
	return true
}
