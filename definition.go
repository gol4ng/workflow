package workflow

import (
	"fmt"
)

type definition struct {
	places        map[Place]Place
	initialPlaces []Place
	transitions   map[string]Transition
}

func NewDefinition(places []Place, initialPlace []Place, transitions []Transition) (definition, error) {
	definition := &definition{
		places:      map[Place]Place{},
		transitions: map[string]Transition{},
	}
	definition.addPlaces(places)
	if err := definition.addTransitions(transitions); err != nil {
		return *definition, err
	}
	if err := definition.setInitialPlaces(initialPlace); err != nil {
		return *definition, err
	}
	definition.initialPlaces = initialPlace
	return *definition, nil
}

func MustNewDefinition(places []Place, initialPlaces []Place, transitions []Transition) definition {
	definition, err := NewDefinition(places, initialPlaces, transitions)
	if err != nil {
		panic(err)
	}
	return definition
}

func (definition *definition) setInitialPlaces(places []Place) error {
	for _, place := range places {
		if _, ok := definition.places[place]; !ok {
			return fmt.Errorf("place %s cannot be the initial place as it does not exist", place)
		}
	}
	return nil
}

func (definition *definition) addPlaces(places []Place) {
	for _, place := range places {
		definition.addPlace(place)
	}
}

func (definition *definition) addPlace(place Place) {
	if _, ok := definition.places[place]; !ok {
		definition.places[place] = place
	}
}

func (definition *definition) addTransitions(transitions []Transition) error {
	for _, transition := range transitions {
		if err := definition.addTransition(transition); err != nil {
			return err
		}
	}
	return nil
}

func (definition *definition) addTransition(transition Transition) error {
	transitionName := transition.Name
	for _, place := range transition.Froms {
		if _, ok := definition.places[place]; !ok {
			return fmt.Errorf("place %s referenced in transition %s does not exist", place, transitionName)
		}
	}
	for _, place := range transition.Tos {
		if _, ok := definition.places[place]; !ok {
			return fmt.Errorf("place %s referenced in transition %s does not exist", place, transitionName)
		}
	}
	definition.transitions[transitionName] = transition
	return nil
}
