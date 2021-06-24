package fixtures

import (
	"log"
	"os/exec"

	"github.com/pkg/errors"
)

// Run runs a command and returns its output (stdout and stderr combined)
func run(name string, args ...string) (string, error) {
	log.Printf("Running command `%s` with args %#q", name, args)
	c := exec.Command(name, args...)
	out, err := c.CombinedOutput()
	if err != nil {
		return "", errors.WithMessagef(err, "failed trying to run '%s' with args '%#q'", name, args)
	}
	result := string(out)
	// TODO: stream output to stdout as well, look into v10x emulator library
	log.Println(result)
	return result, nil
}
