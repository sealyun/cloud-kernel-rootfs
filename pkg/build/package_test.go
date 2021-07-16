package build

import (
	"testing"
)

func TestPackageDocker(t *testing.T) {
	Package("1.19.9", true)
}

func TestSaveImage(t *testing.T) {
	//Package(vars.Docker, "1.20.0", "19.03.12")
	k8s := NewDockerK8s("8.210.233.63", "1.19.9")
	k8s.SaveImages()
}
