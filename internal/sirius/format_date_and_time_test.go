package sirius

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFormatDateTimeStringIntoDateTime(t *testing.T) {
	unsortedData := "2020-10-18 10:11:08"
	expectedResponse, err := time.Parse("2006-01-02 15:04:05", "2020-10-18 10:11:08")
	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, FormatDateTimeStringIntoDateTime("2006-01-02 15:04:05", unsortedData))
}
