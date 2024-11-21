// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cvm

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
)

type stepShareImage struct {
	ShareAccounts     []string
	IsShareOrgMembers bool
}

func (s *stepShareImage) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	if len(s.ShareAccounts) == 0 && !s.IsShareOrgMembers {
		return multistep.ActionContinue
	}

	client := state.Get("cvm_client").(*cvm.Client)

	imageId := state.Get("image").(*cvm.Image).ImageId
	Say(state, strings.Join(s.ShareAccounts, ","), "Trying to share image to")

	req := cvm.NewModifyImageSharePermissionRequest()
	req.ImageId = imageId
	req.Permission = common.StringPtr("SHARE")
	accounts := []*string{}
	for _, account := range s.ShareAccounts {
		accounts = append(accounts, common.StringPtr(account))
	}

	if s.IsShareOrgMembers {
		accountList, err := s.getOrgAccounts(ctx, state)
		if err != nil {
			return Halt(state, err, "Failed to get org accounts")
		}
		accounts = append(accounts, accountList...)
	}

	if len(accounts) == 0 {
		return multistep.ActionContinue
	}

	req.AccountIds = accounts
	err := Retry(ctx, func(ctx context.Context) error {
		_, e := client.ModifyImageSharePermission(req)
		return e
	})
	if err != nil {
		return Halt(state, err, "Failed to share image")
	}

	Message(state, "Image shared", "")

	return multistep.ActionContinue
}

func (s *stepShareImage) getOrgAccounts(ctx context.Context, state multistep.StateBag) ([]*string, error) {

	currentAccount, err := s.getUserId(ctx, state)
	if err != nil {
		return nil, err
	}

	req := organization.NewDescribeOrganizationMembersRequest()
	resp := organization.NewDescribeOrganizationMembersResponse()

	var limit uint64 = 50
	var offset uint64 = 0

	req.Limit = &limit
	req.Offset = &offset

	accounts := []*string{}
	for {
		client := state.Get("org_client").(*organization.Client)
		err := Retry(ctx, func(ctx context.Context) error {
			var e error
			resp, e = client.DescribeOrganizationMembers(req)
			return e
		})
		if err != nil {
			return nil, nil
		}
		if resp.Response == nil {
			return nil, nil
		}
		items := resp.Response.Items
		for _, v := range items {
			if v.MemberUin != nil {
				if strconv.FormatInt(*v.MemberUin, 10) == currentAccount {
					continue
				}
				accounts = append(accounts, common.StringPtr(strconv.Itoa(int(*v.MemberUin))))
			}
		}

		if len(items) < int(limit) {
			break
		}

		offset += limit
	}

	return accounts, nil
}

func (s *stepShareImage) getUserId(ctx context.Context, state multistep.StateBag) (string, error) {
	req := cam.NewGetUserAppIdRequest()
	resp := cam.NewGetUserAppIdResponse()

	client := state.Get("cam_client").(*cam.Client)
	err := Retry(ctx, func(ctx context.Context) error {
		var e error
		resp, e = client.GetUserAppId(req)
		return e
	})
	if err != nil {
		return "", err
	}

	if resp.Response == nil {
		return "", nil
	}

	if resp.Response.Uin != nil {
		return *resp.Response.Uin, nil
	}

	return "", nil
}

func (s *stepShareImage) Cleanup(state multistep.StateBag) {
	_, cancelled := state.GetOk(multistep.StateCancelled)
	_, halted := state.GetOk(multistep.StateHalted)
	if !cancelled && !halted {
		return
	}

	ctx := context.TODO()
	client := state.Get("cvm_client").(*cvm.Client)

	imageId := state.Get("image").(*cvm.Image).ImageId
	SayClean(state, "image share")

	req := cvm.NewModifyImageSharePermissionRequest()
	req.ImageId = imageId
	req.Permission = common.StringPtr("CANCEL")
	accounts := make([]*string, 0, len(s.ShareAccounts))
	for _, account := range s.ShareAccounts {
		account := account
		accounts = append(accounts, &account)
	}
	req.AccountIds = accounts
	err := Retry(ctx, func(ctx context.Context) error {
		_, e := client.ModifyImageSharePermission(req)
		return e
	})
	if err != nil {
		Error(state, err, fmt.Sprintf("Failed to cancel share image(%s), please delete it manually", *imageId))
	}
}
