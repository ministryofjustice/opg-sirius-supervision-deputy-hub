package sirius

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatDateAndTime(t *testing.T) {
	unsortedData := "2020-10-18 10:11:08"
	expectedResponse := "18/10/2020 10:11:08"
	assert.Equal(t, expectedResponse, FormatDateAndTime("2006-01-02 15:04:05", unsortedData, "02/01/2006 15:04:05"))
}
