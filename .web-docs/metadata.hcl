# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name = "Tencent Cloud"
  description = "The Tencent Cloud plugin provides the capability to build customized images based on an existing base images."
  identifier = "packer/hashicorp/tencentcloud"
  component {
    type = "builder"
    name = "Tencent Cloud Builder"
    slug = "cvm"
  }
}
