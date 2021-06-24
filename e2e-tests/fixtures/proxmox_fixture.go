package fixtures

import (
	"testing"
)

type ProxmoxTestFixture struct {
	BaseFixture
	Endpoint string
}

// NewProxmoxTestFixture creates a new test fixture for working with Proxmox.
// endpoint - URL of Proxmox API
func NewProxmoxTestFixture(t *testing.T, endpoint string) chan *ProxmoxTestFixture {
	c := make(chan *ProxmoxTestFixture, 1)
	func() {
		f := &ProxmoxTestFixture{
			BaseFixture: NewBaseFixture(t),
			Endpoint:    endpoint,
		}
		err := f.start()
		f.Require.NoError(err, "failed starting Proxmox VM with Vagrant")
		c <- f
	}()
	return c
}

// start brings up the Proxmox VM
func (f *ProxmoxTestFixture) start() error {
	_, err := run("make", "up")
	f.Require.NoError(err, "expected `make up` to work")
	return nil
}

// TearDown removes every trace the test fixture.
// It should be called with defer right after creating the fixture.
func (f *ProxmoxTestFixture) TearDown() {
	if !f.ShouldClean(f) {
		return
	}
	_, err := run("make", "destroy")
	f.Assert.NoError(err, "expected `make destroy` to work")
}
