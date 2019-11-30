package cli

import "github.com/spf13/cobra"

// NewRootCmd ..
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "<subcommand>",
		Short:   "runs util commands with VCL",
		Version: version,
	}

	cmd.AddCommand(newREPLCmd())
	cmd.AddCommand(newVersionCmd())
	return cmd
}

// Execute ...
func Execute() error {
	rootCmd := newRootCmd()
	return rootCmd.Execute()
}
