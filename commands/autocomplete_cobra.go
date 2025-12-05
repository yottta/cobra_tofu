package commands

import (
	"embed"
	"fmt"
	"os"

	"github.com/posener/complete/cmd/install"
	"github.com/spf13/cobra"
)

//go:embed autocompletion_scripts/tofu.*
var autocompletionFs embed.FS

// newCobraCompletionCommand shows how we can customise autocompletion for different shells with cobra
// and with posener/complete.
func newCobraCompletionCommand(rootCmd *cobra.Command) {
	completionCmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",

		Long:      fmt.Sprintf(longDescription, rootCmd.Root().Name()),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	}
	legacy := completionCmd.Flags().Bool("legacy", false, "indicate that you want to install the legacy autocompletion scripts")
	uninstall := completionCmd.Flags().Bool("uninstall", false, "remove the legacy autocompletion scripts")
	outStream := completionCmd.Flags().String("out", "", "specify the file where this ")

	completionCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if legacy != nil && *legacy {
			if uninstall != nil && *uninstall {
				if err := install.Uninstall(cmd.Root().Name()); err != nil {
					return fmt.Errorf("failed to uninstall legacy scripts: %w", err)
				}
				return nil
			}
			if err := install.Install(cmd.Root().Name()); err != nil {
				return fmt.Errorf("failed to install legacy scripts: %w", err)
			}
			return nil
		}

		dat, err := autocompletionFs.ReadFile(fmt.Sprintf("autocompletion_scripts/tofu.%s", args[0]))
		if err != nil {
			return fmt.Errorf("could not generate autocompletion script: %w", err)
		}
		out := os.Stdout
		if outStream != nil && *outStream != "" {
			f, err := os.OpenFile(*outStream, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
			if err != nil {
				return fmt.Errorf("provided file %s could not be opened: %w", *outStream, err)
			}
			out = f
		}

		_, err = out.Write(dat)
		return err
	}
	rootCmd.AddCommand(completionCmd)
}

const longDescription = `To load completions:
Bash:

  $ source <(%[1]s completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ %[1]s completion bash > /etc/bash_completion.d/%[1]s
  # macOS:
  $ %[1]s completion bash > $(brew --prefix)/etc/bash_completion.d/%[1]s

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ %[1]s completion zsh > "${fpath[1]}/_%[1]s"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ %[1]s completion fish | source

  # To load completions for each session, execute once:
  $ %[1]s completion fish > ~/.config/fish/completions/%[1]s.fish

PowerShell:

  PS> %[1]s completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> %[1]s completion powershell > %[1]s.ps1
  # and source this file from your PowerShell profile.
`
