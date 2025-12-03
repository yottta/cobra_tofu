package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newCobraInitCommand(rootCmd *cobra.Command, ui *uiFlags) {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Prepare your working directory for other commands",
		Long: `Initialize a new or existing OpenTofu working directory by creating initial files, loading any remote state, downloading modules, etc.

This is the first command that should be run for any new or existing OpenTofu configuration per machine. This sets up all the local data necessary to run OpenTofu that is typically not committed to version control.

This command is always safe to run multiple times. Though subsequent runs may give errors, this command will never delete your configuration or state. Even so, if you have important information, please back it up prior to running this command, just in case.`,
		DisableFlagParsing: false,
		GroupID:            commandGroupIdMain.id(),
		SilenceErrors:      true,
		SilenceUsage:       true,
	}

	flags := &initFlags{}
	flags.configureFlags(cmd.Flags())
	rootCmd.AddCommand(cmd)

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		flags.backendFlagSet = FlagIsSet(cmd.Flags(), "backend")
		flags.cloudFlagSet = FlagIsSet(cmd.Flags(), "cloud")

		switch {
		case flags.backendFlagSet && flags.cloudFlagSet:
			return ExitCodeErr(1, fmt.Errorf("The -backend and -cloud options are aliases of one another and mutually-exclusive in their use"))
		case flags.backendFlagSet:
			flags.flagCloud = flags.flagBackend
		case flags.cloudFlagSet:
			flags.flagBackend = flags.flagCloud
		}
		return nil
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		uiOut := ui.build()

		uiOut(fmt.Sprintf("from init runE: %#v, %v", args, ui))
		return nil
	}
}

// FlagIsSet returns whether a flag is explicitly set in a set of flags
func FlagIsSet(flags *pflag.FlagSet, name string) bool {
	isSet := false
	// only the set flags are visited
	flags.Visit(func(f *pflag.Flag) {
		if f.Name == name {
			isSet = true
		}
	})
	return isSet
}
