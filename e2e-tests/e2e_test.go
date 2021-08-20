package e2e_tests

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/danitso/terraform-provider-proxmox/e2e-tests/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TestUserName = "root@pam"
const TestPassword = "proxmox"
const ProviderEnvVar = "PROVIDER"
const ProxmoxEndpointEnvVar = "PROXMOX_ENDPOINT"
const HTTPProxyEnvVar = "HTTP_PROXY"
const TestCasesPath = "cases"

var testScenarios = []struct {
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

	provider := os.Getenv(ProviderEnvVar)
	require.NotEmptyf(provider, "Must define env var '%s' that defines the Vagrant provider", ProviderEnvVar)

	endpoint := os.Getenv(ProxmoxEndpointEnvVar)
	require.NotEmptyf(endpoint, "Must define env var '%s' that defines the Proxmox endpoint", ProxmoxEndpointEnvVar)

	proxy := os.Getenv(HTTPProxyEnvVar)

	pve := <-fixtures.NewProxmoxTestFixture(
		t,
		provider,
		endpoint,
		proxy,
		"Main suite from files",
		TestUserName,
		TestPassword)

	defer pve.TearDown()

	baseStateSnapshotName := t.Name() + " Base"

	// If there's an existing snapshot from a previous run, reuse it. This speeds up the debugging
	// cycle, you can continually run `make debug-test` and not have to wait for the VM to be created.
	hasSnapshot, err := pve.HasSnapshot(baseStateSnapshotName)
	require.NoErrorf(err, "could not determine if snapshot '%s' exists, try again and run `make clean` before running tests")
	if hasSnapshot {
		require.NoErrorf(pve.RestoreSnapshot(baseStateSnapshotName), "unable to restore existing snapshot '%s'", baseStateSnapshotName)
	} else {
		require.NoError(pve.SaveSnapshot(baseStateSnapshotName), "unable to save snapshot at start of test suite")
	}

	// Give a unique name to each test run so that future runs do not collide with it
	beforeScenarioSnapshotName := fmt.Sprintf("Start of %s at %s", t.Name(), timestamp())
	require.NoError(pve.SaveSnapshot(beforeScenarioSnapshotName), "unable to save snapshot at start of test suite")

	for _, scenario := range testScenarios {
		// Wrap with a function so that defer statements below are scoped to each iteration of the loop
		func() {
			log.Printf("Now testing using Terraform version '%s'\n", scenario.TFVersion)

			startOfScenarioSnapshotName := fmt.Sprintf("Start of %s %s %s", t.Name(), scenario.Name, timestamp())
			require.NoError(pve.SaveSnapshot(startOfScenarioSnapshotName), "unable to save snapshot at start of scenario '%s'", startOfScenarioSnapshotName)

			// When the test is complete, save the current state (so that it can be inspected later) and
			// revert back to the starting state in preparation for the next test case.
			// --- DO NOT change order of these defer statements ---
			defer require.NoErrorf(pve.RestoreSnapshot(beforeScenarioSnapshotName), "unable to restore snapshot back to suite start at the end of scenario '%s'", scenario.Name)
			// TODO: This is more useful when snapshotting is disabled on individual test cases, which provides some more
			//       potential test coverage by way of not starting from a fresh state.
			defer require.NoErrorf(pve.SaveSnapshot(fmt.Sprintf("After test scenario '%s' %s", scenario.Name, timestamp())), "unable to save snapshot at end of scenario '%s'", scenario.Name)
			// -----------------------------------------------------

			dirs, err := os.ReadDir(TestCasesPath)
			require.NoErrorf(err, "could not open test cases from directory '%s'", TestCasesPath)
			for _, dir := range dirs {
				// TODO: Use an abstraction like AeroFS
				testCasePath := filepath.Join(TestCasesPath, dir.Name())
				fullTestName := fmt.Sprintf("%s_%s", scenario.Name, testCasePath)

				// Individual test cases start being run here
				t.Run(fullTestName, func(t *testing.T) {
					tf := fixtures.NewTerraformTestFixture(t, testCasePath, scenario.TFVersion, pve.Endpoint, TestUserName, TestPassword)
					expected := fixtures.LoadExpectedResults(t, tf.Directory)
					t.Log(expected)

					// --- DO NOT change order of these defer statements ---
					// TODO: move to end of range dirs, after t.Run() func? May be OK if testing recovers
					//       from any panic within this function (it should, but verify)
					defer require.NoErrorf(pve.RestoreSnapshot(startOfScenarioSnapshotName), "unable to restore snapshot back to start of scenario '%s'", scenario.Name)
					defer tf.TearDown()
					// Save a snapshot before teardown, for debugging purposes.
					defer require.NoErrorf(
						pve.SaveSnapshot(
							fmt.Sprintf("Before Terraform destroy for test case '%s' %s", fullTestName, timestamp())),
						"unable to save snapshot at end of test '%s'", scenario.Name)
					// -----------------------------------------------------

					tf.Init().Apply()

					for apiKey, apiVal := range expected {
						var resp map[string]interface{}
						err := pve.PVEClient.DoRequest(http.MethodGet, apiKey, nil, &resp)
						assert.NoErrorf(t, err, "Unexpected error when calling API '%s'", apiKey)
						data := resp["data"].([]interface{})
						assert.Subsetf(t, data, apiVal, "API result does not match expected test data for items under '%s'")
					}
				})
			}
		}()
	}
}

func timestamp() string {
	return time.Now().Format(time.Stamp)
}
