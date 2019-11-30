package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	version = "0.1.0"
)

// NewVersionCmd ..
func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Shows version",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runversionCmd(cmd, args)
		},
	}

	return cmd
}

func runversionCmd(cmd *cobra.Command, args []string) error {
	fmt.Printf("%s\n", version)
	return nil
}
