package main

import (
	"github.com/ddouglas/todoer/pkg/types"
	"github.com/kelseyhightower/envconfig"
)

var c *types.Config = new(types.Config)

func LoadConfig() error {
	return envconfig.Process("TODOER", c)
}
