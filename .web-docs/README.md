The Tencent Cloud plugin provides the capability to build customized images based on an existing base images.

### Installation
To install this plugin add this code into your Packer configuration and run [packer init](/packer/docs/commands/init)
```hcl
packer {
  required_plugins {
    tencentcloud = {
      version = "~> 1"
      source  = "github.com/hashicorp/tencentcloud"
    }
  }
}
```

Alternatively, you can use `packer plugins install` to manage installation of this plugin.

```sh
packer plugins install github.com/hashicorp/ansible
```

### Components

#### Builders
- [tencentcloud-cvm](/packer/integrations/hashicorp/tencentcloud/latest/components/builder/cvm) - The `tencentcloud-cvm` builder plugin provides the capability to build customized images based on an existing base images.
