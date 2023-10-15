package config

import (
	"fmt"
	"os"
)

func ResolveConfig(class int) (*Configurer, error) {
	var cfg Configurer
	switch class {
	case FLAG:
		cfg = NewEnvConfig()
		eCfg := cfg.(*EnvConfig)
		set := CommandLine
		set.StringVar(&eCfg.logLevel, FlagLogLevel, DefaultFLagLogLevel, "-")
		set.StringVar(&eCfg.port, FlagPort, DefaultFlagPort, "-")
		set.StringVar(&eCfg.enabledDB, FlagDB, DefaultFlagDB, "-")
		set.StringVar(&eCfg.dBHost, FlagDBHost, DefaultFlagDBHost, "-")
		set.StringVar(&eCfg.cloud.provider, FlagProvider, DefaultFlagProvider, "-")
		set.StringVar(&eCfg.cloud.loc, FlagLOC, DefaultFlagLOC, "-")
		err := set.Parse(os.Args[1:])
		if err != nil {
			return nil, fmt.Errorf("unable to parse arguments: %w", err)
		}
	}
	return &cfg, nil
}
