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
	"github.com/sealyun/cloud-kernel-rootfs/pkg/sshcmd/sshutil"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/templates"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/utils"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/vars"
	"strings"
)

type upload struct {
	ssh        sshutil.SSH
	publicIP   string
	k8sVersion string
}

func NewUpload(publicIP string, k8sVersion string) *upload {

	return &upload{
		ssh: sshutil.SSH{
			User:     "root",
			Password: vars.EcsPassword,
			Timeout:  nil,
		},
		publicIP:   publicIP,
		k8sVersion: k8sVersion,
	}
}

func (d *upload) Upload() error {
	rootfs := vars.Bin.Rootfs
	imageName := rootfs.FetchWgetURL()
	if vars.OSSAkID != "" && vars.OSSAkSK != "" {
		writeCalico := `cd cloud-kernel  && echo '%s' >  oss-config`
		ossConfig := &templates.OSSConfig{
			KeyId:     vars.OSSAkID,
			KeySecret: vars.OSSAkSK,
		}
		writeShell := fmt.Sprintf(writeCalico, ossConfig.TemplateConvert())
		err := d.ssh.CmdAsync(d.publicIP, writeShell)
		if err != nil {
			return utils.ProcessError(err)
		}
		//kubernetes-amd64:v1.15 > kube-amd64-v1.15.tar
		tarName := strings.ReplaceAll(imageName, "kubernetes", "kube")
		tarName = strings.ReplaceAll(imageName, ":", "-")
		tarName = tarName + ".tar"
		tarShell := `cd cloud-kernel && sealer save -o  %s %s && \
ossutil64 -c  oss-config cp -f %s oss://%s/%s`
		err = d.ssh.CmdAsync(d.publicIP, fmt.Sprintf(tarShell, tarName, imageName, tarName, vars.OSSRepo, tarName))
		if err != nil {
			return utils.ProcessError(err)
		}
	}
	pushShell := `sealer login  %s -u %s -p %s && \
sealer tag  %s  %s && \
sealer push %s`
	addr := vars.RegistryAddress

	registryImage := fmt.Sprintf("%s/%s/%s", vars.RegistryAddress, vars.RegistryRepo, imageName)

	err := d.ssh.CmdAsync(d.publicIP, fmt.Sprintf(pushShell, addr, vars.RegistryUserName, vars.RegistryPassword, imageName, registryImage, registryImage))
	if err != nil {
		return utils.ProcessError(err)
	}
	return nil
}
