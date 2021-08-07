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

func NewDocker(version, rootfs string, platform Platform) Vars {
	v := &docker{
		info: info{
			Version: version,
			Rootfs:  rootfs,
		},
	}
	v.setWgetURL(platform)
	return v
}

type docker struct {
	info
	wgetURL string
}

func (c *docker) setWgetURL(platform Platform) {
	//https://github.com/osemp/moby/releases/download/v19.03.14/docker-amd64.tar.gz
	//https://github.com/osemp/moby/releases/download/v19.03.14/docker-arm64.tar.gz
	wurl := "https://github.com/osemp/moby/releases/download/v%s/docker-%s.tar.gz"
	splatform := ""
	switch platform {
	case LinuxAmd64: //nolint:typecheck
		splatform = "amd64" //nolint:typecheck
	case LinuxArm64: //nolint:typecheck
		splatform = "arm64" //nolint:typecheck
	default:
		return
	}
	c.wgetURL = fmt.Sprintf(wurl, c.Version, splatform)
}

func (c *docker) FetchWgetURL() string {
	return c.wgetURL
}

func (c *docker) FinalShell() string { //nolint:typecheck
	shell := "wget %s -O  docker.tar.gz"
	return fmt.Sprintf(shell, c.wgetURL)
}
