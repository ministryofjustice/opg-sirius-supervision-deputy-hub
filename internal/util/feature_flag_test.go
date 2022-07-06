package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsFeatureFlagged(t *testing.T) {
	features := []string{"flag", "banner"}
	featureFlagged := IsFeatureFlagged(features)

	assert.True(t, featureFlagged("flag"))
	assert.False(t, featureFlagged("flog"))
}
