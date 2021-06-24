package fixtures

import (
	"io/ioutil"
	"os"
	"testing"
)

type TerraformTestFixture struct {
	BaseFixture
	Directory        string
	Name             string
	TerraformVersion string
}

// NewTerraformTestFixture creates a new test fixture for testing Terraform providers.
// It creates a temporary directory where .tf files will be written.
//
// tfVersion - version of TF to use
// name      - an optional string for extra description.
func NewTerraformTestFixture(t *testing.T, name, tfVersion string) *TerraformTestFixture {
	dir, err := ioutil.TempDir("", name)
	if err != nil {
		t.Fatal("Unable to create temp directory for Terraform test fixture")
	}
	t.Logf("Created TF test fixture named '%s' at '%s', TF version '%s'", name, dir, tfVersion)
	return &TerraformTestFixture{
		BaseFixture:      NewBaseFixture(t),
		Name:             name,
		Directory:        dir,
		TerraformVersion: tfVersion,
	}
}

// TearDown removes every trace the test fixture.
// It should be called with defer right after creating the fixture.
func (f *TerraformTestFixture) TearDown() {
	if !f.ShouldClean(f) {
		return
	}
	if os.RemoveAll(f.Directory) != nil {
		// Fatal is used here because this should never really happen,
		// and if it did, may indicate something is very wrong.
		f.T.Fatal("Unable to clean up temp directory")
	}
}
