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

import "github.com/sealyun/cloud-kernel-rootfs/pkg/logger"

///etc/kubeadm.yaml
const kubeadmYamlTemplate = `apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration
networking:
  podSubnet: 100.64.0.0/10
kubernetesVersion: {{.Version}}
---
apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
mode: "ipvs"
`

type Kubeadm struct {
	K8sVersion string
}

func (k *Kubeadm) Template() string {
	return kubeadmYamlTemplate
}

func (k *Kubeadm) TemplateConvert() string {
	p := map[string]interface{}{
		"Version": k.K8sVersion,
	}

	data, err := templateFromContent(k.Template(), p)
	if err != nil {
		logger.Error(err) //nolint:typecheck
		return ""
	}
	return data
}
