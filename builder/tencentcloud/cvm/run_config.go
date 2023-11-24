// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type tencentCloudDataDisk

package cvm

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/hashicorp/packer-plugin-sdk/uuid"
	"github.com/pkg/errors"
)

type tencentCloudDataDisk struct {
	DiskType   string `mapstructure:"disk_type"`
	DiskSize   int64  `mapstructure:"disk_size"`
	SnapshotId string `mapstructure:"disk_snapshot_id"`
}

type TencentCloudRunConfig struct {
	// Whether allocate public ip to your cvm.
	// Default value is false.
	AssociatePublicIpAddress bool `mapstructure:"associate_public_ip_address" required:"false"`
	// The base image id of Image you want to create
	// your customized image from.
	SourceImageId string `mapstructure:"source_image_id" required:"false"`
	// The base image name of Image you want to create your
	// customized image from.Conflict with SourceImageId.
	SourceImageName string `mapstructure:"source_image_name" required:"false"`
	// Charge type of cvm, values can be `POSTPAID_BY_HOUR` (default) `SPOTPAID`
	InstanceChargeType string `mapstructure:"instance_charge_type" required:"false"`
	// The instance type your cvm will be launched by.
	// You should reference [Instance Type](https://intl.cloud.tencent.com/document/product/213/11518)
	// for parameter taking.
	InstanceType string `mapstructure:"instance_type" required:"true"`
	// Instance name.
	InstanceName string `mapstructure:"instance_name" required:"false"`
	// Root disk type your cvm will be launched by, default is `CLOUD_PREMIUM`. you could
	// reference Disk Type
	// for parameter taking.
	DiskType string `mapstructure:"disk_type" required:"false"`
	// Root disk size your cvm will be launched by. values range(in GB):
	// - LOCAL_BASIC: 50
	// - Other: 50 ~ 1000 (need whitelist if > 50)
	DiskSize int64 `mapstructure:"disk_size" required:"false"`
	// Add one or more data disks to the instance before creating the image.
	// Note that if the source image has data disk snapshots, this argument
	// will be ignored, and the running instance will use source image data
	// disk settings, in such case, `disk_type` argument will be used as disk
	// type for all data disks, and each data disk size will use the origin
	// value in source image.
	// The data disks allow for the following argument:
	// -  `disk_type` - Type of the data disk. Valid choices: `CLOUD_BASIC`, `CLOUD_PREMIUM` and `CLOUD_SSD`.
	// -  `disk_size` - Size of the data disk.
	// -  `disk_snapshot_id` - Id of the snapshot for a data disk.
	DataDisks []tencentCloudDataDisk `mapstructure:"data_disks"`
	// Specify vpc your cvm will be launched by.
	VpcId string `mapstructure:"vpc_id" required:"false"`
	// Specify vpc name you will create. if vpc_id is not set, packer will
	// create a vpc for you named this parameter.
	VpcName string `mapstructure:"vpc_name" required:"false"`
	VpcIp   string `mapstructure:"vpc_ip"`
	// Specify subnet your cvm will be launched by.
	SubnetId string `mapstructure:"subnet_id" required:"false"`
	// Specify subnet name you will create. if subnet_id is not set, packer will
	// create a subnet for you named this parameter.
	SubnetName string `mapstructure:"subnet_name" required:"false"`
	// Specify cider block of the vpc you will create if vpc_id not set
	CidrBlock string `mapstructure:"cidr_block" required:"false"` // 10.0.0.0/16(default), 172.16.0.0/12, 192.168.0.0/16
	// Specify cider block of the subnet you will create if
	// subnet_id not set
	SubnectCidrBlock string `mapstructure:"subnect_cidr_block" required:"false"`
	// Internet charge type of cvm, values can be TRAFFIC_POSTPAID_BY_HOUR, BANDWIDTH_POSTPAID_BY_HOUR, BANDWIDTH_PACKAGE
	InternetChargeType string `mapstructure:"internet_charge_type" required:"false"`
	// Max bandwidth out your cvm will be launched by(in MB).
	// values can be set between 1 ~ 100.
	InternetMaxBandwidthOut int64 `mapstructure:"internet_max_bandwidth_out" required:"false"`
	// When internet_charge_type is BANDWIDTH_PACKAGE, bandwidth_package_id is required
	BandwidthPackageId string `mapstructure:"bandwidth_package_id" required:"false"`
	// Specify securitygroup your cvm will be launched by.
	SecurityGroupId string `mapstructure:"security_group_id" required:"false"`
	// Specify security name you will create if security_group_id not set.
	SecurityGroupName string `mapstructure:"security_group_name" required:"false"`
	// userdata.
	UserData string `mapstructure:"user_data" required:"false"`
	// userdata file.
	UserDataFile string `mapstructure:"user_data_file" required:"false"`
	// host name.
	HostName string `mapstructure:"host_name" required:"false"`
	// CAM role name.
	CamRoleName string `mapstructure:"cam_role_name" required:"false"`
	// Key/value pair tags to apply to the instance that is *launched* to
	// create the image. These tags are *not* applied to the resulting image.
	RunTags map[string]string `mapstructure:"run_tags" required:"false"`
	// Same as [`run_tags`](#run_tags) but defined as a singular repeatable
	// block containing a `key` and a `value` field. In HCL2 mode the
	// [`dynamic_block`](/packer/docs/templates/hcl_templates/expressions#dynamic-blocks)
	// will allow you to create those programatically.
	RunTag config.KeyValues `mapstructure:"run_tag" required:"false"`

	// Communicator settings
	Comm         communicator.Config `mapstructure:",squash"`
	SSHPrivateIp bool                `mapstructure:"ssh_private_ip"`
}

var ValidCBSType = []string{
	"LOCAL_BASIC", "LOCAL_SSD", "CLOUD_BASIC", "CLOUD_SSD", "CLOUD_PREMIUM",
}

func (cf *TencentCloudRunConfig) Prepare(ctx *interpolate.Context) []error {
	packerId := fmt.Sprintf("packer_%s", uuid.TimeOrderedUUID()[:8])
	if cf.Comm.SSHKeyPairName == "" && cf.Comm.SSHTemporaryKeyPairName == "" &&
		cf.Comm.SSHPrivateKeyFile == "" && cf.Comm.SSHPassword == "" && cf.Comm.WinRMPassword == "" {
		//tencentcloud support key pair name length max to 25
		cf.Comm.SSHTemporaryKeyPairName = packerId
	}

	errs := cf.Comm.Prepare(ctx)
	if cf.SourceImageId == "" && cf.SourceImageName == "" {
		errs = append(errs, errors.New("source_image_id or source_image_name must be specified"))
	}

	if cf.SourceImageId != "" && !CheckResourceIdFormat("img", cf.SourceImageId) {
		errs = append(errs, errors.New("source_image_id wrong format"))
	}

	if cf.InstanceType == "" {
		errs = append(errs, errors.New("instance_type must be specified"))
	}

	if cf.UserData != "" && cf.UserDataFile != "" {
		errs = append(errs, errors.New("only one of user_data or user_data_file can be specified"))
	} else if cf.UserDataFile != "" {
		if _, err := os.Stat(cf.UserDataFile); err != nil {
			errs = append(errs, errors.New("user_data_file not exist"))
		}
	}

	if (cf.VpcId != "" || cf.CidrBlock != "") && cf.SubnetId == "" && cf.SubnectCidrBlock == "" {
		errs = append(errs, errors.New("if vpc cidr_block is specified, then "+
			"subnet_cidr_block must also be specified."))
	}

	if cf.VpcId == "" {
		if cf.VpcName == "" {
			cf.VpcName = packerId
		}
		if cf.CidrBlock == "" {
			cf.CidrBlock = "10.0.0.0/16"
		}
		if cf.SubnetId != "" {
			errs = append(errs, errors.New("can't set subnet_id without set vpc_id"))
		}
	}

	if cf.SubnetId == "" {
		if cf.SubnetName == "" {
			cf.SubnetName = packerId
		}
		if cf.SubnectCidrBlock == "" {
			cf.SubnectCidrBlock = "10.0.8.0/24"
		}
	}

	if cf.SecurityGroupId == "" && cf.SecurityGroupName == "" {
		cf.SecurityGroupName = packerId
	}

	if cf.DiskType != "" && !checkDiskType(cf.DiskType) {
		errs = append(errs, errors.New(fmt.Sprintf("specified disk_type(%s) is invalid", cf.DiskType)))
	} else if cf.DiskType == "" {
		cf.DiskType = "CLOUD_PREMIUM"
	}

	if cf.DiskSize <= 0 {
		cf.DiskSize = 50
	}

	validChargeTypes := map[string]int{
		"TRAFFIC_POSTPAID_BY_HOUR":   0,
		"BANDWIDTH_POSTPAID_BY_HOUR": 0,
		"BANDWIDTH_PACKAGE":          0,
	}

	if cf.InternetChargeType != "" {
		if _, ok := validChargeTypes[cf.InternetChargeType]; !ok {
			errs = append(errs, fmt.Errorf("specified internet_charge_type(%s) is invalid.", cf.InternetChargeType))
		}
	}

	if cf.InternetChargeType == "BANDWIDTH_PACKAGE" {
		if cf.BandwidthPackageId == "" {
			errs = append(errs,
				fmt.Errorf("bandwidth_package_id is required when internet_charge_type is BANDWIDTH_PACKAGE"))
		}
	}

	if cf.AssociatePublicIpAddress && cf.InternetMaxBandwidthOut <= 0 {
		cf.InternetMaxBandwidthOut = 1
	}

	if cf.InstanceName == "" {
		cf.InstanceName = packerId
	}

	if cf.HostName == "" {
		cf.HostName = cf.InstanceName
	}

	if len(cf.HostName) > 15 {
		cf.HostName = cf.HostName[:15]
	}
	cf.HostName = strings.Replace(cf.HostName, "_", "-", -1)

	if cf.RunTags == nil {
		cf.RunTags = make(map[string]string)
	}

	errs = append(errs, cf.RunTag.CopyOn(&cf.RunTags)...)

	return errs
}

func checkDiskType(diskType string) bool {
	for _, valid := range ValidCBSType {
		if valid == diskType {
			return true
		}
	}

	return false
}
