package e2e_tests

import (
	"testing"

	"github.com/danitso/terraform-provider-proxmox/e2e-tests/fixtures"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var endpoint = "http://localhost:8000"

type E2ESuite struct {
	suite.Suite
	ProxmoxFixture *fixtures.ProxmoxTestFixture
}

func TestE2ESuite(t *testing.T) {
	suite.Run(t, new(E2ESuite))
}

func (s *E2ESuite) SetupSuite() {
	s.ProxmoxFixture = <-fixtures.NewProxmoxTestFixture(s.T(), fixtures.ProxmoxTestFixtureOptions{})
}

func (s *E2ESuite) TearDownSuite() {
	s.ProxmoxFixture.TearDown()
}

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
	pve := <-fixtures.NewProxmoxTestFixture(t, fixtures.ProxmoxTestFixtureOptions{})
	defer pve.TearDown()

	require.NoError(pve.SaveSnapshot("StartOfSuite"), "unable to save snapshot at start of test suite")
	for _, testCase := range testCases {
		snapshotName := "TestCase-" + testCase.Name
		require.NoErrorf(pve.SaveSnapshot(snapshotName), "unable to save snapshot at start of test '%s'", testCase.Name)
		defer require.NoErrorf(pve.RestoreSnapshot(snapshotName), "unable to restore snapshot at end of test '%s'", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			tf := fixtures.NewTerraformTestFixture(t, "cases/simple", testCase.TFVersion, endpoint, "root@pam", "proxmox")
			defer tf.TearDown()
			tf.Init().Apply()
		})
	}
}
