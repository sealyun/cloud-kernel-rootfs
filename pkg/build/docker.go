package build

import (
	"errors"
	"fmt"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/logger"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/retry"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/sshcmd/sshutil"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/templates"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/utils"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/vars"
	"strings"
	"time"
)

//k8s dockerShell k8s
var dockerShell = `wget http://gosspublic.alicdn.com/ossutil/1.7.3/ossutil64 &&  chmod a+x ossutil64 && \
mv ossutil64 /usr/bin/ && \
yum install -y git conntrack tree && \
git clone https://github.com/alibaba/sealer && \
git clone https://github.com/sealyun/cloud-kernel-rootfs && mv cloud-kernel-rootfs cloud-kernel && \
cd cloud-kernel && git checkout %s && mkdir -p rootfs && mkdir -p rootfs/bin &&  mkdir -p rootfs/registry \
cp -rf runtime/rootfs/* rootfs/ && cp -rf runtime/docker/* rootfs/   && \
cp -rf ../sealer/rootfs/rootfs/* rootfs/ && cp -rf ../sealer/rootfs/docker/* rootfs/   && \
%s && \
%s && \
%s && \
%s && \
%s && \
cd rootfs/scripts && chmod a+x * && sh init.sh  &&\
docker pull fanux/lvscare &&  \
cp /usr/sbin/conntrack ../bin/`

var dockerSaveShell = `cd cloud-kernel/rootfs/scripts &&  \
sh save-images.sh && \
rm save-images.sh && \
cd ../ && tree -L 3  && %s`

type dockerK8s struct {
	ssh        sshutil.SSH
	publicIP   string
	k8sVersion string
}

func NewDockerK8s(publicIP string, k8sVersion string) _package {

	return &dockerK8s{
		ssh: sshutil.SSH{
			User:     "root",
			Password: vars.EcsPassword,
			Timeout:  nil,
		},
		publicIP:   publicIP,
		k8sVersion: k8sVersion,
	}
}
func (d *dockerK8s) InitK8sServer() error {
	k8s := vars.Bin.Kubernetes.FinalShell()
	cri := vars.Bin.CriCtl.FinalShell()
	nerd := vars.Bin.NerdCtl.FinalShell()
	sealUtil := vars.Bin.SealUtil.FinalShell()
	sealer := vars.Bin.Sealer.FinalShell()
	shell := fmt.Sprintf(dockerShell, vars.Branch, k8s, cri, nerd, sealUtil, sealer)
	//logger.Debug(shell)
	//os.Exit(-1)
	err := d.ssh.CmdAsync(d.publicIP, shell)
	if err != nil {
		return utils.ProcessError(err)
	}
	return nil
}
func (d *dockerK8s) RunK8sServer() error {
	version := strings.Join([]string{"v", d.k8sVersion}, "")
	writeKubeadm := `cd cloud-kernel/rootfs && echo '%s' > etc/kubeadm-config.yaml`

	//
	kubeadm := &templates.Kubeadm{K8sVersion: version}
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
func (d *dockerK8s) WaitImages() error {
	if err := d.ssh.CmdAsync(d.publicIP, "docker images"); err != nil {
		_ = utils.ProcessError(err)
		return err
	}
	err := retry.Do(func() error {
		logger.Debug(fmt.Sprintf("%d. retry wait k8s  pod is running :%s", 4, d.publicIP))
		checkShell := "docker images   | grep  \"lvscare\" | wc -l"
		podNum := d.ssh.CmdToString(d.publicIP, checkShell, "")
		if podNum == "0" {
			return errors.New("retry error")
		}
		return nil
	}, 100, 500*time.Millisecond, false)
	if err != nil {
		return utils.ProcessError(err)
	}
	return nil
}

func (d *dockerK8s) SaveImages() error {
	rootfs := vars.Bin.Rootfs
	err := d.ssh.CmdAsync(d.publicIP, fmt.Sprintf(dockerSaveShell, rootfs.FinalShell()))
	if err != nil {
		return utils.ProcessError(err)
	}

	return nil
}
