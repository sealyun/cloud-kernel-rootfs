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
package build

import (
	"fmt"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/logger"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/sshcmd/sshutil"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/templates"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/utils"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/vars"
	"os"
	"path"
	"strings"
)

type install struct {
	ssh        sshutil.SSH
	publicIP   string
	k8sVersion string
}

func NewInstall(publicIP string, k8sVersion string) *install {

	return &install{
		ssh: sshutil.SSH{
			User:     "root",
			Password: vars.EcsPassword,
			Timeout:  nil,
		},
		publicIP:   publicIP,
		k8sVersion: k8sVersion,
	}
}

func (i *install) pull() error {
	oss := vars.Bin.OSSUtil.FinalShell()
	pull := `%s && yum install -y git conntrack tree && \
git clone https://github.com/alibaba/sealer && \
mkdir -p cloud-kernel`
	shell := fmt.Sprintf(pull, oss)
	err := i.ssh.CmdAsync(i.publicIP, shell)
	if err != nil {
		return utils.ProcessError(err)
	}
	return nil
}

func (d *install) merge() error {
	merge := `cd cloud-kernel && mkdir -p rootfs && \
mkdir -p rootfs/bin &&  mkdir -p rootfs/registry && \
cp -rf ../sealer/filesystem/rootfs/rootfs/* rootfs/ && \
%s && \
%s && \
%s && \
%s && \
%s && \
cp /usr/sbin/conntrack rootfs/bin/`
	k8s := vars.Bin.Kubernetes.FinalShell()
	cri := vars.Bin.CriCtl.FinalShell()
	nerd := vars.Bin.NerdCtl.FinalShell()
	sealUtil := vars.Bin.SealUtil.FinalShell()
	sealer := vars.Bin.Sealer.FinalShell()
	shell := fmt.Sprintf(merge, k8s, cri, nerd, sealUtil, sealer)
	err := d.ssh.CmdAsync(d.publicIP, shell)
	if err != nil {
		return utils.ProcessError(err)
	}
	return nil
}

func (d *install) init() error {
	init := `cd cloud-kernel/rootfs/scripts && chmod a+x * && sh init.sh`
	err := d.ssh.CmdAsync(d.publicIP, init)
	if err != nil {
		return utils.ProcessError(err)
	}
	return nil
}

func (d *install) runK8sServer() error {
	version := strings.Join([]string{"v", d.k8sVersion}, "")
	writeKubeadm := `cd cloud-kernel/rootfs && echo '%s' > etc/kubeadm-config.yaml`
	//
	kubeadm := &templates.Kubeadm{K8sVersion: version}
	var cgroupShell string
	if templates.VersionCompare(version, templates.V1200) {
		cgroupShell = templates.ContainerdShell
	} else {
		cgroupShell = templates.DockerShell
	}
	cgroupDriver := d.ssh.CmdToString(d.publicIP, cgroupShell, " ")
	kubeadm.CriCGroupDriver = cgroupDriver

	writeShell := fmt.Sprintf(writeKubeadm, kubeadm.TemplateConvert())
	err := d.ssh.CmdAsync(d.publicIP, writeShell)
	if err != nil {
		return utils.ProcessError(err)
	}
	writeCluster := `cd cloud-kernel/rootfs && echo '%s' > etc/Clusterfile`
	cluster := &templates.Cluster{}
	writeShell = fmt.Sprintf(writeCluster, cluster.Template())
	err = d.ssh.CmdAsync(d.publicIP, writeShell)
	if err != nil {
		return utils.ProcessError(err)
	}
	writeShell = `cd cloud-kernel/rootfs && mkdir -p /var/lib/sealer && cp etc/Clusterfile /var/lib/sealer/`
	err = d.ssh.CmdAsync(d.publicIP, writeShell)
	if err != nil {
		return utils.ProcessError(err)
	}
	writeMetadata := `cd cloud-kernel/rootfs && echo '%s' > Metadata`
	metadata := &templates.Metadata{K8sVersion: version}
	writeShell = fmt.Sprintf(writeMetadata, metadata.TemplateConvert())
	err = d.ssh.CmdAsync(d.publicIP, writeShell)
	if err != nil {
		return utils.ProcessError(err)
	}
	//kubeadm init
	writeShell = `cd cloud-kernel/rootfs/etc && kubeadm init --config kubeadm-config.yaml && \
mkdir ~/.kube && cp /etc/kubernetes/admin.conf ~/.kube/config && \
kubectl taint nodes --all node-role.kubernetes.io/master- `
	err = d.ssh.CmdAsync(d.publicIP, writeShell)
	if err != nil {
		return utils.ProcessError(err)
	}
	return nil
}
func (d *install) initSaveImageShellFile() error {
	logger.Debug("initSaveImageShellFile init .cloud-kernel-rootfs dir for host ")
	if !utils.FileExist(path.Join(utils.GetUserHomeDir(), ".cloud-kernel-rootfs")) {
		err := os.MkdirAll(path.Join(utils.GetUserHomeDir(), ".cloud-kernel-rootfs"), os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *install) save() error {

	rootfs := vars.Bin.Rootfs
	saveShell := `cd cloud-kernel/rootfs/scripts &&  \
sh save-images.sh && \
rm save-images.sh && \
cd ../ && tree -L 3  && %s`
	err := d.ssh.CmdAsync(d.publicIP, fmt.Sprintf(saveShell, rootfs.FinalShell()))
	if err != nil {
		return utils.ProcessError(err)
	}

	return nil
}
