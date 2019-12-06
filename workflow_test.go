package workflow_test

import (
	"fmt"

	"github.com/gol4ng/workflow"
)

type MyObject struct {
	*workflow.Marking
}

func (m *MyObject) GetMarking() *workflow.Marking {
	return m.Marking
}

func (m *MyObject) Hello() {
	fmt.Println("Hello from MyObject")
}

func ExampleMergingWorkflow() {
	w := workflow.MustNewWorkflow(
		"start",
		map[string]workflow.Transition{
			"submit": {
				[]workflow.Place{"start"},
				[]workflow.Place{"test"},
			},
			"update": {
				[]workflow.Place{"coding", "test", "review"},
				[]workflow.Place{"test"},
			},
			"wait_for_review": {
				[]workflow.Place{"test"},
				[]workflow.Place{"review"},
			},
			"request_change": {
				[]workflow.Place{"review"},
				[]workflow.Place{"coding"},
			},
			"accept": {
				[]workflow.Place{"review"},
				[]workflow.Place{"merged"},
			},
			"reject": {
				[]workflow.Place{"review"},
				[]workflow.Place{"closed"},
			},
			"reopen": {
				[]workflow.Place{"closed"},
				[]workflow.Place{"review"},
			},
		},
		[]workflow.Place{
			"start",
			"coding",
			"test",
			"review",
			"merged",
			"closed",
		},
		map[string]workflow.Callback{},
	)

	o := &MyObject{workflow.NewMarking(w.InitialMarking)}
	fmt.Println(o.Marking.GetPlaces())

	result := w.Can(o, "update")
	fmt.Println(result)

	result = w.Can(o, "submit")
	fmt.Println(result)

	result = w.Apply(o, "submit")
	fmt.Println(result)
	fmt.Println(o.Marking.GetPlaces())
	// Output: [start]
	//false
	//true
	//true
	//[test]

}

func ExamplePublishingWorkflow() {
	w := workflow.MustNewWorkflow(
		"draft",
		map[string]workflow.Transition{
			"request_review": {
				[]workflow.Place{"draft"},
				[]workflow.Place{"waiting_for_journalist", "waiting_for_spellchecker"},
			},
			"journalist_approval": {
				[]workflow.Place{"waiting_for_journalist"},
				[]workflow.Place{"approved_by_journalist"},
			},
			"spellchecker_approval": {
				[]workflow.Place{"waiting_for_spellchecker"},
				[]workflow.Place{"approved_by_spellchecker"},
			},
			"publish": {
				[]workflow.Place{"approved_by_journalist", "approved_by_spellchecker"},
				[]workflow.Place{"published"},
			},
		},
		[]workflow.Place{
			"draft",
			"waiting_for_journalist",
			"waiting_for_spellchecker",
			"approved_by_journalist",
			"approved_by_spellchecker",
			"published",
		},
		map[string]workflow.Callback{
			"transition.request_review": func(markingStorer workflow.MarkingStorer) {
				fmt.Println("entered.request_review", markingStorer.GetMarking().GetPlaces())
				markingStorer.(*MyObject).Hello()
			},
			"transition.journalist_approval": func(markingStorer workflow.MarkingStorer) {
				fmt.Println("entered.journalist_approval", markingStorer.GetMarking().GetPlaces())
				markingStorer.(*MyObject).Hello()
			},
			"transition.spellchecker_approval": func(markingStorer workflow.MarkingStorer) {
				fmt.Println("entered.spellchecker_approval", markingStorer.GetMarking().GetPlaces())
				markingStorer.(*MyObject).Hello()
			},
			"transition.publish": func(markingStorer workflow.MarkingStorer) {
				fmt.Println("entered.publish", markingStorer.GetMarking().GetPlaces())
				markingStorer.(*MyObject).Hello()
			},
		},
	)

	o := &MyObject{workflow.NewMarking(w.InitialMarking)}
	fmt.Println(o.Marking.GetPlaces())

	w.Apply(o, "request_review")
	fmt.Println(o.Marking.GetPlaces())

	w.Apply(o, "journalist_approval")
	fmt.Println(o.Marking.GetPlaces())

	w.Apply(o, "spellchecker_approval")
	fmt.Println(o.Marking.GetPlaces())

	w.Apply(o, "publish")
	fmt.Println(o.Marking.GetPlaces())
	// Output: [draft]
	//entered.request_review []
	//Hello from MyObject
	//[waiting_for_journalist waiting_for_spellchecker]
	//entered.journalist_approval [waiting_for_spellchecker]
	//Hello from MyObject
	//[waiting_for_spellchecker approved_by_journalist]
	//entered.spellchecker_approval [approved_by_journalist]
	//Hello from MyObject
	//[approved_by_journalist approved_by_spellchecker]
	//entered.publish []
	//Hello from MyObject
	//[published]
}
