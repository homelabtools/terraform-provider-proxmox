package fixtures

import (
	"os"
	"reflect"
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

func (f *BaseFixture) ShouldClean(fixture interface{}) bool {
	if os.Getenv("SKIP_CLEANUP") != "" {
		if fixture != nil {
			f.T.Logf("SKIP_CLEANUP env var found, skipping cleanup of %s", reflect.TypeOf(fixture))
		}
		return false
	}
	return true
}
