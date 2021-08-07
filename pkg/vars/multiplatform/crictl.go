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

func NewCriCTL(version, rootfs string, platform Platform) Vars {
	v := &crictl{
		info: info{
			Version: version,
			Rootfs:  rootfs,
		},
	}
	v.setWgetURL(platform)
	return v
}

type crictl struct { //nolint:typecheck
	info    //nolint:Recheck
	wgetURL string
}

func (c *crictl) setWgetURL(platform Platform) {
	wurl := "https://github.com/kubernetes-sigs/cri-tools/releases/download/v%s/crictl-v%s-%s.tar.gz"
	splatform := ""
	switch platform {
	case LinuxAmd64: //nolint:typecheck
		splatform = "linux-amd64" //nolint:typecheck
	case LinuxArm64: //nolint:typecheck
		splatform = "linux-arm64" //nolint:typecheck
	default:
		return
	}
	c.wgetURL = fmt.Sprintf(wurl, c.Version, c.Version, splatform)
}

func (c *crictl) FetchWgetURL() string {
	return c.wgetURL
}

func (c *crictl) FinalShell() string { //nolint:typecheck
	shell := "wget %s -O  crictl.tar.gz &&  tar xf crictl.tar.gz && cp crictl %s/bin/"
	return fmt.Sprintf(shell, c.wgetURL, c.Rootfs)
}
