/*
Copyright 2021 cuisongliu@qq.com.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package multiplatform //nolint:gofmt

import (
	"fmt"
)

func NewRootfs(version string, platform Platform, release bool) Vars {
	v := &rootfs{
		info: info{
			Version: version,
		},
		release: release,
	}
	v.setPlatform(platform)
	return v
}

type rootfs struct {
	info
	imageName string
	tarName   string
	release   bool
}

func (c *rootfs) setPlatform(platform Platform) {
	splatform := ""
	switch platform {
	case LinuxAmd64:
		splatform = "amd64"
	case LinuxArm64:
		splatform = "arm64"
	default:
		return
	}
	Version := c.Version
	if !c.release {
		Version = c.Version + "_develop"
	}
	c.imageName = fmt.Sprintf("kubernetes-%s:v%s", splatform, Version)
	c.tarName = fmt.Sprintf("kube-%s-v%s.tar", splatform, Version)
}

func (c *rootfs) FetchWgetURL() string {
	return c.imageName
}

func (c *rootfs) FinalShell() string {
	shell := `sealer build -f Kubefile -b "local" -t %s .`
	return fmt.Sprintf(shell, c.imageName)
}
