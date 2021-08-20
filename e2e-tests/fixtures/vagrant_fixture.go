package fixtures

import (
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
)

const DisableSnapshotsEnvVar = "DISABLE_SNAPSHOTS"

type VagrantTestFixture struct {
	Provider string
	// Optionally disable snapshots. This logic is handled here so that
	// callers can keep test code concise and do not need to wrap every
	// snapshot operation in an IF statement.
	disableSnapshots bool
}

func NewVagrantTestFixture(provider string) VagrantTestFixture {
	return VagrantTestFixture{
		Provider:         provider,
		disableSnapshots: os.Getenv(DisableSnapshotsEnvVar) != "",
	}
}

func (f *VagrantTestFixture) Up() error {
	err := runStdout("vagrant", "up", "--provider", f.Provider)
	return errors.WithStack(err)
}

func (f *VagrantTestFixture) Halt() error {
	err := runStdout("vagrant", "halt")
	return errors.WithStack(err)
}

func (f *VagrantTestFixture) HasSnapshot(name string) (bool, error) {
	if f.disableSnapshots {
		return false, nil
	}
	snapshots, err := f.ListSnapshots()
	if err != nil {
		return false, errors.WithStack(err)
	}
	for _, snapshot := range snapshots {
		if snapshot == name {
			return true, nil
		}
	}
	return false, nil
}

func (f *VagrantTestFixture) ListSnapshots() ([]string, error) {
	if f.disableSnapshots {
		return []string{}, nil
	}
	lines, err := runCaptureLines("vagrant", "--machine-readable", "snapshot", "list")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// Parse the output. Snapshot list output looks like this:
	//   1626843761,default,metadata,provider,virtualbox
	//   1626843761,default,ui,output,==> default:
	//   1626843761,default,ui,detail,snapshotName1
	//   1626843761,default,ui,detail,snapshotName2

	// Output is expected to be at least 3 lines, two for the provider and VM info, one or more for the snapshots.
	if len(lines) < 3 || strings.Contains(lines[1], "No snapshots have been taken yet") {
		return []string{}, nil
	}

	// Skip the first two lines which show provider name and VM name.
	lines = lines[2:]
	result := make([]string, len(lines))
	for _, line := range lines {
		columns := strings.Split(line, ",")
		if len(columns) < 5 {
			continue
		}
		result = append(result, columns[4])
	}
	return result, nil
}

func (f *VagrantTestFixture) RestoreSnapshot(name string) error {
	if f.disableSnapshots {
		return nil
	}
	log.Printf("Attempting to restore snapshot '%s'...", name)
	err := runStdout("vagrant", "snapshot", "restore", name)
	if err != nil {
		log.Printf("Failed restoring snapshot '%s': %s\n", name, err.Error())
	}
	return errors.WithStack(err)
}

func (f *VagrantTestFixture) SaveSnapshot(name string) error {
	if f.disableSnapshots {
		return nil
	}
	log.Printf("Attempting to save snapshot '%s'...", name)
	err := runStdout("vagrant", "snapshot", "save", name)
	if err != nil {
		log.Printf("Failed saving snapshot '%s': %s\n", name, err.Error())
	}
	return errors.WithStack(err)
}
