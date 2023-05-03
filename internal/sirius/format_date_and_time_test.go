package sirius

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatDateAndTime(t *testing.T) {
	unsortedData := "2020-10-18 10:11:08"
	expectedResponse := ""

	if isDST() {
		expectedResponse = "18/10/2020 11:11:08"
	} else {
		expectedResponse = "18/10/2020 10:11:08"
	}

	assert.Equal(t, expectedResponse, FormatDateAndTime(TimelineDateTimeFormat, unsortedData, TimelineDateTimeDisplayFormat))
}
