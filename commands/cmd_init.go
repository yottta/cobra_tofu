package commands

import (
	"flag"
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

	cfg := &initCfg{}
	rootCmd.AddCommand(cmd)

	flagSet := configureInitCobraFlags(cfg, cmd.Flags())
	// flagSet.Usage = func() {
	// 	helpText := commandHelp()(cmd)
	// 	fmt.Println(helpText)
	// }
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		uiOut := ui.build()
		nonFlagArgs, exitCode := parseFlags(flagSet, args, cfg)
		if exitCode > 0 {
			return &ExitCodeError{ExitCode: exitCode}
		}

		acts := initActs{
			initCfg: cfg,
		}
		uiOut(nonFlagArgs, acts)
		return &ExitCodeError{ExitCode: 0}
	}

}

func configureInitCobraFlags(flags *initCfg, cmdFlags *pflag.FlagSet) *flag.FlagSet {
	basicFlags := &flag.FlagSet{}
	cmdFlags.BoolVar(&flags.flagBackend, "backend", true, "Disable backend or cloud backend initialization for this configuration and use what was previously initialized instead.")
	cmdFlags.BoolVar(&flags.flagCloud, "cloud", true, "")
	cmdFlags.Var(&flags.flagConfigExtra, "backend-config", "Configuration to be merged with what is in the configuration file's 'backend' block. This can be either a path to an HCL file with key/value assignments (same format as terraform.tfvars) or a 'key=value' format, and can be specified multiple times. The backend type must be in the configuration itself.")
	cmdFlags.StringVar(&flags.flagFromModule, "from-module", "", "Copy the contents of the given module into the target directory before initialization.")
	cmdFlags.BoolVar(&flags.flagGet, "get", true, "Disable downloading modules for this configuration.")
	cmdFlags.BoolVar(&flags.flagUpgrade, "upgrade", false, "Install the latest module and provider versions allowed within configured constraints, overriding the default behavior of selecting exactly the version recorded in the dependency lockfile.")
	cmdFlags.Var(&flags.flagPluginPath, "plugin-dir", "Directory containing plugin binaries. This overrides all default search paths for plugins, and prevents the automatic installation of plugins. This flag can be used multiple times.")
	cmdFlags.StringVar(&flags.flagLockfile, "lockfile", "", `Set a dependency lockfile mode. Currently only "readonly" is valid.`)
	cmdFlags.StringVar(&flags.testsDirectory, "test-directory", "tests", `Set the OpenTofu test directory, defaults to "tests". When set, the test command will search for test files in the current directory and in the one specified by the flag.`)

	cmdFlags.CopyToGoFlagSet(basicFlags)
	return basicFlags
}

func parseFlags(cmdFlags *flag.FlagSet, args []string, flags *initCfg) ([]string, int) {
	if err := cmdFlags.Parse(args); err != nil {
		return nil, 1
	}

	flags.backendFlagSet = FlagIsSet(cmdFlags, "backend")
	flags.cloudFlagSet = FlagIsSet(cmdFlags, "cloud")

	switch {
	case flags.backendFlagSet && flags.cloudFlagSet:
		fmt.Println("The -backend and -cloud options are aliases of one another and mutually-exclusive in their use")
		return nil, 1
	case flags.backendFlagSet:
		flags.flagCloud = flags.flagBackend
	case flags.cloudFlagSet:
		flags.flagBackend = flags.flagCloud
	}
	return cmdFlags.Args(), 0
}

// FlagIsSet returns whether a flag is explicitly set in a set of flags
func FlagIsSet(flags *flag.FlagSet, name string) bool {
	isSet := false
	flags.Visit(func(f *flag.Flag) {
		if f.Name == name {
			isSet = true
		}
	})
	return isSet
}
