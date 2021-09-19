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

func NewNerdctl(version, rootfs string, platform Platform) Vars {
	v := &nerdctl{
		info: info{
			Version: version,
			Rootfs:  rootfs,
		},
	}
	v.setWgetURL(platform)
	return v
}

type nerdctl struct {
	info
	wgetURL string
}

func (c *nerdctl) setWgetURL(platform Platform) { //nolint:typecheck
	wurl := "https://github.com/containerd/nerdctl/releases/download/v%s/nerdctl-%s-%s.tar.gz"
	splatform := ""
	switch platform {
	case LinuxAmd64: //nolint:typecheck
		splatform = "linux-amd64"
	case LinuxArm64: //nolint:typecheck
		splatform = "linux-arm64"
	default:
		return
	}
	c.wgetURL = fmt.Sprintf(wurl, c.Version, c.Version, splatform)
}
func (c *nerdctl) FetchWgetURL() string {
	return c.wgetURL
}
func (c *nerdctl) FinalShell() string {
	shell := "wget %s -O  nerdctl.tar.gz &&  tar xf nerdctl.tar.gz && cp nerdctl %s/bin/ && cp containerd-rootless* %s/bin/"
	return fmt.Sprintf(shell, c.wgetURL, c.Rootfs, c.Rootfs)
}
