package state

import (
	"gorim/internal/state"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialState(t *testing.T) {
	assert := assert.New(t)
	s := state.Initialize()
	assert.Nil(s.ActiveProfile)
	assert.True(len(s.ModList) == 0)
}
