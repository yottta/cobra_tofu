package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func CobraCommands() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:  "cobra_tofu",
		Long: "The available commands for execution are listed below. The primary workflow commands are given first, followed by less common or more advanced commands.",

		// These 2 are needed to disable printing usage and errors because each command will return
		// an error type with the exit code, even the command execution succeeds.
		// This goes hand in hand with [ExitCodeError].
		SilenceUsage:  true,
		SilenceErrors: true,

		// TraverseChildren is needed to ensure that the flags from the parent command are not passed
		// as arguments in any of its subcommands.
		// Instead, by having this enabled, the flags are parsed for any parent command into their
		// dedicated pointers.
		TraverseChildren: true,
	}
	// adding groups to be able to group the commands accordingly in the help text
	rootCmd.AddGroup(commandGroupIdMain.group(), commandGroupIdOther.group())

	// allows customisation of error messages when a flag parsing failed
	rootCmd.SetFlagErrorFunc(func(command *cobra.Command, err error) error {
		return fmt.Errorf("failed parsing flag while executing %s: %w", command.Name(), err)
	})

	// if any subcommand needs to use globally defined flags, we can:
	// * use persistent flags and any subcommand will have to register the flags themselves
	// * use a struct like it's used here and pass the struct to other commands initialisation to be able to build the
	//   components dependent on those values
	uif := &uiFlags{}
	uif.configureFlags(rootCmd)

	// basic flags for the root command
	chdir := rootCmd.Flags().String("chdir", "", "Switch to a different working directory before executing the given subcommand")
	rootCmd.Flags().Bool("help", false, "Show this help output, or the help for a specified subcommand")
	rootCmd.Flags().Bool("version", false, `Alias to "version" command`)

	// this way we can have common actions defined once and executed for any other invoked command.
	// We could have here pretty much what we have in realMain.
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if chdir != nil && *chdir != "" {
			err := os.Chdir(*chdir)
			if err != nil {
				return fmt.Errorf("error handling -chdir option: %w", err)
			}
		}
		return nil
	}
	rootCmd.PersistentPostRun = func(cmd *cobra.Command, args []string) {
		wd, _ := os.Getwd()
		fmt.Println("working directory: ", wd)
	}

	// customisation of the help text that is applied to this command and any sub-command of it
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		helpText := commandHelp()(cmd)
		fmt.Println(helpText)
	})

	// create the other sub commands
	newCobraInitCommand(rootCmd, uif) // <- example of passing the UI struct that holds the arguments for the UI to allow the command execution to use those.
	newCobraMainCommands(rootCmd)
	newCobraOtherCommands(rootCmd)
	newCobraCompletionCommand(rootCmd)

	return rootCmd
}

// this struct is used just as POC to see how we could pass
// globally parsed flags to sub commands without delegating the
// flags there
// For example we could use this to generate the view based on
// flags like `-concise` or `-no-color`.
type uiFlags struct {
	NoColor bool
	Concise bool
}

func (uif *uiFlags) configureFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&uif.Concise, "concise", false, "Disables progress-related messages in the output.")
	cmd.Flags().BoolVar(&uif.NoColor, "no-color", false, "If specified, output won't contain any color.")
	// TODO add more flags for the ui
}

func (uif *uiFlags) build() func(...any) {
	var tags []string
	if uif.NoColor {
		tags = append(tags, "[no-color]")
	}
	if uif.Concise {
		tags = append(tags, "[concise]")
	}
	return func(s ...any) {
		all := []any{strings.Join(tags, " ")}
		all = append(all, s...)
		fmt.Println(all...)
	}
}
