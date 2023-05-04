// Copyright 2022 SaferPlace

package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"

	cmd "safer.place/realtime/internal/cmd/realtime"
)

func main() {
	if err := run(); err != nil {
		// There must be a more elegant way, but that's for later.
		panic(err)
	}
}

func run() error {
	var cfg cmd.Config
	if err := envconfig.Process("saferplace", &cfg); err != nil {
		return fmt.Errorf("unable to parse config: %w", err)
	}
	return cmd.Run(cfg)
}
