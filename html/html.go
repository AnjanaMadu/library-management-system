package html

import (
	"strings"
)

type HTMLElement struct {
	Tag        string
	Content    string
	Class      []string
	Attributes map[string]string
}

func (e *HTMLElement) getString() string {
	s := "<" + e.Tag
	if len(e.Class) > 0 {
		s += " class=\"" + strings.Join(e.Class, " ") + "\""
	}
	for k, v := range e.Attributes {
		s += " " + k + "=\"" + v + "\""
	}
	s += ">" + e.Content + "</" + e.Tag + ">"
	return s
}

func (e *HTMLElement) addClass(class string) {
	e.Class = append(e.Class, class)
}

func (e *HTMLElement) setAttribute(key, value string) {
	e.Attributes[key] = value
}

func (e *HTMLElement) appendHTML(html string) {
	e.Content += html
}

func createElement(tag string) *HTMLElement {
	return &HTMLElement{
		Tag:        tag,
		Class:      make([]string, 0),
		Attributes: make(map[string]string),
	}
}
