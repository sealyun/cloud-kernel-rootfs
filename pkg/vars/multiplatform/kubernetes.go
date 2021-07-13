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
package multiplatform

import (
	"fmt"
)

func NewKubernetes(version, rootfs string, platform Platform) Vars {
	v := &kubernetes{
		info: info{
			Version: version,
			Rootfs:  rootfs,
		},
	}
	v.setWgetURL(platform)
	return v
}

type kubernetes struct {
	info    //nolint:Recheck
	wgetURL string
}

func (c *kubernetes) setWgetURL(platform Platform) { //nolint:typecheck
	wurl := "https://dl.k8s.io/v%s/kubernetes-server-%s.tar.gz"
	splatform := ""
	switch platform {
	case LinuxAmd64: //nolint:typecheck
		splatform = "linux-amd64" //nolint:typecheck
	case LinuxArm64: //nolint:typecheck
		splatform = "linux-arm64" //nolint:typecheck
	default:
		return
	}
	c.wgetURL = fmt.Sprintf(wurl, c.Version, splatform)
}

func (c *kubernetes) FetchWgetURL() string {
	return c.wgetURL
}

func (c *kubernetes) FinalShell() string { //nolint:typecheck
	shell := "wget %s -O  kubernetes-server.tar.gz &&  tar xf kubernetes-server.tar.gz && " +
		"cp kubernetes/server/bin/kubectl %s/bin/ && " +
		"cp kubernetes/server/bin/kubelet %s/bin/ && " +
		"cp kubernetes/server/bin/kubeadm %s/bin/ "
	return fmt.Sprintf(shell, c.wgetURL, c.Rootfs, c.Rootfs, c.Rootfs)
}
