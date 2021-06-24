package e2e_tests

import (
	"testing"

	"github.com/danitso/terraform-provider-proxmox/e2e-tests/fixtures"
)

func TestMain(t *testing.T) {
	tf := fixtures.NewTerraformTestFixture(t, "test", "1.0.0")
	defer tf.TearDown()
	pve := fixtures.NewProxmoxTestFixture(t, "https://localhost:8006")
	pve.Start()
	defer pve.TearDown()
}
