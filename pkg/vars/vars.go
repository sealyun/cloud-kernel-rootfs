package vars

import (
	"github.com/sealyun/cloud-kernel-rootfs/pkg/logger"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/sshcmd/sshutil"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/vars/multiplatform"
	"os"
	"strings"
)

type DownloadBin struct {
	CriCtl     multiplatform.Vars
	Kubernetes multiplatform.Vars
	MarketCtl  multiplatform.Vars
	NerdCtl    multiplatform.Vars
	SealUtil   multiplatform.Vars
	Sealer     multiplatform.Vars
	SSHCmd     multiplatform.Vars
	Rootfs     multiplatform.Vars
}

var (
	DingDing         string
	AkID             string
	AkSK             string
	OSSAkID          string
	OSSAkSK          string
	OSSRepo          string
	RegistryUserName string
	RegistryAddress  string
	RegistryRepo     string
	RegistryPassword string
	Platform         string
	Bin              DownloadBin

	defaultSealVersion    = "0.2.1"
	defaultSSHCmdVersion  = "1.5.5"
	defaultNerdctlVersion = "0.7.3"
	defaultCriCtlVersion  = "1.20.0"
)

const (
	EcsPassword = "Fanux#123"
	Branch      = "main"
)

func LoadAKSK() {
	if OSSAkID == "" {
		if v := os.Getenv("OSS_AKID"); v != "" {
			OSSAkID = v
		}
	}
	if OSSAkSK == "" {
		if v := os.Getenv("OSS_AKSK"); v != "" {
			OSSAkSK = v
		}
	}
	if OSSRepo == "" {
		if v := os.Getenv("OSS_REPO"); v != "" {
			OSSRepo = v
		}
	}
	if AkID == "" {
		if v := os.Getenv("ECS_AKID"); v != "" {
			AkID = v
		}
	}
	if AkSK == "" {
		if v := os.Getenv("ECS_AKSK"); v != "" {
			AkSK = v
		}
	}
}

func loadEnv() {
	if v := os.Getenv("SEALER_VERSION"); v != "" {
		defaultSealVersion = v
	}
	if v := os.Getenv("SSH_CMD_VERSION"); v != "" {
		defaultSSHCmdVersion = v
	}
	if v := os.Getenv("NERD_CTL_VERSION"); v != "" {
		defaultNerdctlVersion = v
	}
	if v := os.Getenv("CRI_CTL_VERSION"); v != "" {
		defaultCriCtlVersion = v
	}
}

func LoadVars(k8sVersion, publicIP string, s sshutil.SSH) error {
	loadVersion(publicIP, s)
	loadEnv()
	p := multiplatform.Platform(Platform)
	rootfs := "rootfs"
	Bin = DownloadBin{
		CriCtl:     multiplatform.NewCriCTL(defaultCriCtlVersion, rootfs, p),
		Kubernetes: multiplatform.NewKubernetes(k8sVersion, rootfs, p),
		Rootfs:     multiplatform.NewRootfs(k8sVersion, p),
		NerdCtl:    multiplatform.NewNerdctl(defaultNerdctlVersion, rootfs, p),
		SealUtil:   multiplatform.NewSeautil(defaultSealVersion, rootfs, p),
		Sealer:     multiplatform.NewSealer(defaultSealVersion, p),
		SSHCmd:     multiplatform.NewSSHCmd(defaultSSHCmdVersion, rootfs, p),
	}
	return nil
}

func loadVersion(publicIP string, s sshutil.SSH) {
	//install jq
	err := s.CmdAsync(publicIP, "yum install -y jq")
	if err != nil {
		logger.Error("安装jq失败: %v", err)
		return
	}
	sealerVersion := "curl -LsSf https://api.github.com/repos/alibaba/sealer/releases/latest | jq -r \".tag_name\""
	sshcmdVersion := "curl -LsSf https://api.github.com/repos/cuisongliu/sshcmd/releases/latest | jq -r \".tag_name\""
	nerdVersion := "curl -LsSf https://api.github.com/repos/containerd/nerdctl/releases/latest | jq -r \".tag_name\""

	if version := s.CmdToString(publicIP, sealerVersion, ""); version != "" {
		defaultSealVersion = strings.ReplaceAll(version, "v", "")
		logger.Info("获取sealer最新版本: %s", defaultSealVersion)
	}
	if version := s.CmdToString(publicIP, sshcmdVersion, ""); version != "" {
		defaultSSHCmdVersion = strings.ReplaceAll(version, "v", "")
		logger.Info("获取sshcmd最新版本: %s", defaultSSHCmdVersion)
	}
	if version := s.CmdToString(publicIP, nerdVersion, ""); version != "" {
		defaultNerdctlVersion = strings.ReplaceAll(version, "v", "")
		logger.Info("获取nerdctl最新版本: %s", defaultNerdctlVersion)
	}
}
