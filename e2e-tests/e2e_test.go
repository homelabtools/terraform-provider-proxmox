package e2e_tests

import (
	"fmt"
	"testing"

	"github.com/danitso/terraform-provider-proxmox/e2e-tests/fixtures"
	"github.com/stretchr/testify/require"
)

const TestUserName = "root@pam"
const TestPassword = "proxmox"

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
	pve := <-fixtures.NewProxmoxTestFixture(t, provider, "http://localhost:8000", "Main suite from files", TestUserName, TestPassword)

	defer pve.TearDown()

	startSnapshotName := t.Name() + " Start"

	// If there's an existing snapshot from a previous run, reuse it. This speeds up the debugging
	// cycle, you can continually run `make debug-test` and not have to wait for the VM to be created.
	hasSnapshot, err := pve.HasSnapshot(startSnapshotName)
	require.NoErrorf(err, "could not determine if snapshot '%s' exists, try again and run `make clean` before running tests")
	if hasSnapshot {
		require.NoErrorf(pve.RestoreSnapshot(startSnapshotName), "unable to restore existing snapshot '%s'", startSnapshotName)
	} else {
		require.NoError(pve.SaveSnapshot(startSnapshotName), "unable to save snapshot at start of test suite")
	}

	for _, testCase := range testCases {
		// When the test is complete, save the current state (so that it can be inspected later) and
		// revert back to the starting state in preparation for the next test case.
		// --- DO NOT change the order of these defer statements.
		defer require.NoErrorf(pve.RestoreSnapshot(startSnapshotName), "unable to restore snapshot back to suite start at the end of test '%s'", testCase.Name)
		defer require.NoErrorf(pve.SaveSnapshot(fmt.Sprintf("After test case '%s'", testCase.Name)), "unable to save snapshot at end of test '%s'", testCase.Name)
		// ---

		t.Run(testCase.Name, func(t *testing.T) {
			// TODO: Take test cases from files
			tf := fixtures.NewTerraformTestFixture(t, "cases/simple", testCase.TFVersion, pve.Endpoint, TestUserName, TestPassword)
			//expected := fixtures.LoadExpectedResults(t, tf.Directory)

			// --- DO NOT change order of these defer statements
			defer tf.TearDown()
			// Save a snapshot before teardown for debugging purposes.
			defer require.NoErrorf(pve.SaveSnapshot(fmt.Sprintf("Before Terraform destroy for test case '%s'", testCase.Name)), "unable to save snapshot at end of test '%s'", testCase.Name)
			// ---

			tf.Init().Apply()

			// TODO: Evaluate results
		})
	}
}

func evaluateResults(t *testing.T, expected map[string]interface{}) {
	for apiKey, apiVal := range expected {
		for _, item := range apiVal.(map[string]interface{}) {

			//for attrName, attrValue := range item.(map[string]interface{}) {
			//}
		}
	}
}
