package workflow

type MarkingStorer interface {
	GetMarking() *Marking
}

type Marking struct {
	places []Place
}

func NewMarking(representation ...Place) *Marking {
	marking := new(Marking)
	for _, place := range representation {
		marking.Mark(place)
	}

	return marking
}

func (m *Marking) Mark(place Place) {
	m.places = append(m.places, place)
}

func (m *Marking) Unmark(p Place) {
	var newPlaces []Place
	for _, place := range m.places {
		if place != p {
			newPlaces = append(newPlaces, place)
		}
	}

	m.places = newPlaces
}

func (m *Marking) Has(place Place) bool {
	for _, currentPlace := range m.places {
		if currentPlace == place {
			return true
		}
	}

	return false
}

func (m *Marking) GetPlaces() []Place {
	return m.places
}
