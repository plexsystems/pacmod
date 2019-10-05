package commands

import "github.com/spf13/cobra"

// NewDefaultCommand creates the default command
func NewDefaultCommand() *cobra.Command {

	cmd := cobra.Command{
		Use:   "pacmod <subcommand>",
		Short: "Command line tool to assist with packaging Go modules",
	}

	cmd.AddCommand(NewPackCommand())

	return &cmd
}
