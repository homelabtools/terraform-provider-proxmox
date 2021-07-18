package e2e_tests

import (
	"testing"

	"github.com/danitso/terraform-provider-proxmox/e2e-tests/fixtures"
)

var endpoint = "http://localhost:8000"

func TestMain(t *testing.T) {
	tf := fixtures.NewTerraformTestFixture(t, "cases/simple", "1.0.1", endpoint, "root@pam", "proxmox")
	defer tf.TearDown()

	var pve *fixtures.ProxmoxTestFixture
	defer func() { pve.TearDown() }()
	pve = <-fixtures.NewProxmoxTestFixture(t, fixtures.ProxmoxTestFixtureOptions{})

	tf.Init().Apply()
}
