package workflow_test

import (
	"fmt"

	"github.com/gol4ng/workflow"
)

type MyObject struct {
	*workflow.Marking
}

func (m *MyObject) GetMarking() (workflow.Marking, error) {
	if m.Marking == nil {
		m.Marking = workflow.NewMarking()
	}
	return *m.Marking, nil
}

func (m *MyObject) SetMarking(marking workflow.Marking) error {
	for _, place := range marking.GetPlaces() {
		m.Marking.Mark(place)
	}
	return nil
}

func (m *MyObject) Hello() {
	fmt.Println("Hello from MyObject")
}

func ExampleMergingWorkflow() {
	w, _ := workflow.NewWorkflow(
		"merging",
		workflow.MustNewDefinition(
			[]workflow.Place{
				"start",
				"coding",
				"test",
				"review",
				"merged",
				"closed",
			},
			[]workflow.Place{"start"},
			[]workflow.Transition{
				{
					"submit",
					[]workflow.Place{"start"},
					[]workflow.Place{"test"},
				},
				{
					"update",
					[]workflow.Place{"coding", "test", "review"},
					[]workflow.Place{"test"},
				},
				{
					"wait_for_review",
					[]workflow.Place{"test"},
					[]workflow.Place{"review"},
				},
				{
					"request_change",
					[]workflow.Place{"review"},
					[]workflow.Place{"coding"},
				},
				{
					"accept",
					[]workflow.Place{"review"},
					[]workflow.Place{"merged"},
				},
				{
					"reject",
					[]workflow.Place{"review"},
					[]workflow.Place{"closed"},
				},
				{
					"reopen",
					[]workflow.Place{"closed"},
					[]workflow.Place{"review"},
				},
			},
		),
		&workflow.MethodMarkingStorer{},
		map[string]workflow.Callback{},
	)

	o := &MyObject{}
	fmt.Println(w.GetMarking(o))
	fmt.Println("Can(update)", w.Can(o, "update"))
	fmt.Println("Can(submit)", w.Can(o, "submit"))
	fmt.Println("Apply(submit)", w.Apply(o, "submit"))
	fmt.Println(w.GetMarking(o))

	// Output:
	//&{map[start:start]} <nil>
	//Can(update) true
	//Can(submit) true
	//Apply(submit) <nil>
	//&{map[test:test]} <nil>
}

//                          ┏━(waiting_for_journalist)━━━>[journalist_approval]━━>(approved_by_journalist)━━━━┓
// (main)━━>[request_review]┫                                                                                 ┣━>[publish]━>(publish)
//                          ┗━(waiting_for_spellchecker)━>[spellchecker_approval]━>(approved_by_spellchecker)━┛
func ExamplePublishingWorkflow() {
	w, _ := workflow.NewWorkflow(
		"review",
		workflow.MustNewDefinition(
			[]workflow.Place{
				"draft",
				"waiting_for_journalist",
				"waiting_for_spellchecker",
				"approved_by_journalist",
				"approved_by_spellchecker",
				"published",
			},
			[]workflow.Place{"draft"},
			[]workflow.Transition{
				{
					"request_review",
					[]workflow.Place{"draft"},
					[]workflow.Place{"waiting_for_journalist", "waiting_for_spellchecker"},
				},
				{
					"journalist_approval",
					[]workflow.Place{"waiting_for_journalist"},
					[]workflow.Place{"approved_by_journalist"},
				},
				{
					"spellchecker_approval",
					[]workflow.Place{"waiting_for_spellchecker"},
					[]workflow.Place{"approved_by_spellchecker"},
				},
				{
					"publish",
					[]workflow.Place{"approved_by_journalist", "approved_by_spellchecker"},
					[]workflow.Place{"published"},
				},
			},
		),
		&workflow.MethodMarkingStorer{},
		map[string]workflow.Callback{
			"workflow.review.transition.request_review": func(subject interface{}) {
				//fmt.Println("entered.request_review", markingStorer.GetMarking().GetPlaces())
				subject.(*MyObject).Hello()
			},
			"workflow.review.transition.journalist_approval": func(subject interface{}) {
				//fmt.Println("entered.journalist_approval", subject.GetMarking().GetPlaces())
				subject.(*MyObject).Hello()
			},
			"workflow.review.transition.spellchecker_approval": func(subject interface{}) {
				//fmt.Println("entered.spellchecker_approval", subject.GetMarking().GetPlaces())
				subject.(*MyObject).Hello()
			},
			"workflow.review.transition.publish": func(subject interface{}) {
				//fmt.Println("entered.publish", subject.GetMarking().GetPlaces())
				subject.(*MyObject).Hello()
			},
		},
	)

	o := &MyObject{}

	fmt.Println(w.GetMarking(o))
	fmt.Println("apply(request_review)", w.Apply(o, "request_review"))
	fmt.Println(w.GetMarking(o))
	fmt.Println("apply(journalist_approval)", w.Apply(o, "journalist_approval"))
	fmt.Println(w.GetMarking(o))
	fmt.Println("apply(spellchecker_approval)", w.Apply(o, "spellchecker_approval"))
	fmt.Println(w.GetMarking(o))
	fmt.Println("apply(publish)", w.Apply(o, "publish"))
	fmt.Println(w.GetMarking(o))

	// Output:
	//&{map[draft:draft]} <nil>
	//Hello from MyObject
	//apply(request_review) <nil>
	//&{map[waiting_for_journalist:waiting_for_journalist waiting_for_spellchecker:waiting_for_spellchecker]} <nil>
	//Hello from MyObject
	//apply(journalist_approval) <nil>
	//&{map[approved_by_journalist:approved_by_journalist waiting_for_spellchecker:waiting_for_spellchecker]} <nil>
	//Hello from MyObject
	//apply(spellchecker_approval) <nil>
	//&{map[approved_by_journalist:approved_by_journalist approved_by_spellchecker:approved_by_spellchecker]} <nil>
	//Hello from MyObject
	//apply(publish) <nil>
	//&{map[published:published]} <nil>
}
