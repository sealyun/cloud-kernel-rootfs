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

func NewRegistry(version string, platform Platform) Vars {
	v := &registry{
		info: info{
			Version: version,
		},
	}
	v.setWgetURL(platform)
	return v
}

type registry struct {
	info
	wgetURL string
}

func (c *registry) setWgetURL(platform Platform) {
	//pull ghcr.io/osemp/distribution-arm64/distribution:2.7.1
	//pull ghcr.io/osemp/distribution-amd64/distribution:2.7.1
	wurl := "ghcr.io/osemp/distribution-%s/distribution:%s"
	splatform := ""
	switch platform {
	case LinuxAmd64: //nolint:typecheck
		splatform = "amd64" //nolint:typecheck
	case LinuxArm64: //nolint:typecheck
		splatform = "arm64" //nolint:typecheck
	default:
		return
	}
	c.wgetURL = fmt.Sprintf(wurl, splatform, c.Version)
}

func (c *registry) FetchWgetURL() string {
	return c.wgetURL
}

func (c *registry) FinalShell() string { //nolint:typecheck
	return ""
}
