package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Placeholders for most of the "other" commands. Once a command from this list is properly implemented, remove it from here.
// This function creates placeholder commands
func newCobraOtherCommands(rootCmd *cobra.Command) {
	other := map[string]string{
		// "console":      "Try OpenTofu expressions at an interactive command prompt",
		// "fmt":          "Reformat your configuration in the standard style",
		// "force-unlock": "Release a stuck lock on the current workspace",
		// "get":          "Install or upgrade remote OpenTofu modules",
		// "graph":        "Generate a Graphviz graph of the steps in an operation",
		// "import":       "Associate existing infrastructure with a OpenTofu resource",
		// "login":        "Obtain and save credentials for a remote host",
		// "logout":       "Remove locally-stored credentials for a remote host",
		// "metadata":     "Metadata related commands",
		// "output":       "Show output values from your root module",
		// "providers":    "Show the providers required for this configuration",
		// "refresh":      "Update the state to match remote systems",
		// "show":         "Show the current state or a saved plan",
		// "state":        "Advanced state management",
		// "taint":        "Mark a resource instance as not fully functional",
		// "test":         "Execute integration tests for OpenTofu modules",
		// "untaint":      "Remove the 'tainted' state from a resource instance",
		"version": "Show the current OpenTofu version",
	}
	for cmdName, desc := range other {
		cmd := &cobra.Command{
			Use:                cmdName,
			Short:              desc,
			DisableFlagParsing: false,
			GroupID:            commandGroupIdOther.id(),
		}
		cmd.Run = func(cmd *cobra.Command, args []string) {
			fmt.Println("execute", cmdName)
		}
		rootCmd.AddCommand(cmd)
	}
	newWorkspaceCommand(rootCmd)
}

// This highlights how we can use `ValidArgsFunction` to read valid args from a remote source
// to be returned as suggestions for autocompletion.
func newWorkspaceCommand(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:                "workspace",
		Short:              "Workspace management",
		DisableFlagParsing: false,
		GroupID:            commandGroupIdOther.id(),
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: false,
			HiddenDefaultCmd:  false,
		},
	}
	cmd.Run = func(cmd *cobra.Command, args []string) {
		fmt.Println("execute", "workspace")
	}
	wsSelectCmd := &cobra.Command{
		Use:                "select",
		DisableFlagParsing: false,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
			return []cobra.Completion{"default"}, cobra.ShellCompDirectiveKeepOrder
		},
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	cmd.AddCommand(
		wsSelectCmd,
		&cobra.Command{Use: "delete",
			DisableFlagParsing: false,
			Run: func(cmd *cobra.Command, args []string) {

			}},
		&cobra.Command{Use: "add",
			DisableFlagParsing: false,
			Run: func(cmd *cobra.Command, args []string) {

			}},
		&cobra.Command{Use: "list",
			DisableFlagParsing: false,
			Run: func(cmd *cobra.Command, args []string) {

			}},
	)
	rootCmd.AddCommand(cmd)
}
