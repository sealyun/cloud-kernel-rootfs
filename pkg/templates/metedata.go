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
package templates

import (
	"github.com/sealyun/cloud-kernel-rootfs/pkg/logger"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/vars"
	"strings"
)

///Metadata
const metadataTemplate = `{
  "version": "{{.Version}}",
  "arch": "{{.Arch}}"
}
`

type Metadata struct {
	K8sVersion string
}

func (m *Metadata) Template() string {
	return metadataTemplate
}

func (m *Metadata) TemplateConvert() string {
	vp := vars.Platform
	arch := strings.Split(vp, "/")
	p := map[string]interface{}{
		"Version": m.K8sVersion,
		"Arch":    arch[1],
	}

	data, err := templateFromContent(m.Template(), p)
	if err != nil {
		logger.Error(err) //nolint:typecheck
		return ""
	}
	return data
}
