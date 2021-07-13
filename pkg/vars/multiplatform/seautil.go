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

func NewSeautil(version, rootfs string, platform Platform) Vars {
	v := &seautil{
		info: info{
			Version: version,
			Rootfs:  rootfs,
		},
	}
	v.setWgetURL(platform)
	return v
}

type seautil struct {
	info
	wgetURL string
}

func (c *seautil) setWgetURL(platform Platform) {
	wurl := "https://github.com/alibaba/sealer/releases/download/v%s/seautil-v%s-%s.tar.gz"
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
func (c *seautil) FetchWgetURL() string {
	return c.wgetURL
}
func (c *seautil) FinalShell() string {
	shell := "wget %s -O  seautil.tar.gz &&  tar xf seautil.tar.gz && cp seautil %s/bin/"
	return fmt.Sprintf(shell, c.wgetURL, c.Rootfs)
}
