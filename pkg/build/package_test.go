package build

import (
	"testing"
)

func TestPackageDocker(t *testing.T) {
	Package("1.19.9", true)
}
