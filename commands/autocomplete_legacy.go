package commands

import (
	"github.com/posener/complete"
	"github.com/spf13/cobra"
)

func ExecuteLegacyAutocompletion(cmd *cobra.Command) bool {
	// We do this to check if the user still has the old autocompletion scripts.
	// If it does and autocompletion got triggered, this will return "true" allowing the
	// logic below to execute and provide legacy autocompletion.
	if !complete.New(cmd.Name(), complete.Command{}).Complete() {
		return false
	}

	// Build the root command
	autocomplete := complete.New(cmd.Name(), generateAutocompletionData(cmd))
	return autocomplete.Complete()
}

func generateAutocompletionData(cmd *cobra.Command) complete.Command {
	var completeCmd complete.Command
	walkFn := func(k string, cmd *cobra.Command) {
		// Ignore the empty key which can be present for default commands.
		if !cmd.HasParent() {
			return
		}
		if cmd.Hidden {
			return
		}

		if _, ok := completeCmd.Sub[k]; ok {
			// If we already tracked this subcommand then ignore
			return
		}

		if completeCmd.Sub == nil {
			completeCmd.Sub = make(map[string]complete.Command)
		}
		subCmd := generateAutocompletionData(cmd)

		completeCmd.Sub[k] = subCmd
		return
	}

	walkCommand := func(cmd *cobra.Command, f func(k string, cmd *cobra.Command)) {
		if cmd == nil {
			return
		}
		for _, c := range cmd.Commands() {
			f(c.Name(), c)
		}
	}
	walkCommand(cmd, walkFn)
	return completeCmd
}
