package workflow

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
