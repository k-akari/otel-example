package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type env struct {
	Port           int    `envconfig:"PORT" required:"true"`
	EndpointJaeger string `envconfig:"ENDPOINT_JAEGER" required:"true"`
}

func mustNewConfig() *env {
	var e env
	if err := envconfig.Process("", &e); err != nil {
		panic(fmt.Errorf("failed to process envconfig: %w", err))
	}

	return &e
}
