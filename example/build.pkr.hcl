# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "secret_id" {
  type    = string
  default = "${env("TENCENTCLOUD_SECRET_ID")}"
}

variable "secret_key" {
  type    = string
  default = "${env("TENCENTCLOUD_SECRET_KEY")}"
}

source "tencentcloud-cvm" "example" {
  disk_type     = "CLOUD_PREMIUM"
  image_name    = "PackerTest"
  instance_type = "SA2.MEDIUM2"
  packer_debug  = true
  region        = "ap-guangzhou"
  run_tags      = {
    good = "luck"
  }
  secret_id                   = "${var.secret_id}"
  secret_key                  = "${var.secret_key}"
  source_image_id             = "img-9qrfy1xt"
  ssh_username                = "root"
  zone                        = "ap-guangzhou-3"
  security_group_id           = "sg-r7kju7cf"
  associate_public_ip_address = true
  image_tags                  = {
    createdBy = "packer"
    usedBy    = "test"
  }
}

build {
  sources = ["source.tencentcloud-cvm.example"]

  provisioner "shell" {
    inline = ["sleep 30", "yum install redis.x86_64 -y"]
  }
}
