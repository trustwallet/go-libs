package testy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsMethods(t *testing.T) {
	rt := parseTags("unit,integration")

	assert.True(t, rt.contains("unit"))
	assert.True(t, rt.contains("integration"))
	assert.False(t, rt.contains(""))
	assert.False(t, rt.contains("UNIT"))

	assert.True(t, rt.containsAll("unit"))
	assert.True(t, rt.containsAll("unit", "integration"))
	assert.False(t, rt.containsAll("unit", "integration", "something-else"))
	assert.False(t, rt.containsAll("unit", "integration", ""))

	assert.True(t, rt.containsAny("unit", "something-else"))
	assert.True(t, rt.containsAny("whatever", "unit", "something-else"))
	assert.True(t, rt.containsAny("whatever", "unit", "something-else", "integration"))
	assert.False(t, rt.containsAny("whatever", "", "something-else"))
}
