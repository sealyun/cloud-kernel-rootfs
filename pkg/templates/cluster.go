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

///etc/Clusterfile
const clusterFileTemplate = `apiVersion: zlink.aliyun.com/v1alpha1
kind: Cluster
metadata:
  name: my-cluster
spec:
  image: kubernetes:v1.19.9
  provider: ALI_CLOUD
  ssh:
    passwd: Seadent123
    pk: xxx
    pkPasswd: xxx
    user: root
  network:
    interface: eth0
    cniName: calico
    podCIDR: 100.64.0.0/10
    svcCIDR: 10.96.0.0/22
    withoutCNI: false
  certSANS:
    - aliyun-inc.com
    - 10.0.0.2
  masters:
    cpu: 4
    memory: 4
    count: 1
    systemDisk: 100
    dataDisks:
    - 100
  nodes:
    cpu: 4
    memory: 4
    count: 0
    systemDisk: 100
    dataDisks:
    - 100
`

type Cluster struct {
}

func (c *Cluster) Template() string {
	return clusterFileTemplate
}

func (c *Cluster) TemplateConvert() string {
	panic("implement me")
}
