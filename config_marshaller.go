package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigReaction struct {
	Name   string   `yaml:"name"`
	Rate   string   `yaml:"rate"`
	Input  []string `yaml:"input"`
	Output []string `yaml:"output"`
}
type ConfigState struct {
	Name  string `yaml:"name"`
	Value int    `yaml:"value"`
}
type RunConfig struct {
	Until float64 `yaml:"until"`
}
type Config struct {
	Reactions []ConfigReaction `yaml:"reactions"`
	States    []ConfigState    `yaml:"states"`
	Run       RunConfig        `yaml:"run"`
}

func readConfig(filename string) (Config, error) {
	fmt.Println(filename)
	config, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, errors.New("failed to read yaml file")
	}
	fmt.Println("read", filename)
	var yml Config
	err = yaml.Unmarshal(config, &yml)
	fmt.Println("Unmarshalled")
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	fmt.Println(yml)

	return yml, nil

}
