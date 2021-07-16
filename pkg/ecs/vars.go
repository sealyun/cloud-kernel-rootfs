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
package ecs

//HW 代表华为参数
//ALI 代表阿里参数
var (
	HWProjectID string
	HWSubnetId  string
	HWVpcId     string
)

var (
	ALISecurityGroupId string
	ALIVSwitchId       string
)

const (
	ALIRegionId                = "cn-hongkong"
	ALIZoneId                  = "cn-hongkong-c"
	ALIImageId                 = "centos_7_04_64_20G_alibase_201701015.vhd"
	ALIInstanceType            = "ecs.c5.xlarge"
	ALIInternetChargeType      = "PayByTraffic"
	ALIInternetMaxBandwidthIn  = "100"
	ALIInternetMaxBandwidthOut = "100"
	ALIInstanceChargeType      = "PostPaid"
	ALISpotStrategy            = "SpotAsPriceGo"
	HWZone                     = "ap-southeast-1a"
	HWIpType                   = "5_bgp"
	HWImageId                  = "04678140-fcc1-465d-ba36-3a2b19d155f9"
	HWFlavorRef                = "kc1.large.2"
)
