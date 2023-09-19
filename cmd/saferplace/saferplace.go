package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"safer.place/internal/cmd/saferplace"
	"safer.place/internal/config"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() error {
	configFile := flag.String("config", "/etc/saferplace/config.yaml", "Config file")
	flag.Parse()

	components := saferplace.AllComponents()
	if len(flag.Args()) > 0 {
		if flag.Arg(0) != "all" {
			components = saferplace.StringsToComponents(strings.Split(flag.Arg(0), ","))
		}
	}

	cfg, err := config.Parse(*configFile)
	if err != nil {
		return err
	}
	return saferplace.Run(context.Background(), components, cfg)
}
