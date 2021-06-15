module github.com/hashicorp/packer-plugin-tencentcloud

go 1.16

require (
	github.com/hashicorp/hcl/v2 v2.10.0
	github.com/hashicorp/packer-plugin-sdk v0.2.3
	github.com/pkg/errors v0.9.1
	github.com/tencentcloud/tencentcloud-sdk-go v1.0.152
	github.com/zclconf/go-cty v1.8.3
)

// This version contained an invalid version for github.com/tencentcloud/tencentcloud-sdk-go
retract v1.0.0
