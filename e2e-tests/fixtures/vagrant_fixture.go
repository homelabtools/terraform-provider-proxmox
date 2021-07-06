package fixtures

import "github.com/pkg/errors"

type VagrantTestFixture struct {
	Provider string
}

func NewVagrantTestFixture(provider string) VagrantTestFixture {
	return VagrantTestFixture{
		Provider: provider,
	}
}

func (f *VagrantTestFixture) Up() error {
	err := run("vagrant", "up", "--provider", f.Provider)
	return errors.WithStack(err)
}

func (f *VagrantTestFixture) Halt() error {
	err := run("vagrant", "halt")
	return errors.WithStack(err)
}

func (f *VagrantTestFixture) SaveSnapshot(name string) error {
	err := run("vagrant", "snapshot", "save", name)
	return errors.WithStack(err)
}

func (f *VagrantTestFixture) RestoreSnapshot(name string) error {
	err := run("vagrant", "snapshot", "restore", name)
	return errors.WithStack(err)
}
