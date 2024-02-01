package pseudos

import (
	"errors"
	"kaiju/uimarkup/css/rules"
	"kaiju/uimarkup/markup"
)

func (p Current) Process(elm markup.DocElement, value rules.SelectorPart) ([]markup.DocElement, error) {
	return []markup.DocElement{elm}, errors.New("not implemented")
}
