# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name = "Tencent Cloud"
  description = "TODO"
  identifier = "packer/BrandonRomano/tencentcloud"
  component {
    type = "builder"
    name = "Tencentcloud Image Builder"
    slug = "cvm"
  }
}
