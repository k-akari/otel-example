package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type env struct {
}

func mustNewConfig() *env {
	var e env
	if err := envconfig.Process("", &e); err != nil {
		panic(fmt.Errorf("failed to process envconfig: %w", err))
	}

	return &e
}
