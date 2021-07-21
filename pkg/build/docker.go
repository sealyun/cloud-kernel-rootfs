package build

import (
	"errors"
	"fmt"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/logger"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/retry"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/sshcmd/sshutil"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/utils"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/vars"
	"time"
)

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
func (d *dockerK8s) InitCRI() error {
	var dockerShell = `
cd cloud-kernel && cp -rf runtime/docker/* rootfs/   && \
cp -rf ../sealer/filesystem/rootfs/docker/* rootfs/`
	err := d.ssh.CmdAsync(d.publicIP, dockerShell)
	if err != nil {
		return utils.ProcessError(err)
	}
	return nil
}
func (d *dockerK8s) PullImages() error {
	var dockerShell = `docker pull fanux/lvscare`
	err := d.ssh.CmdAsync(d.publicIP, dockerShell)
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
