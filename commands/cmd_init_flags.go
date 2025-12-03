package commands

import (
	"context"
)

type initCfg struct {
	flagFromModule, flagLockfile, testsDirectory string
	flagBackend, flagCloud, flagGet, flagUpgrade bool
	flagPluginPath                               FlagStringSlice
	backendFlagSet                               bool
	cloudFlagSet                                 bool
	flagConfigExtra                              rawFlags
}

type initActs struct {
	*initCfg
}

func (c *initActs) initBackend(ctx context.Context, extraConfig rawFlags) error {
	return nil
}
