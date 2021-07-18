package e2e_tests

import (
	"testing"

	"github.com/danitso/terraform-provider-proxmox/e2e-tests/fixtures"
)

var endpoint = "http://localhost:8000"

// TODO: Have Proxmox fixture optionally created during suite setup in addition to test setup.

func TestMain(t *testing.T) {
	// Proxmox fixture should come first. Specifically, Proxmox TearDown should happen *after* TF TearDown.
	// Otherwise TF TearDown fails because the VM no longer exists.
	var pve *fixtures.ProxmoxTestFixture
	defer func() { pve.TearDown() }()
	pve = <-fixtures.NewProxmoxTestFixture(t, fixtures.ProxmoxTestFixtureOptions{})

	tf := fixtures.NewTerraformTestFixture(t, "cases/simple", "1.0.1", endpoint, "root@pam", "proxmox")
	defer tf.TearDown()

	tf.Init().Apply()
}
