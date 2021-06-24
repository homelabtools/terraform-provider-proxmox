package fixtures

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type BaseFixture struct {
	T       *testing.T
	Assert  *assert.Assertions
	Require *require.Assertions
}

func NewBaseFixture(t *testing.T) BaseFixture {
	return BaseFixture{
		T:       t,
		Assert:  assert.New(t),
		Require: require.New(t),
	}
}
