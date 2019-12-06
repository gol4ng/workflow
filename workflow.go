package workflow

import (
	"fmt"
)

type Callback func(markingStorer MarkingStorer)

type Place string

type Transition struct {
	Froms []Place
	Tos []Place
}

type workflow struct {
	InitialMarking Place

	transitions map[string]Transition
	places []Place

	callbacks map[string]Callback
}

func NewWorkflow(initialMarking Place, transitions map[string]Transition, places []Place, callbacks map[string]Callback) (*workflow, error) {
	w := &workflow{
		initialMarking,
		transitions,
		places,
		callbacks,
	}

	return w, w.validate()
}

func MustNewWorkflow(initialMarking Place, transitions map[string]Transition, places []Place, callbacks map[string]Callback) *workflow {
	w, err := NewWorkflow(initialMarking, transitions, places, callbacks)
	if err != nil {
		panic(err)
	}

	return w
}

func (w *workflow) validate() error {
	for transitionName, transition := range w.transitions{
		for _, from := range transition.Froms {
			if !w.hasPlace(from) {
				return fmt.Errorf("place %s not declared in workflow on transition %s in from", from, transitionName)
			}
		}

		for _, to := range transition.Tos {
			if !w.hasPlace(to) {
				return fmt.Errorf("place %s not declared in workflow on transition %s in to", to, transitionName)
			}
		}
	}

	return nil
}

func (w *workflow) hasPlace(place Place) bool {
	for _, definedPlace := range w.places {
		if definedPlace == place {
			return true
		}
	}

	return false
}

func (w *workflow) Can(markingStorer MarkingStorer, transitionName string) bool {
	currentPlaces := markingStorer.GetMarking().GetPlaces()
	for _, currentPlace := range currentPlaces{
		for _, from := range w.getTransition(transitionName).Froms {
			if from == currentPlace {
				return true
			}
		}
	}

	return false
}

func (w *workflow) Apply(markingStorer MarkingStorer, transitionName string) bool {
	if !w.Can(markingStorer, transitionName) {
		return false
	}

	transition := w.getTransition(transitionName)
	w.leave(markingStorer, transition)
	w.transition(markingStorer, transitionName)
	w.enter(markingStorer, transition)

	for _, to := range transition.Tos {
		markingStorer.GetMarking().Mark(to)
	}

	w.entered(markingStorer, transition)
	w.completed(markingStorer, transitionName)
	w.announce(markingStorer, transition)

	return true
}

func (w *workflow) leave(markingStorer MarkingStorer, transition Transition) {
	w.callCallback("leave", markingStorer)

	places := transition.Froms
	for _, place := range places {
		w.callCallback("leave."+string(place), markingStorer)
	}

	for _, place := range places {
		markingStorer.GetMarking().Unmark(place)
	}
}

func (w *workflow) transition(markingStorer MarkingStorer, transitionName string) {
	w.callCallback("transition", markingStorer)
	w.callCallback("transition."+transitionName, markingStorer)
}

func (w *workflow) enter(markingStorer MarkingStorer, transition Transition) {
	w.callCallback("enter", markingStorer)
	for _, place := range transition.Tos {
		w.callCallback("enter."+string(place), markingStorer)
	}
}

func (w *workflow) entered(markingStorer MarkingStorer, transition Transition) {
	w.callCallback("entered", markingStorer)
	for _, place := range transition.Tos {
		w.callCallback("entered."+string(place), markingStorer)
	}
}

func (w *workflow) completed(markingStorer MarkingStorer, transitionName string) {
	w.callCallback("completed", markingStorer)
	w.callCallback("completed."+transitionName, markingStorer)
}

func (w *workflow) announce(markingStorer MarkingStorer, transition Transition) {
	w.callCallback("announce", markingStorer)
	for transitionName, _ := range w.getEnabledTransitions(markingStorer) {
		w.callCallback("announce."+transitionName, markingStorer)
	}
}

func (w *workflow) callCallback(key string, markingStorer MarkingStorer) {
	if callback, ok := w.callbacks[key]; ok {
		callback(markingStorer)
	}
}

func (w *workflow) getTransition(transitionName string) Transition {
	return w.transitions[transitionName]
}

func (w *workflow) getEnabledTransitions(markingStorer MarkingStorer) map[string]Transition {
	enabledTransitions := map[string]Transition {}
	for transitionName, transition := range w.transitions {
		if w.Can(markingStorer, transitionName) {
			enabledTransitions[transitionName] = transition
		}
	}

	return enabledTransitions
}
