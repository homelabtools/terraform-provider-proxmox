package fixtures

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/imdario/mergo"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// ProxmoxTestFixture is a test helper for bringing up Vagrant VMs that run Proxmox.
type ProxmoxTestFixture struct {
	BaseFixture
	VagrantTestFixture
	UseSnapshots      bool
	FixtureName       string
	vagrantProvider   string
	snapshotStartName string
	snapshotEndName   string
}

// ProxmoxTestFixtureOptions is the options struct for NewProxmoxTestFixture.
type ProxmoxTestFixtureOptions struct {
	// The Vagrant provider to use, defaults to virtualbox
	VagrantProvider string
	// FixtureName is a descriptive name for this test fixture.
	FixtureName  string
	UseSnapshots bool
}

var defaultOptions = ProxmoxTestFixtureOptions{
	VagrantProvider: "virtualbox",
	FixtureName:     fmt.Sprintf("fixture-%d", rand.Intn(1000)),
	UseSnapshots:    false,
}

// NewProxmoxTestFixture creates a new Vagrant-based test fixture for working with Proxmox.
// Calling this function will asynchronously bring up a VM for running Proxmox.
func NewProxmoxTestFixture(t *testing.T, opts ProxmoxTestFixtureOptions) chan *ProxmoxTestFixture {
	f := NewBaseFixture(t)
	f.Require.NoError(mergo.Merge(&opts, defaultOptions), "failed merging default options")
	c := make(chan *ProxmoxTestFixture, 1)
	func() {
		now := time.Now().Format(time.RFC822)
		f := &ProxmoxTestFixture{
			BaseFixture:        f,
			VagrantTestFixture: NewVagrantTestFixture(opts.VagrantProvider),
			snapshotStartName:  opts.FixtureName + " start " + now,
			snapshotEndName:    opts.FixtureName + " end " + now,
			vagrantProvider:    opts.VagrantProvider,
			FixtureName:        opts.FixtureName,
			UseSnapshots:       opts.UseSnapshots,
		}
		f.start()
		c <- f
	}()
	return c
}

// start brings up the Proxmox VM
func (f *ProxmoxTestFixture) start() {
	// Bring up the VM
	err := f.Up()
	f.Require.NoErrorf(err, "failed to bring up VM for fixture '%s'", f.FixtureName)
}

// TearDown removes every trace the test fixture.
// It should be called with defer right after creating the fixture.
func (f *ProxmoxTestFixture) TearDown() {
	if !f.ShouldClean(f) {
		return
	}
	// Turn off the VM.
	err := f.Halt()
	f.Assert.NoErrorf(err, "failed shutting down VM for fixture '%s'", f.FixtureName)
}
