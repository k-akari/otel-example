package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type env struct {
	Port           int    `envconfig:"PORT" required:"true"`
	EndpointJaeger string `envconfig:"ENDPOINT_JAEGER" required:"true"`
	DBUser         string `envconfig:"DB_USER" required:"true"`
	DBPass         string `envconfig:"DB_PASS" required:"true"`
	DBHost         string `envconfig:"DB_HOST" required:"true"`
	DBName         string `envconfig:"DB_NAME" required:"true"`
	DBPort         int    `envconfig:"DB_PORT" required:"true"`
}

func mustNewConfig() *env {
	var e env
	if err := envconfig.Process("", &e); err != nil {
		panic(fmt.Errorf("failed to process envconfig: %w", err))
	}

	return &e
}
