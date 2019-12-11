package workflow

import (
	"fmt"
)

type Callback func(subject interface{})

type Place string

type Transition struct {
	Name  string
	Froms []Place
	Tos   []Place
}

type workflow struct {
	name string

	definition    definition
	markingStorer MarkingStorer

	// todo event dispatcher ?
	callbacks map[string]Callback
}

func NewWorkflow(name string, definition definition, markingStore MarkingStorer, callbacks map[string]Callback) (*workflow, error) {
	w := &workflow{
		name,
		definition,
		markingStore,
		callbacks,
	}
	return w, nil
}

func (workflow *workflow) hasPlace(place Place) bool {
	for _, definedPlace := range workflow.definition.places {
		if definedPlace == place {
			return true
		}
	}
	return false
}

func (workflow *workflow) GetMarking(subject interface{}) (*Marking, error) {
	marking, err := workflow.markingStorer.GetMarking(subject)
	if err != nil {
		return nil, err
	}
	// Init marking if empty
	if len(marking.GetPlaces()) == 0 {
		marking = NewMarking(workflow.definition.initialPlaces...)
		workflow.markingStorer.SetMarking(subject, *marking)

		//TODO workflow.entered()
	}
	return marking, nil
}

func (workflow *workflow) Can(subject interface{}, transitionName string) bool {
	workflow.GetMarking(subject)
	for _, transition := range workflow.definition.transitions {
		if transition.Name != transitionName {
			continue
		}
		// TODO transition blocker
		return true
	}
	return false
}

func (workflow *workflow) Apply(subject interface{}, transitionName string) error {
	if _, ok := workflow.definition.transitions[transitionName]; !ok {
		return fmt.Errorf("transition %s is not defined for workflow %s", transitionName, workflow.name)
	}
	marking, err := workflow.GetMarking(subject)
	if err != nil {
		return err
	}
	var approvedTransitions []Transition
	for _, transition := range workflow.definition.transitions {
		if transition.Name != transitionName {
			continue
		}
		// TODO transition blocker
		approvedTransitions = append(approvedTransitions, transition)
	}

	if len(approvedTransitions) == 0 {
		return fmt.Errorf("transition %s is not enable for workflow %s", transitionName, workflow.name)
	}

	for _, transition := range approvedTransitions {
		workflow.leave(subject, transition, marking)
		workflow.transition(subject, transition, marking)
		workflow.enter(subject, transition, marking)

		workflow.markingStorer.SetMarking(subject, *marking)

		workflow.entered(subject, transition, marking)
		workflow.completed(subject, transition, marking)
		workflow.announce(subject, transition, marking)
	}
	return nil
}

func (workflow *workflow) leave(subject interface{}, transition Transition, marking *Marking) {
	workflow.callCallback("workflow.leave", subject)
	workflow.callCallback("workflow."+workflow.name+".leave", subject)

	places := transition.Froms
	for _, place := range places {
		workflow.callCallback("workflow."+workflow.name+".leave."+string(place), subject)
	}

	for _, place := range places {
		marking.Unmark(place)
	}
}

func (workflow *workflow) transition(subject interface{}, transition Transition, marking *Marking) {
	workflow.callCallback("workflow.transition", subject)
	workflow.callCallback("workflow."+workflow.name+".transition", subject)
	workflow.callCallback("workflow."+workflow.name+".transition."+transition.Name, subject)
}

func (workflow *workflow) enter(subject interface{}, transition Transition, marking *Marking) {
	workflow.callCallback("workflow.enter", subject)
	workflow.callCallback("workflow."+workflow.name+".enter", subject)
	places := transition.Tos
	for _, place := range places {
		workflow.callCallback("workflow."+workflow.name+".enter."+string(place), subject)
	}

	for _, place := range places {
		marking.Mark(place)
	}
}

func (workflow *workflow) entered(subject interface{}, transition Transition, marking *Marking) {
	workflow.callCallback("workflow.entered", subject)
	workflow.callCallback("workflow."+workflow.name+".entered", subject)
	for _, place := range transition.Tos {
		workflow.callCallback("workflow."+workflow.name+".entered."+string(place), subject)
	}
}

func (workflow *workflow) completed(subject interface{}, transition Transition, marking *Marking) {
	workflow.callCallback("workflow.completed", subject)
	workflow.callCallback("workflow."+workflow.name+".completed", subject)
	workflow.callCallback("workflow."+workflow.name+".completed."+transition.Name, subject)
}

func (workflow *workflow) announce(subject interface{}, transition Transition, marking *Marking) {
	workflow.callCallback("workflow.announce", subject)
	workflow.callCallback("workflow."+workflow.name+".announce", subject)
	for transitionName, _ := range workflow.getEnabledTransitions(subject) {
		workflow.callCallback("workflow."+workflow.name+".announce."+transitionName, subject)
	}
}

func (workflow *workflow) callCallback(key string, subject interface{}) {
	if callback, ok := workflow.callbacks[key]; ok {
		callback(subject)
	}
}

func (workflow *workflow) hasTransition(transitionName string) bool {
	_, hasTransition := workflow.definition.transitions[transitionName]
	return hasTransition
}

func (workflow *workflow) getTransition(transitionName string) Transition {
	return workflow.definition.transitions[transitionName]
}

func (workflow *workflow) getEnabledTransitions(subject interface{}) map[string]Transition {
	workflow.GetMarking(subject)
	enabledTransitions := map[string]Transition{}
	for transitionName, transition := range workflow.definition.transitions {
		if workflow.Can(subject, transitionName) {
			enabledTransitions[transitionName] = transition
		}
	}

	return enabledTransitions
}
