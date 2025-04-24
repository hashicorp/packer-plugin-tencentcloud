// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc struct-markdown

package cvm

import (
	"fmt"
	"unicode/utf8"

	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type TencentCloudImageConfig struct {
	// The name you want to create your customize image,
	// it should be composed of no more than 60 characters, of letters, numbers
	// or minus sign.
	ImageName string `mapstructure:"image_name" required:"true"`
	// Image description. It should no more than 60 characters.
	ImageDescription string `mapstructure:"image_description" required:"false"`
	// Indicates whether to perform a forced shutdown to
	// create an image when soft shutdown fails. Default value is `false`.
	ForcePoweroff bool `mapstructure:"force_poweroff" required:"false"`
	// Whether enable Sysprep during creating windows image.
	Sysprep bool `mapstructure:"sysprep" required:"false"`
	// regions that will be copied to after
	// your image created.
	ImageCopyRegions []string `mapstructure:"image_copy_regions" required:"false"`
	// accounts that will be shared to
	// after your image created.
	ImageShareAccounts []string `mapstructure:"image_share_accounts" required:"false"`
	// Key/value pair tags that will be applied to the resulting image.
	ImageTags      map[string]string `mapstructure:"image_tags" required:"false"`
	skipValidation bool
	// Skip creating an image. When set to true, you don't need to enter target image information, share, copy, etc. The default value is false.
	SkipCreateImage bool `mapstructure:"skip_create_image" required:"false"`
}

func (cf *TencentCloudImageConfig) Prepare(ctx *interpolate.Context) []error {
	var errs []error

	if cf.SkipCreateImage {
		return nil
	}

	if cf.ImageName == "" {
		errs = append(errs, fmt.Errorf("image_name must be specified"))
	} else if utf8.RuneCountInString(cf.ImageName) > 60 {
		errs = append(errs, fmt.Errorf("image_name length should not exceed 60 characters"))
	}

	if utf8.RuneCountInString(cf.ImageDescription) > 60 {
		errs = append(errs, fmt.Errorf("image_description length should not exceed 60 characters"))
	}

	if len(cf.ImageCopyRegions) > 0 {
		regionSet := make(map[string]struct{})
		regions := make([]string, 0, len(cf.ImageCopyRegions))

		for _, region := range cf.ImageCopyRegions {
			if _, ok := regionSet[region]; ok {
				continue
			}

			regionSet[region] = struct{}{}

			if !cf.skipValidation {
				if err := validRegion(region); err != nil {
					errs = append(errs, err)
					continue
				}
			}
			regions = append(regions, region)
		}
		cf.ImageCopyRegions = regions
	}

	if cf.ImageTags == nil {
		cf.ImageTags = make(map[string]string)
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
