package commands

import "github.com/spf13/cobra"

// NewDefaultCommand creates a new default command
func NewDefaultCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "pacmod <subcommand>",
		Short: "Command line tool to assist with packaging Go modules",
	}

	cmd.AddCommand(newPackCommand())

	return &cmd
}
