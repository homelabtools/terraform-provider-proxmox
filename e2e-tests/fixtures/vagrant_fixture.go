package fixtures

import (
	"strings"

	"github.com/pkg/errors"
)

type VagrantTestFixture struct {
	Provider string
}

func NewVagrantTestFixture(provider string) VagrantTestFixture {
	return VagrantTestFixture{
		Provider: provider,
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
	err := runStdout("vagrant", "snapshot", "restore", name)
	return errors.WithStack(err)
}

func (f *VagrantTestFixture) SaveSnapshot(name string) error {
	err := runStdout("vagrant", "snapshot", "save", name)
	return errors.WithStack(err)
}
