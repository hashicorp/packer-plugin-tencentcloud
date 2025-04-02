# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

data "tencentcloud-image" "test-image" {
  filters = {
    image-type = "PUBLIC_IMAGE"
    platform = "Rocky Linux"
    image-name = "Rocky Linux 9.4 64bit"
  }
  most_recent = true
}

locals {
  id = data.tencentcloud-image.test-image.id
  name = data.tencentcloud-image.test-image.name
}

source "null" "basic-example" {
  communicator = "none"
}

build {
  sources = [
    "source.null.basic-example"
  ]

  provisioner "shell-local" {
    inline = [
      "echo id: ${local.id}",
      "echo name: ${local.name}",
    ]
  }
}
