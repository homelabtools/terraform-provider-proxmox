package fixtures

import (
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// Timestamp format that works with filenames Windows (has no colon)
const FileTimestampFormat = "Mon_Jan_02_15-04-05PM"

// run runs a command and echoes the output to stdout
func run(name string, args ...string) error {
	return errors.WithStack(runf(os.Stdout, name, args...))
}

// runf runs a command and echoes the output to the specified stream
func runf(out io.Writer, name string, args ...string) error {
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
	// Stream output asynchronously
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		out.Write(tmp)
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
