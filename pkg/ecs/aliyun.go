package ecs

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/utils"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/logger"
	cutils "github.com/sealyun/cloud-kernel-rootfs/pkg/utils"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/vars"
	"strconv"
	"sync"
)

func (a *AliyunEcs) getClient() *ecs.Client {
	a.ecsOnce.Do(func() {
		var err error
		a.ecsHKCli, err = ecs.NewClientWithAccessKey("", vars.AkID, vars.AkSK)
		if err != nil {
			_ = cutils.ProcessCloudError(err)
		}
	})
	return a.ecsHKCli
}

type AliyunEcs struct {
	ecsOnce  sync.Once
	ecsHKCli *ecs.Client
}

func (a *AliyunEcs) Healthy() error {
	cli, err := ecs.NewClientWithAccessKey("", vars.AkID, vars.AkSK)
	if err != nil {
		return err
	}
	r := ecs.CreateDescribeZonesRequest()
	r.RegionId = ALIRegionId
	_, err = cli.DescribeZones(r)
	if err != nil {
		return errors.New("阿里云 " + err.Error())
	}
	return nil
}

func (a *AliyunEcs) New(amount int, dryRun bool, bandwidthOut bool) []string {
	loadVars()
	client := a.getClient()
	// 创建请求并设置参数
	hk := ecs.CreateRunInstancesRequest()
	hk.ImageId = ALIImageId
	hk.InstanceType = ALIInstanceType
	hk.InternetChargeType = ALIInternetChargeType
	hk.InternetMaxBandwidthIn = ALIInternetMaxBandwidthIn
	hk.InternetMaxBandwidthOut = ALIInternetMaxBandwidthOut
	hk.InstanceChargeType = ALIInstanceChargeType
	hk.SpotStrategy = ALISpotStrategy
	hk.RegionId = ALIRegionId
	hk.SecurityGroupId = ALISecurityGroupId
	hk.VSwitchId = ALIVSwitchId
	hk.ZoneId = ALIZoneId
	hk.Password = vars.EcsPassword
	hk.Amount = requests.Integer(strconv.Itoa(amount))
	hk.ClientToken = utils.GetUUID()
	if !bandwidthOut {
		hk.InternetMaxBandwidthOut = "0"
	}
	hk.DryRun = requests.Boolean(strconv.FormatBool(dryRun))
	response, err := client.RunInstances(hk)
	if err != nil {
		_ = cutils.ProcessCloudError(err)
		return nil
	}
	return response.InstanceIdSets.InstanceIdSet
}

func (a *AliyunEcs) Delete(instanceId []string, maxCount int) {
	client := a.getClient()
	// 创建请求并设置参数
	request := ecs.CreateDeleteInstancesRequest()
	request.DryRun = requests.Boolean(strconv.FormatBool(false))
	request.Force = "true"
	request.RegionId = "cn-hongkong"
	request.InstanceId = &instanceId
	var response *ecs.DeleteInstancesResponse
	var err error
	for i := 0; i < maxCount; i++ {
		logger.Info("递归删除ecs")
		response, err = client.DeleteInstances(request)
		if err != nil {
			_ = cutils.ProcessCloudError(err)
		} else {
			break
		}
	}
	if err == nil {
		logger.Info("删除ecs成功: %s", response.RequestId)
	} else {
		logger.Error("删除ecs失败: %v", instanceId)
	}
}

func (a *AliyunEcs) Describe(instanceId string) (*CloudInstanceResponse, error) {
	client := a.getClient()
	request := ecs.CreateDescribeInstanceAttributeRequest()
	request.RegionId = "cn-hongkong"
	request.InstanceId = instanceId
	attr, err := client.DescribeInstanceAttribute(request)
	if err != nil {
		return nil, err
	}
	iResponse := &CloudInstanceResponse{
		IsOk: attr.Status == "Running",
	}
	if len(attr.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
		iResponse.PrivateIP = attr.VpcAttributes.PrivateIpAddress.IpAddress[0]
	}
	if len(attr.PublicIpAddress.IpAddress) > 0 {
		iResponse.PublicIP = attr.PublicIpAddress.IpAddress[0]
	}
	return iResponse, nil
}
