// Copyright 2022 SaferPlace

package main

import (
	cmd "safer.place/internal/cmd/realtime"
	"safer.place/internal/config"
)

func main() {
	if err := run(); err != nil {
		// There must be a more elegant way, but that's for later.
		panic(err)
	}
}

func run() error {
	cfg, err := config.Parse("saferplace")
	if err != nil {
		return err
	}
	return cmd.Run(cfg)
}
