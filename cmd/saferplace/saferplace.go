package main

import (
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
	configFile := flag.String("config", "", "Config file")
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
	return saferplace.Run(components, cfg)
}
