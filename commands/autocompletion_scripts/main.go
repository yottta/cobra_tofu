package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yottta/cobra_tofu/commands"
)

//go:generate go run ./main.go gen_completion

func main() {
	cmd := commands.CobraCommands()
	newCobraCompletionCommand(cmd)
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}

// newCobraCompletionCommand shows how we can customise autocompletion for different shells with cobra
// and with posener/complete.
func newCobraCompletionCommand(rootCmd *cobra.Command) {
	completionCmd := &cobra.Command{
		Use:  "gen_completion",
		Args: cobra.NoArgs,
	}

	completionCmd.RunE = func(cmd *cobra.Command, args []string) error {
		for _, target := range []string{"bash", "zsh", "fish", "powershell"} {
			fileName := fmt.Sprintf("tofu.%s", target)
			f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
			if err != nil {
				return fmt.Errorf("could not create file to generate script for %s: %w", target, err)
			}
			switch target {
			case "bash":
				if err := cmd.Root().GenBashCompletion(f); err != nil {
					return err
				}
			case "zsh":
				if err := cmd.Root().GenZshCompletion(f); err != nil {
					return err
				}
			case "fish":
				if err := cmd.Root().GenFishCompletion(f, true); err != nil {
					return err
				}
			case "powershell":
				if err := cmd.Root().GenPowerShellCompletionWithDesc(f); err != nil {
					return err
				}
			}
		}

		return nil
	}
	rootCmd.AddCommand(completionCmd)
}
