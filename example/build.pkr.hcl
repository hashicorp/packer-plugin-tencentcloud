variable "secret_id" {
  type    = string
  default = "${env("TENCENTCLOUD_SECRET_ID")}"
}

variable "secret_key" {
  type    = string
  default = "${env("TENCENTCLOUD_SECRET_KEY")}"
}

source "tencentcloud-cvm" "example" {
  associate_public_ip_address = true
  disk_type                   = "CLOUD_PREMIUM"
  image_name                  = "PackerTest"
  instance_type               = "S4.SMALL1"
  packer_debug                = true
  region                      = "ap-guangzhou"
  run_tags = {
    good = "luck"
  }
  secret_id       = "${var.secret_id}"
  secret_key      = "${var.secret_key}"
  source_image_id = "img-oikl1tzv"
  ssh_username    = "root"
  zone            = "ap-guangzhou-4"
}

build {
  sources = ["source.tencentcloud-cvm.example"]

  provisioner "shell" {
    inline = ["sleep 30", "yum install redis.x86_64 -y"]
  }
}
