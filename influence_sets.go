package main

type Set struct {
	list map[string]struct{}
}

func (s *Set) Has(v string) bool {
	_, ok := s.list[v]
	return ok
}
func (s *Set) Add(v string) {
	s.list[v] = struct{}{}
}
func (s *Set) Size() int {
	return len(s.list)
}
func NewSet() *Set {
	s := &Set{}
	s.list = make(map[string]struct{})
	return s
}

func MakeInfluenceSets(reactions map[string]Reaction) map[string][]string {
	dependentReactions := make(map[string]*Set)
	changedSpecies := make(map[string]*Set)

	for rxnName, rxnInfo := range reactions {
		if _, exists := changedSpecies[rxnName]; !exists {
			changedSpecies[rxnName] = NewSet()
		}
		for _, species := range rxnInfo.inputSpecies {
			if _, exists := dependentReactions[species]; !exists {
				dependentReactions[species] = NewSet()
			}
			dependentReactions[species].Add(rxnName)

			changedSpecies[rxnName].Add(species)
		}
		for _, species := range rxnInfo.outputSpecies {
			changedSpecies[rxnName].Add(species)
		}
		for _, species := range rxnInfo.rateSpecies {
			if _, exists := dependentReactions[species]; !exists {
				dependentReactions[species] = NewSet()
			}
			dependentReactions[species].Add(rxnName)
		}
	}

	InfluenceSet := make(map[string][]string, len(reactions))

	for rxnName, speciesArr := range changedSpecies {
		influence := NewSet()
		for species := range speciesArr.list {
			for reaction := range dependentReactions[species].list {
				influence.Add(reaction)
			}
		}
		for influencedReaction := range influence.list {
			InfluenceSet[rxnName] = append(InfluenceSet[rxnName], influencedReaction)
		}
	}

	return InfluenceSet
}
