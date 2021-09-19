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
	"strings"
)

const (
	InitTemplateText = string(InitConfigurationDefault +
		ClusterConfigurationDefault +
		kubeproxyConfigDefault +
		kubeletConfigDefault)

	InitConfigurationDefault = `apiVersion: {{.KubeadmAPI}}
kind: InitConfiguration
localAPIEndpoint:
  advertiseAddress: {{.Master0}}
  bindPort: 6443
nodeRegistration:
  criSocket: {{.CriSocket}}
`

	ClusterConfigurationDefault = `---
apiVersion: {{.KubeadmAPI}}
kind: ClusterConfiguration
kubernetesVersion: {{.Version}}
networking:
  podSubnet: 100.64.0.0/10
`
	kubeproxyConfigDefault = `
---
apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
mode: "ipvs"
ipvs:
  excludeCIDRs:
  - "{{.VIP}}/32"
`

	kubeletConfigDefault = `
---
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
authentication:
  anonymous:
    enabled: false
  webhook:
    cacheTTL: 2m0s
    enabled: true
  x509:
    clientCAFile: /etc/kubernetes/pki/ca.crt
authorization:
  mode: Webhook
  webhook:
    cacheAuthorizedTTL: 5m0s
    cacheUnauthorizedTTL: 30s
cgroupDriver: {{ .CgroupDriver}}
cgroupsPerQOS: true
clusterDomain: cluster.local
configMapAndSecretChangeDetectionStrategy: Watch
containerLogMaxFiles: 5
containerLogMaxSize: 10Mi
contentType: application/vnd.kubernetes.protobuf
cpuCFSQuota: true
cpuCFSQuotaPeriod: 100ms
cpuManagerPolicy: none
cpuManagerReconcilePeriod: 10s
enableControllerAttachDetach: true
enableDebuggingHandlers: true
enforceNodeAllocatable:
- pods
eventBurst: 10
eventRecordQPS: 5
evictionHard:
  imagefs.available: 15%
  memory.available: 100Mi
  nodefs.available: 10%
  nodefs.inodesFree: 5%
evictionPressureTransitionPeriod: 5m0s
failSwapOn: true
fileCheckFrequency: 20s
hairpinMode: promiscuous-bridge
healthzBindAddress: 127.0.0.1
healthzPort: 10248
httpCheckFrequency: 20s
imageGCHighThresholdPercent: 85
imageGCLowThresholdPercent: 80
imageMinimumGCAge: 2m0s
iptablesDropBit: 15
iptablesMasqueradeBit: 14
kubeAPIBurst: 10
kubeAPIQPS: 5
makeIPTablesUtilChains: true
maxOpenFiles: 1000000
maxPods: 110
nodeLeaseDurationSeconds: 40
nodeStatusReportFrequency: 10s
nodeStatusUpdateFrequency: 10s
oomScoreAdj: -999
podPidsLimit: -1
port: 10250
registryBurst: 10
registryPullQPS: 5
rotateCertificates: true
runtimeRequestTimeout: 2m0s
serializeImagePulls: true
staticPodPath: /etc/kubernetes/manifests
streamingConnectionIdleTimeout: 4h0m0s
syncFrequency: 1m0s
volumeStatsAggPeriod: 1m0s`
	ContainerdShell = `if grep "SystemdCgroup = true"  /etc/containerd/config.toml &> /dev/null; then  
driver=systemd
else
driver=cgroupfs
fi
echo ${driver}`
	DockerShell = `driver=$(docker info -f "{{.CgroupDriver}}")
	echo "${driver}"`
)

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
	K8sVersion      string
	CriCGroupDriver string
}

func (k *Kubeadm) Template() string {
	return kubeadmYamlTemplate
}

func (k *Kubeadm) TemplateConvert() string {
	p := map[string]interface{}{
		"Version": k.K8sVersion,
	}
	criScoket, adminAPI := setKubeadmAPIByVersion(k.K8sVersion)
	p[KubeadmAPI] = adminAPI
	p[CriSocket] = criScoket
	// we need to Dynamic get cgroup driver on ervery join nodes.
	p[CriCGroupDriver] = k.CriCGroupDriver

	data, err := templateFromContent(k.Template(), p)
	if err != nil {
		logger.Error(err) //nolint:typecheck
		return ""
	}
	return data
}

const (
	V1991                      = "v1.19.1"
	V1992                      = "v1.19.2"
	V1150                      = "v1.15.0"
	V1200                      = "v1.20.0"
	V1230                      = "v1.23.0"
	DefaultDockerCRISocket     = "/var/run/dockershim.sock"
	DefaultContainerdCRISocket = "/run/containerd/containerd.sock"
	DefaultSystemdCgroupDriver = "systemd"
	DefaultCgroupDriver        = "cgroupfs"

	// kubeadm api version
	KubeadmV1beta1 = "kubeadm.k8s.io/v1beta1"
	KubeadmV1beta2 = "kubeadm.k8s.io/v1beta2"
	KubeadmV1beta3 = "kubeadm.k8s.io/v1beta3"

	KubeadmAPI      = "KubeadmAPI"
	CriSocket       = "CriSocket"
	CriCGroupDriver = "CriCGroupDriver"
)

func setKubeadmAPIByVersion(k8sVersion string) (criSocket, kubeadmAPI string) {
	switch {
	case VersionCompare(k8sVersion, V1150) && !VersionCompare(k8sVersion, V1200):
		criSocket = DefaultDockerCRISocket
		kubeadmAPI = KubeadmV1beta2
	// kubernetes gt 1.20, use Containerd instead of docker
	case VersionCompare(k8sVersion, V1200):
		kubeadmAPI = KubeadmV1beta2
		criSocket = DefaultContainerdCRISocket
	default:
		// Compatible with versions 1.14 and 1.13. but do not recommended.
		kubeadmAPI = KubeadmV1beta1
		criSocket = DefaultDockerCRISocket
	}
	return
}
func VersionCompare(v1, v2 string) bool {
	v1 = strings.Replace(v1, "v", "", -1)
	v2 = strings.Replace(v2, "v", "", -1)
	v1 = strings.Split(v1, "-")[0]
	v2 = strings.Split(v2, "-")[0]
	v1List := strings.Split(v1, ".")
	v2List := strings.Split(v2, ".")

	if len(v1List) != 3 || len(v2List) != 3 {
		logger.Error("error version format %s %s", v1, v2)
		return false
	}
	if v1List[0] > v2List[0] {
		return true
	} else if v1List[0] < v2List[0] {
		return false
	}
	if v1List[1] > v2List[1] {
		return true
	} else if v1List[1] < v2List[1] {
		return false
	}
	if v1List[2] > v2List[2] {
		return true
	}
	return true
}
