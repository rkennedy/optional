//go:build ignore

// Usage:
//   go run mage.go

package main

import (
	"os"

	"github.com/magefile/mage/mage"
)

func main() {
	os.Exit(mage.Main())
}
