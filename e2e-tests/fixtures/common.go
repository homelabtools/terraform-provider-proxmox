package fixtures

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/pkg/errors"
)

// Run runs a command and echoes the output to stdout
func run(name string, args ...string) error {
	log.Printf("Running command `%s` with args %#q", name, args)
	c := exec.Command(name, args...)
	c.Stderr = c.Stdout
	stdout, err := c.StdoutPipe()
	if err != nil {
		return errors.WithStack(err)
	}
	err = c.Start()
	if err != nil {
		return errors.WithMessagef(err, "failed trying to run '%s' with args '%#q'", name, args)
	}
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		fmt.Print(string(tmp))
		if err != nil {
			break
		}
	}
	err = c.Wait()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
