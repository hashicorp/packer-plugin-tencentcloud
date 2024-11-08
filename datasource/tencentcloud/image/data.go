// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type Config,DatasourceOutput
package image

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	builder "github.com/hashicorp/packer-plugin-tencentcloud/builder/tencentcloud/cvm"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/zclconf/go-cty/cty"
)

type ImageFilterOptions struct {
	// Filters used to select an image. Any filter described in the documentation for
	//[DescribeImages](https://www.tencentcloud.com/document/product/213/33272) can be used.
	Filters map[string]string `mapstructure:"filters"`
	// Selects the most recently created image when multiple results are returned. Note that
	// public images don't have a creation date, so this flag is only really useful for private
	// images.
	MostRecent bool `mapstructure:"most_recent"`
}

type Config struct {
	common.PackerConfig              `mapstructure:",squash"`
	builder.TencentCloudAccessConfig `mapstructure:",squash"`
	ImageFilterOptions               `mapstructure:",squash"`
}

type Datasource struct {
	config Config
}

type DatasourceOutput struct {
	// The image ID
	ID string `mapstructure:"id"`
	// The image name
	Name string `mapstructure:"name"`
}

func (d *Datasource) ConfigSpec() hcldec.ObjectSpec {
	return d.config.FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Configure(raws ...interface{}) error {
	err := config.Decode(&d.config, nil, raws...)
	if err != nil {
		return err
	}

	var errs *packersdk.MultiError
	errs = packersdk.MultiErrorAppend(errs, d.config.TencentCloudAccessConfig.Prepare()...)

	if len(d.config.Filters) == 0 {
		errs = packersdk.MultiErrorAppend(errs, fmt.Errorf("`filters` must be specified"))
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}
	return nil
}

func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Execute() (cty.Value, error) {
	client, _, err := d.config.Client()
	if err != nil {
		return cty.NullVal(cty.EmptyObject), err
	}

	req := cvm.NewDescribeImagesRequest()

	var filters []*cvm.Filter
	for k, v := range d.config.Filters {
		k := k
		v := v
		filters = append(filters, &cvm.Filter{
			Name:   &k,
			Values: []*string{&v},
		})
	}
	req.Filters = filters

	resp, err := client.DescribeImages(req)
	if err != nil {
		return cty.NullVal(cty.EmptyObject), err
	}

	if *resp.Response.TotalCount == 0 {
		return cty.NullVal(cty.EmptyObject), nil
	}

	if *resp.Response.TotalCount > 1 && !d.config.MostRecent {
		return cty.NullVal(cty.EmptyObject), fmt.Errorf("Your image query returned more than result. Please try a more specific search, or set `most_recent` to `true`.")
	}

	var image *cvm.Image

	if d.config.MostRecent {
		image = mostRecentImage(resp.Response.ImageSet)
	} else {
		image = resp.Response.ImageSet[0]
	}

	output := DatasourceOutput{
		ID:   *image.ImageId,
		Name: *image.ImageName,
	}
	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}
