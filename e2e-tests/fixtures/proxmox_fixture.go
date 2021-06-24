package fixtures

import (
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/pkg/errors"
)

type ProxmoxTestFixture struct {
	BaseFixture
	Endpoint string
}

// NewProxmoxTestFixture creates a new test fixture for working with Proxmox.
// endpoint - URL of Proxmox API
func NewProxmoxTestFixture(t *testing.T, endpoint string) *ProxmoxTestFixture {
	return &ProxmoxTestFixture{
		BaseFixture: NewBaseFixture(t),
		Endpoint:    endpoint,
	}
}

// Start brings up the Proxmox VM
func (f *ProxmoxTestFixture) Start() error {
	_, err := run("make", "up")
	f.Require.NoError(err, "expected `make up` to work")
	return nil
}

// TearDown removes every trace the test fixture.
// It should be called with defer right after creating the fixture.
func (f *ProxmoxTestFixture) TearDown() {
	if os.Getenv("SKIP_CLEANUP") != "" {
		f.T.Logf("SKIP_CLEANUP env var found, skipping cleanup of '%s'", f.Endpoint)
		return
	}
	_, err := run("make", "destroy")
	f.Assert.NoError(err, "expected `make destroy` to work")
}

// Run runs a command and returns its output (stdout and stderr combined)
func run(name string, args ...string) (string, error) {
	log.Printf("Running command `%s` with args %#q", name, args)
	c := exec.Command(name, args...)
	out, err := c.CombinedOutput()
	if err != nil {
		return "", errors.WithMessagef(err, "failed trying to run '%s' with args '%#q'", name, args)
	}
	result := string(out)
	// TODO: stream output to stdout as well, look into v10x emulator library
	log.Println(result)
	return result, nil
}
