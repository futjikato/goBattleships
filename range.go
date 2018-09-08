package main

import (
	"fmt"
	"strconv"
)

type RangePromter struct {
	id       string
	next     string
	question string
	min      int
	max      int
	field    *int
}

func (cp *RangePromter) ID() string {
	return cp.id
}

func (cp *RangePromter) PromptString() string {
	return cp.question
}

func (cp *RangePromter) NextOnError(e error) string {
	return cp.ID()
}

func (cp *RangePromter) NextOnSuccess(choice string) string {
	*cp.field, _ = strconv.Atoi(choice)
	return cp.next
}

func (cp *RangePromter) Parse(choice string) error {
	p, err := strconv.Atoi(choice)
	if err != nil {
		return err
	}
	if p < cp.min {
		return fmt.Errorf("Invalid port %d. Must be at least %d", p, cp.min)
	}
	if p > cp.max {
		return fmt.Errorf("Invalid port %d. Must can not be higher then %d", p, cp.max)
	}

	return nil
}

func NewRangePromter(ID, question, next string, min, max int, field *int) *RangePromter {
	return &RangePromter{
		id:       ID,
		next:     next,
		question: question,
		min:      min,
		max:      max,
		field:    field,
	}
}
