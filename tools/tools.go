// +build tools

// The purpose of this file is to allow tools meant for only building and testing,
// such as ginkgo or mocking tools, to be brought in and managed as Go modules.
// This makes builds more reproducible and easier to run.
// Tools can be run with `go run <github.com/name/of/package>`

package tools

import (
	_ "github.com/onsi/ginkgo/ginkgo"
)
