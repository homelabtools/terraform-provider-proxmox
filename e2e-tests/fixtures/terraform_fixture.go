package fixtures

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

type TerraformTestFixture struct {
	BaseFixture
	Directory        string
	TerraformVersion string
	Options          *terraform.Options
}

// NewTerraformTestFixture creates a new test fixture for testing Terraform providers.
// It creates a temporary directory where .tf files will be written.
//
// tfVersion - version of TF to use
// endpoint  - URL of Proxmox instance
// username  - Proxmox username
// password  - Proxmox password
func NewTerraformTestFixture(t *testing.T, testSourceDir, tfVersion, endpoint, username, password string) *TerraformTestFixture {
	name := filepath.Base(testSourceDir)
	dir := filepath.Join("testbed", fmt.Sprintf("%s-%s", name, time.Now().Format(FileTimestampFormat)))
	require.NoError(t, os.MkdirAll(filepath.Dir(dir), 0755))

	// TODO: Don't rely on cp, use library
	cmd := exec.Command("cp", "-a", testSourceDir, dir)
	require.NoError(t, cmd.Run(), "failed to copy test case at '%s' into testbed directory at '%s'", testSourceDir, dir)

	t.Logf("Created TF test fixture at '%s', TF version '%s'", dir, tfVersion)
	// TODO: actually handle TF version, fetch TF binaries
	f := &TerraformTestFixture{
		BaseFixture:      NewBaseFixture(t),
		Directory:        dir,
		TerraformVersion: tfVersion,
		Options: terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: dir,
		}),
	}

	f.writeProviderTF(tfVersion, endpoint, username, password)
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
	log.Println(terraform.Destroy(f.T, f.Options))
	if os.RemoveAll(f.Directory) != nil {
		// Fatal is used here because this should never really happen,
		// and if it did, may indicate something is very wrong.
		f.T.Fatal("Unable to clean up temp directory")
	}
}

func (f *TerraformTestFixture) Init() *TerraformTestFixture {
	log.Println(terraform.Init(f.T, f.Options))
	return f
}

func (f *TerraformTestFixture) Apply() *TerraformTestFixture {
	log.Println(terraform.ApplyAndIdempotent(f.T, f.Options))
	return f
}

func (f *TerraformTestFixture) writeProviderTF(tfVersion, endpoint, username, password string) {
	f.WriteFile("provider.tf", fmt.Sprintf(`
provider "proxmox"{
  virtual_environment {
    endpoint = "%s"
    username = "%s"
    password = "%s"
    insecure = true
  }
}

terraform {
  # TODO version stuff
  #required_version = "== %s"
  required_providers {
    proxmox = {
      source = "registry.terraform.io/danitso/proxmox"
    }
  }
}

`, endpoint, username, password, tfVersion))
}
