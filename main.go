// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/plugin"

	"github.com/hashicorp/packer-plugin-tencentcloud/builder/tencentcloud/cvm"
	"github.com/hashicorp/packer-plugin-tencentcloud/datasource/tencentcloud/image"
	"github.com/hashicorp/packer-plugin-tencentcloud/version"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder("cvm", new(cvm.Builder))
	pps.RegisterDatasource("image", new(image.Datasource))
	pps.SetVersion(version.PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
