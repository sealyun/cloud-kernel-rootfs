package ecs

import (
	"github.com/sealyun/cloud-kernel-rootfs/pkg/vars"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/vars/multiplatform"
)

type Cloud interface {
	New(amount int, dryRun bool, bandwidthOut bool) []string
	Delete(instanceId []string, maxCount int)
	Describe(instanceId string) (*CloudInstanceResponse, error)
	Healthy() error
}
type CloudInstanceResponse struct {
	IsOk      bool
	PrivateIP string
	PublicIP  string
}

func NewCloud() Cloud {
	var c Cloud
	if multiplatform.Platform(vars.Platform) == multiplatform.LinuxAmd64 {
		c = &AliyunEcs{}
	} else {
		c = &HuaweiEcs{}
	}
	return c
}
