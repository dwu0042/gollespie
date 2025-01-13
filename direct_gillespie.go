package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
)

type Reaction struct {
	rate               float64
	rateSpecies        []string
	inputSpecies       []string
	inputSpeciesCount  []int
	outputSpecies      []string
	outputSpeciesCount []int
}

func computeRate(rx Reaction, state map[string]int) float64 {
	rate := rx.rate
	for _, s := range rx.rateSpecies {
		rate = rate * float64(state[s])
	}
	return rate
}

type RatesMapping struct {
	rates     map[string]float64
	totalrate float64
}

func (rs *RatesMapping) insert(key string, rate float64) {
	existing := rs.rates[key]
	rs.totalrate -= existing
	rs.rates[key] = rate
	rs.totalrate += rate
}
func (rs *RatesMapping) nextTime() float64 {
	return rand.ExpFloat64() / rs.totalrate
}
func (rs *RatesMapping) drawReaction() (string, error) {
	rv := rand.Float64() * rs.totalrate
	for k, v := range rs.rates {
		rv -= v
		if rv < 0 {
			return k, nil
		}
	}
	return "", errors.New("reaction not found")
}
func NewRatesMapping(reactions map[string]Reaction, state map[string]int) *RatesMapping {
	ratesMap := make(map[string]float64)
	totalRate := 0.0
	for name, rxn := range reactions {
		rateValue := computeRate(rxn, state)
		ratesMap[name] = rateValue
		totalRate += rateValue
	}
	return &RatesMapping{ratesMap, totalRate}
}

func directMethod(reactions map[string]Reaction, maxTime float64, initialState map[string]int) {

	time := 0.0
	// init state
	state := make(map[string]int, len(initialState))
	for k, v := range initialState {
		state[k] = v
	}
	// init hazards
	hazards := NewRatesMapping(reactions, state)
	// compute influence sets of reactions
	influenceSets := MakeInfluenceSets(reactions)
	for {
		// draw time of next reaction
		time += hazards.nextTime()

		// draw next reaction
		reaction, drawError := hazards.drawReaction()
		if drawError != nil {
			fmt.Println("Reaction draw failed at time", time, ":", drawError)
			break
		}
		// realise the reaction outcome
		chosenReaction := reactions[reaction]
		for i, reagent := range chosenReaction.inputSpecies {
			cnt := chosenReaction.inputSpeciesCount[i]
			state[reagent] -= cnt
		}
		for i, reagent := range chosenReaction.outputSpecies {
			cnt := chosenReaction.outputSpeciesCount[i]
			state[reagent] += cnt
		}
		// check if simulation is done
		if time > maxTime {
			break
		}
		// update the rates of affected reactions
		for _, rxnName := range influenceSets[reaction] {
			// recompute rate and insert
			newRate := computeRate(reactions[rxnName], state)
			hazards.insert(rxnName, newRate)
		}
		fmt.Print(time, ",", state["A"], "\n")
	}

}

func product(vs []float64) float64 {
	res := 1.0
	for _, v := range vs {
		res *= v
	}
	return res
}
func multiplyStringParse(s string) (float64, []string) {
	runes := strings.Split(s, "*")
	var coefs []float64
	var species []string
	for _, substr := range runes {
		if s, err := strconv.ParseFloat(substr, 64); err == nil {
			coefs = append(coefs, s)
		} else {
			species = append(species, substr)
		}
	}
	return product(coefs), species
}

func main() {

	rawConfig, err := readConfig("reactions.yaml")
	if err != nil {
		panic(err)
	}
	_ = rawConfig
	reactionInfo, initialState, runTime := parseConfig(rawConfig)
	fmt.Println(reactionInfo)
	fmt.Println(initialState)
	// reactionInfo := map[string]Reaction{
	// "replication": {rate: 1.0, rateSpecies: []string{"A"}, inputSpecies: []string{"A"}, inputSpeciesCount: []int{1}, outputSpecies: []string{"A"}, outputSpeciesCount: []int{2}},
	// "death":       {rate: 0.8, rateSpecies: []string{"A"}, inputSpecies: []string{"A"}, inputSpeciesCount: []int{1}, outputSpecies: []string{}, outputSpeciesCount: []int{}},
	// }
	// initialState := map[string]int{
	// "A": 10,
	// }
	// runTime := 100.0
	directMethod(reactionInfo, runTime, initialState)
}
