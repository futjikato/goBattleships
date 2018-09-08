package main

import (
	"fmt"
)

type ChoicePromter struct {
	id       string
	next     map[string]string
	question string
	choices  map[string]string
	field    *string
}

func (cp *ChoicePromter) ID() string {
	return cp.id
}

func (cp *ChoicePromter) PromptString() string {
	return cp.question
}

func (cp *ChoicePromter) NextOnError(e error) string {
	return cp.ID()
}

func (cp *ChoicePromter) NextOnSuccess(choice string) string {
	*cp.field = cp.choices[choice]
	return cp.next[choice]
}

func (cp *ChoicePromter) Parse(choice string) error {
	if _, ok := cp.choices[choice]; ok == false {
		possibleInput := make([]string, 0)
		for k := range cp.choices {
			possibleInput = append(possibleInput, k)
		}
		return fmt.Errorf("Invalid input %s. You need to enter one of %v", choice, possibleInput)
	}

	return nil
}

func NewChoice(ID, question string, next, choices map[string]string, field *string) *ChoicePromter {
	return &ChoicePromter{
		id:       ID,
		next:     next,
		question: question,
		choices:  choices,
		field:    field,
	}
}
