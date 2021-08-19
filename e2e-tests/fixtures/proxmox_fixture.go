package fixtures

import (
	"math/rand"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// ProxmoxTestFixture is a test helper for bringing up Vagrant VMs that run Proxmox.
type ProxmoxTestFixture struct {
	BaseFixture
	VagrantTestFixture
	// The Vagrant provider to use, defaults to virtualbox
	VagrantProvider string
	// Name is a descriptive name for this test fixture.
	Name string
	// URL of Proxmox instance
	Endpoint  string
	PVEClient *proxmox.VirtualEnvironmentClient
}

// NewProxmoxTestFixture creates a new Vagrant-based test fixture for working with Proxmox.
// Calling this function will asynchronously bring up a VM for running Proxmox.
func NewProxmoxTestFixture(t *testing.T, vagrantProvider, endpoint, httpProxy, name, username, password string) chan *ProxmoxTestFixture {
	base := NewBaseFixture(t)
	c := make(chan *ProxmoxTestFixture, 1)
	pveClient, err := proxmox.NewVirtualEnvironmentClient(endpoint,
		username,
		password,
		"",
		true,
		func(r *http.Request) (*url.URL, error) {
			// A custom proxy function is used here, even though the proxy is usually coming from
			// an env var, because the regular http.ProxyFromEnvironment will return nil for no proxy
			// if the HTTP request address is localhost or 127.0.0.1. For testing the Proxmox endpoint
			// is usually on localhost, so http.ProxyFromEnvironment will not work.
			if httpProxy == "" {
				return nil, nil
			}
			return url.Parse(httpProxy)
		},
	)
	require.NoError(t, err)
	func() {
		f := &ProxmoxTestFixture{
			BaseFixture:        base,
			VagrantTestFixture: NewVagrantTestFixture(vagrantProvider),
			VagrantProvider:    vagrantProvider,
			Name:               name,
			Endpoint:           endpoint,
			PVEClient:          pveClient,
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
	f.Require.NoErrorf(err, "failed to bring up VM for fixture '%s'", f.Name)
}

// TearDown removes every trace the test fixture.
// It should be called with defer right after creating the fixture.
func (f *ProxmoxTestFixture) TearDown() {
	if !f.ShouldClean(f) {
		return
	}
	// Turn off the VM.
	err := f.Halt()
	f.Assert.NoErrorf(err, "failed shutting down VM for fixture '%s'", f.Name)
}
