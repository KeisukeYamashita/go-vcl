package cli

import (
	"os"

	"github.com/KeisukeYamashita/go-vcl/repl"
	"github.com/spf13/cobra"
)

// NewVersionCmd ..
func newREPLCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "repl",
		Short:   "Starts REPL",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runREPLCmd(cmd, args)
		},
	}

	return cmd
}

func runREPLCmd(cmd *cobra.Command, args []string) error {
	repl.Start(os.Stdin, os.Stdout)
	return nil
}
