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
