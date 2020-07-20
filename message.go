package main

import "fmt"

type MessageInterface interface {
	ID() string
	Name() string
	Data() string
}

type DataMessage string

func (s DataMessage) ID() string {
	return ""
}

func (s DataMessage) Name() string {
	return ""
}

func (s DataMessage) Data() string {
	return string(s)
}

func EncodeMessage(m MessageInterface) (s string) {
	s = ""
	if m.ID() != "" {
		s += fmt.Sprintf("id:%s\n", m.ID())
	}

	if m.Name() != "" {
		s += fmt.Sprintf("name:%s\n", m.Name())
	}

	if m.Data() != "" {
		s += fmt.Sprintf("data:%s\n", m.Data())
	}

	s += "\n"

	return
}
