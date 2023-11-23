// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cvm

import (
	"testing"
)

func TestTencentCloudAccessConfig_Prepare(t *testing.T) {
	cf := TencentCloudAccessConfig{
		SecretId:  "secret-id",
		SecretKey: "secret-key",
	}

	if err := cf.Prepare(nil); err == nil {
		t.Fatal("should raise error: region not set")
	}

	cf.Region = "ap-guangzhou"
	if err := cf.Prepare(nil); err != nil {
		t.Fatalf("shouldn't raise error: %v", err)
	}

	cf.Region = "unknown-region"
	if err := cf.Prepare(nil); err == nil {
		t.Fatal("should raise error: unknown region")
	}

	cf.skipValidation = true
	if err := cf.Prepare(nil); err != nil {
		t.Fatalf("shouldn't raise error: %v", err)
	}
}

func TestTencentCloudAccessConfig_Profile(t *testing.T) {
	cf := TencentCloudAccessConfig{
		Region: "ap-guangzhou",
	}

	err := cf.Prepare(nil)
	if err == nil {
		t.Fatal("should raise error: aksk not set")
	}

	if cf.Profile != DEFAULT_PROFILE {
		t.Fatal("should raise error: default value should be set")
	}
}


