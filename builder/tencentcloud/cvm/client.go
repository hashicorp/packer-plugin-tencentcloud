// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cvm

import (
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	org "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
	sts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

type TencentCloudClient struct {
	Credential    *common.Credential
	ClientProfile *profile.ClientProfile

	Region string

	vpcConn *vpc.Client
	cvmConn *cvm.Client
	stsConn *sts.Client
	orgConn *org.Client
	camConn *cam.Client
}

func (me *TencentCloudClient) UseVpcClient(cpf *profile.ClientProfile) *vpc.Client {
	if me.vpcConn != nil {
		return me.vpcConn
	}

	me.vpcConn, _ = vpc.NewClient(me.Credential, me.Region, cpf)
	// me.vpcConn.WithHttpTransport(&LogRoundTripper{})

	return me.vpcConn
}

func (me *TencentCloudClient) UseCvmClient(cpf *profile.ClientProfile) *cvm.Client {
	if me.cvmConn != nil {
		return me.cvmConn
	}

	me.cvmConn, _ = cvm.NewClient(me.Credential, me.Region, cpf)

	return me.cvmConn
}

func (me *TencentCloudClient) UseOrgClient(cpf *profile.ClientProfile) *org.Client {
	if me.orgConn != nil {
		return me.orgConn
	}

	me.orgConn, _ = org.NewClient(me.Credential, me.Region, cpf)

	return me.orgConn
}

func (me *TencentCloudClient) UseCamClient() *cam.Client {
	if me.camConn != nil {
		return me.camConn
	}

	cpf := me.ClientProfile
	me.camConn, _ = cam.NewClient(me.Credential, me.Region, cpf)

	return me.camConn
}

func (me *TencentCloudClient) UseStsClient() *sts.Client {
	if me.stsConn != nil {
		return me.stsConn
	}

	cpf := me.ClientProfile
	me.stsConn, _ = sts.NewClient(me.Credential, me.Region, cpf)

	return me.stsConn
}
