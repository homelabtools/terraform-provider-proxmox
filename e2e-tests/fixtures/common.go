package fixtures

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

// Timestamp format that works with filenames Windows (has no colon)
const FileTimestampFormat = "Mon_Jan_02_15-04-05PM"

// run runs a command and echoes the output to stdout
func runStdout(name string, args ...string) error {
	return errors.WithStack(run(os.Stdout, name, args...))
}

// runCaptureLines runs a command and returns the output sliced by newlines
func runCaptureLines(name string, args ...string) ([]string, error) {
	buf := &bytes.Buffer{}
	err := run(buf, name, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	lines := []string{}
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}

// run runs a command and echoes the output to the specified stream
func run(out io.Writer, name string, args ...string) error {
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

// LoadExpectedResults loads expected.json from the given test directory
func LoadExpectedResults(t *testing.T, dir string) map[string]interface{} {
	expectedResults, err := os.ReadFile(filepath.Join(dir, "expected.json"))
	require.NoErrorf(t, err, "Failed to load expected test results")
	expected := map[string]interface{}{}
	require.NoError(t, json.Unmarshal(expectedResults, &expected))
	return expected
}
