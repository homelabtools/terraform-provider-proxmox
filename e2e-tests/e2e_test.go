package e2e_tests

import (
	"testing"

	"github.com/danitso/terraform-provider-proxmox/e2e-tests/fixtures"
)

func TestMain(t *testing.T) {
	tf := fixtures.NewTerraformTestFixture(t, "test", "1.0.1")
	defer tf.TearDown()

	var pve *fixtures.ProxmoxTestFixture
	defer func() { pve.TearDown() }()
	// endpoint := "https://localhost:8006"
	pve = <-fixtures.NewProxmoxTestFixture(t, fixtures.ProxmoxTestFixtureOptions{})
}
