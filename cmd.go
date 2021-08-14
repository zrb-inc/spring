package spring

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Cmds = []*cobra.Command{genCmd}

func preCheck(cmd *cobra.Command, args []string) {
	panic(123)
}

var genCmd = &cobra.Command{
	Use:    "gen",
	Short:  "generate compiled code",
	PreRun: preCheck,
	Run:    func(cmd *cobra.Command, args []string) { fmt.Println(12312231) },
}

func RegisterCmds(root *cobra.Command) {
	for _, c := range Cmds {
		root.AddCommand(c)
	}
}
