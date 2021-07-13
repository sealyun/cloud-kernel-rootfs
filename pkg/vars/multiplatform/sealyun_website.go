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

func NewSealyunWebsite(platform Platform) Vars {
	v := &sealyunWebsite{
		info: info{},
	}
	v.setWgetURL(platform)
	return v
}

type sealyunWebsite struct {
	info
	wgetURL string
}

func (c *sealyunWebsite) setWgetURL(platform Platform) { //nolint:typecheck
	splatform := ""
	switch platform {
	case LinuxAmd64: //nolint:typecheck
		splatform = "kubernetes-image-amd64" //nolint:typecheck
	case LinuxArm64: //nolint:typecheck
		splatform = "kubernetes-image-arm64" //nolint:typecheck
	default:
		return
	}
	c.wgetURL = splatform
}

func (c *sealyunWebsite) FetchWgetURL() string {
	return c.wgetURL
}

func (c *sealyunWebsite) FinalShell() string { //nolint:typecheck
	return ""
}
