package e2e_tests

import (
	"fmt"
	"testing"

	"github.com/danitso/terraform-provider-proxmox/e2e-tests/fixtures"
	"github.com/stretchr/testify/require"
)

var testCases = []struct {
	Name      string
	TFVersion string
}{
	{
		Name:      "TF 1.0.1",
		TFVersion: "1.0.1",
	},
}

func TestMain(t *testing.T) {
	require := require.New(t)
	// TODO: Get this from somewhere not hardcoded?
	provider := "virtualbox"
	// TODO: Better name choice
	pve := <-fixtures.NewProxmoxTestFixture(t, provider, "http://localhost:8000", "Main suite from files")

	defer pve.TearDown()

	suiteSnapshotName := "StartOfSuite"

	// If there's an existing snapshot from a previous run, reuse it. This speeds up the debugging
	// cycle, you can continually run `make debug-test` and not have to wait for the VM to be created.
	hasSnapshot, err := pve.HasSnapshot(suiteSnapshotName)
	require.NoErrorf(err, "could not determine if snapshot '%s' exists, try again and run `make clean` before running tests")
	if hasSnapshot {
		require.NoErrorf(pve.RestoreSnapshot(suiteSnapshotName), "unable to restore existing snapshot '%s'", suiteSnapshotName)
	} else {
		require.NoError(pve.SaveSnapshot(suiteSnapshotName), "unable to save snapshot at start of test suite")
	}

	for _, testCase := range testCases {
		// When the test is complete, save the current state (so that it can be inspected later) and
		// revert back to the starting state in preparation for the next test case.
		defer require.NoErrorf(pve.RestoreSnapshot(suiteSnapshotName), "unable to restore snapshot back to suite start at the end of test '%s'", testCase.Name)
		defer require.NoErrorf(pve.SaveSnapshot(fmt.Sprintf("After test case '%s'", testCase.Name)), "unable to save snapshot at end of test '%s'", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			// TODO: Take test cases from files
			tf := fixtures.NewTerraformTestFixture(t, "cases/simple", testCase.TFVersion, pve.Endpoint, "root@pam", "proxmox")
			// Save a snapshot before teardown for debugging purposes.
			require.NoErrorf(pve.SaveSnapshot(fmt.Sprintf("Before Terraform destroy for test case '%s'", testCase.Name)), "unable to save snapshot at end of test '%s'", testCase.Name)
			defer tf.TearDown()
			tf.Init().Apply()
		})
	}
}
