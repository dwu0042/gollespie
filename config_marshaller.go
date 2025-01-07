package main

import (
	"gopkg.in/yaml.v3"
)

type ConfigReaction struct {
	Name string `yaml:"name"`
	Rate string `yaml:"rate"`
	Transition string `yaml:"transition"`
}
type ConfigState struct {
	Name string `yaml:"name"`
	Value int `yaml:"value"`
}
type Config struct {
	Reactions []ConfigReaction `yaml:"reactions"`
	States []ConfigState `yaml:"states"`
	Run []
}