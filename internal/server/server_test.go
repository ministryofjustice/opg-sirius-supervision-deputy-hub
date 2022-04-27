package server

import (
	"io"

	"github.com/stretchr/testify/mock"
)

type mockTemplates struct {
	mock.Mock
	count    int
	lastName string
	lastVars interface{}
}

func (m *mockTemplates) ExecuteTemplate(w io.Writer, name string, vars interface{}) error {
	m.count += 1
	m.lastName = name
	m.lastVars = vars

	return nil
}
