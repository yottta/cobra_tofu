package commands

import (
	"context"

	"github.com/spf13/pflag"
)

type initFlags struct {
	flagFromModule, flagLockfile, testsDirectory string
	flagBackend, flagCloud, flagGet, flagUpgrade bool
	flagPluginPath                               FlagStringSlice
	backendFlagSet                               bool
	cloudFlagSet                                 bool
	flagConfigExtra                              rawFlags
}

func (c *initFlags) configureFlags(cmdFlags *pflag.FlagSet) {
	cmdFlags.BoolVar(&c.flagBackend, "backend", true, "Disable backend or cloud backend initialization for this configuration and use what was previously initialized instead.")
	cmdFlags.BoolVar(&c.flagCloud, "cloud", true, "")
	cmdFlags.Var(&c.flagConfigExtra, "backend-config", "Configuration to be merged with what is in the configuration file's 'backend' block. This can be either a path to an HCL file with key/value assignments (same format as terraform.tfvars) or a 'key=value' format, and can be specified multiple times. The backend type must be in the configuration itself.")
	cmdFlags.StringVar(&c.flagFromModule, "from-module", "", "Copy the contents of the given module into the target directory before initialization.")
	cmdFlags.BoolVar(&c.flagGet, "get", true, "Disable downloading modules for this configuration.")
	cmdFlags.BoolVar(&c.flagUpgrade, "upgrade", false, "Install the latest module and provider versions allowed within configured constraints, overriding the default behavior of selecting exactly the version recorded in the dependency lockfile.")
	cmdFlags.Var(&c.flagPluginPath, "plugin-dir", "Directory containing plugin binaries. This overrides all default search paths for plugins, and prevents the automatic installation of plugins. This flag can be used multiple times.")
	cmdFlags.StringVar(&c.flagLockfile, "lockfile", "", `Set a dependency lockfile mode. Currently only "readonly" is valid.`)
	cmdFlags.StringVar(&c.testsDirectory, "test-directory", "tests", `Set the OpenTofu test directory, defaults to "tests". When set, the test command will search for test files in the current directory and in the one specified by the flag.`)
}

func (c *initFlags) initBackend(ctx context.Context, extraConfig rawFlags) error {
	// TODO return the backend configured with the information from the flags
	return nil
}
