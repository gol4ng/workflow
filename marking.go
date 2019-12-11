package workflow

import (
	"fmt"
)

type MarkingStorer interface {
	GetMarking(subject interface{}) (*Marking, error)
	SetMarking(subject interface{}, marking Marking) error
}

type MethodMarkingStorer struct {
}

func (m *MethodMarkingStorer) GetMarking(subject interface{}) (*Marking, error) {
	if v, ok := subject.(markedSubject); ok {
		marking, err := v.GetMarking()
		return &marking, err
	}
	return nil, fmt.Errorf("marking not found")
}

func (m *MethodMarkingStorer) SetMarking(subject interface{}, marking Marking) error {
	if v, ok := subject.(markedSubject); ok {
		return v.SetMarking(marking)
	}
	return fmt.Errorf("marking not found")
}

type markedSubject interface {
	GetMarking() (Marking, error)
	SetMarking(marking Marking) error
}

type Marking struct {
	places map[Place]Place
}

func NewMarking(representation ...Place) *Marking {
	marking := &Marking{
		places: map[Place]Place{},
	}
	for _, place := range representation {
		marking.Mark(place)
	}

	return marking
}

func (m *Marking) Mark(place Place) {
	if _, ok := m.places[place]; !ok {
		m.places[place] = place
	}
}

func (m *Marking) Unmark(p Place) {
	if _, ok := m.places[p]; ok {
		delete(m.places, p)
	}
}

func (m *Marking) Has(place Place) bool {
	_, has := m.places[place]
	return has
}

func (m *Marking) GetPlaces() []Place {
	var places []Place
	for _, place := range m.places {
		places = append(places, place)
	}
	return places
}
