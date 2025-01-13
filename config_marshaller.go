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

func parseConfig(rawConfig Config) (map[string]Reaction, map[string]int, float64) {

	runTime := rawConfig.Run.Until

	reactionInfo := make(map[string]Reaction)
	initialState := make(map[string]int)

	for _, state := range rawConfig.States {
		initialState[state.Name] = state.Value
	}

	for _, reaction := range rawConfig.Reactions {
		_ = reaction
		var inputCoefs []int
		var inputSpecies []string
		var outputCoefs []int
		var outputSpecies []string
		rateCoef, rateSpecies := multiplyStringParse(reaction.Rate)
		for _, inputInfo := range reaction.Input {
			coef, species := multiplyStringParse(inputInfo)
			inputCoefs = append(inputCoefs, int(coef))
			inputSpecies = append(inputSpecies, species...)
		}

		for _, outputInfo := range reaction.Output {
			coef, species := multiplyStringParse(outputInfo)
			outputCoefs = append(outputCoefs, int(coef))
			outputSpecies = append(outputSpecies, species...)
		}
		reactionInfo[reaction.Name] = Reaction{
			rate:               rateCoef,
			rateSpecies:        rateSpecies,
			inputSpeciesCount:  inputCoefs,
			inputSpecies:       inputSpecies,
			outputSpeciesCount: outputCoefs,
			outputSpecies:      outputSpecies,
		}

	}

	return reactionInfo, initialState, runTime
}
