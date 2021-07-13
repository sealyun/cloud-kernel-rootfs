package vars

import (
	"github.com/sealyun/cloud-kernel-rootfs/pkg/vars/multiplatform"
	"os"
)

type DownloadBin struct {
	CriCtl         multiplatform.Vars
	Kubernetes     multiplatform.Vars
	MarketCtl      multiplatform.Vars
	NerdCtl        multiplatform.Vars
	SealyunWebsite multiplatform.Vars
	SealUtil       multiplatform.Vars
	Sealer         multiplatform.Vars
	SSHCmd         multiplatform.Vars
	Rootfs         multiplatform.Vars
}

var (
	DingDing         string
	AkID             string
	AkSK             string
	RegistryUserName string
	RegistryPassword string
	MarketCtlToken   string
	Platform         string
	Bin              DownloadBin
	DefaultPrice     float64
	DefaultZeroPrice float64
	DefaultClass     = "cloud_kernel" //cloud_kernel

	defaultSealVersion      = "0.2.1"
	defaultMarketCtlVersion = "1.0.5" //v1.0.5
	defaultSSHCmdVersion    = "1.5.5"
	defaultNerdctlVersion   = "0.7.3"
	defaultCriCtlVersion    = "1.20.0"
)

const (
	EcsPassword = "Fanux#123"
	Branch      = "cluster-image-cri"
)

func loadEnv() {
	if v := os.Getenv("SEALER_VERSION"); v != "" {
		defaultSealVersion = v
	}
	if v := os.Getenv("MARKET_CTL_VERSION"); v != "" {
		defaultMarketCtlVersion = v
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

func LoadVars(k8sVersion string) error {
	loadEnv()
	p := multiplatform.Platform(Platform)
	rootfs := "rootfs"
	Bin = DownloadBin{
		CriCtl:         multiplatform.NewCriCTL(defaultCriCtlVersion, rootfs, p),
		Kubernetes:     multiplatform.NewKubernetes(k8sVersion, rootfs, p),
		Rootfs:         multiplatform.NewRootfs(k8sVersion, p),
		MarketCtl:      multiplatform.NewMarketctl(defaultMarketCtlVersion, rootfs, p),
		NerdCtl:        multiplatform.NewNerdctl(defaultNerdctlVersion, rootfs, p),
		SealyunWebsite: multiplatform.NewSealyunWebsite(p),
		SealUtil:       multiplatform.NewSeautil(defaultSealVersion, rootfs, p),
		Sealer:         multiplatform.NewSealer(defaultSealVersion, p),
		SSHCmd:         multiplatform.NewSSHCmd(defaultSSHCmdVersion, rootfs, p),
	}
	return nil
}

const MarketYaml = `
market:
  body:
    spec:
      name: v%s
      price: %.2f
      product:
        class: %s
        productName: %s
      url: /tmp/kube%s.tar.gz
    status:
      productVersionStatus: ONLINE
  kind: productVersion`
