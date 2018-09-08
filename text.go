package main

type TextPromter struct {
	id       string
	question string
	next     string
	field    *string
}

func (cp *TextPromter) ID() string {
	return cp.id
}

func (cp *TextPromter) PromptString() string {
	return cp.question
}

func (cp *TextPromter) NextOnError(e error) string {
	return cp.ID()
}

func (cp *TextPromter) NextOnSuccess(choice string) string {
	*cp.field = choice
	return cp.next
}

func (cp *TextPromter) Parse(choice string) error {
	return nil
}

func NewTextPromter(ID, question, next string, field *string) *TextPromter {
	return &TextPromter{
		id:       ID,
		question: question,
		next:     next,
		field:    field,
	}
}
