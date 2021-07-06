package fixtures

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	f := &TerraformTestFixture{
		BaseFixture:      NewBaseFixture(t),
		Name:             name,
		Directory:        dir,
		TerraformVersion: tfVersion,
	}
	f.writeMainTF("https://127.0.0.1:8006")
	return f
}

func (f *TerraformTestFixture) WriteFile(filename, contents string) {
	err := ioutil.WriteFile(filepath.Join(f.Directory, filename), []byte(contents), 0644)
	f.Require.NoErrorf(err, "expected to be able to write to file '%s'", filename)
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

func (f *TerraformTestFixture) writeMainTF(endpoint string) {
	f.WriteFile("provider.tf", fmt.Sprintf(`
provider "proxmox"{
  virtual_environment {
    endpoint = %s
    username = "root@pam"
    password = "proxmox"
    insecure = true
  }
}

terraform {
  required_version = ">=1.0.0"
  required_providers {
    proxmox = {
      source = "registry.terraform.io/danitso/proxmox"
    }
  }
}

`, endpoint))

}
