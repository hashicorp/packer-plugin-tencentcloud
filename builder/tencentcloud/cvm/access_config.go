//go:generate packer-sdc struct-markdown

package cvm

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

type Region string

// below would be moved to tencentcloud sdk git repo
const (
	Bangkok       = Region("ap-bangkok")
	Beijing       = Region("ap-beijing")
	Chengdu       = Region("ap-chengdu")
	Chongqing     = Region("ap-chongqing")
	Guangzhou     = Region("ap-guangzhou")
	GuangzhouOpen = Region("ap-guangzhou-open")
	Hongkong      = Region("ap-hongkong")
	Jakarta       = Region("ap-jakarta")
	Mumbai        = Region("ap-mumbai")
	Seoul         = Region("ap-seoul")
	Shanghai      = Region("ap-shanghai")
	Nanjing       = Region("ap-nanjing")
	ShanghaiFsi   = Region("ap-shanghai-fsi")
	ShenzhenFsi   = Region("ap-shenzhen-fsi")
	Singapore     = Region("ap-singapore")
	Tokyo         = Region("ap-tokyo")
	Frankfurt     = Region("eu-frankfurt")
	Moscow        = Region("eu-moscow")
	Ashburn       = Region("na-ashburn")
	Siliconvalley = Region("na-siliconvalley")
	Toronto       = Region("na-toronto")
	SaoPaulo      = Region("sa-saopaulo")
)

var ValidRegions = []Region{
	Bangkok, Beijing, Chengdu, Chongqing, Guangzhou, GuangzhouOpen, Hongkong, Jakarta, Shanghai, Nanjing,
	ShanghaiFsi, ShenzhenFsi,
	Mumbai, Seoul, Singapore, Tokyo, Moscow,
	Frankfurt, Ashburn, Siliconvalley, Toronto, SaoPaulo,
}

type TencentCloudAccessConfig struct {
	// Tencentcloud secret id. You should set it directly,
	// or set the TENCENTCLOUD_SECRET_ID environment variable.
	SecretId string `mapstructure:"secret_id" required:"true"`
	// Tencentcloud secret key. You should set it directly,
	// or set the TENCENTCLOUD_SECRET_KEY environment variable.
	SecretKey string `mapstructure:"secret_key" required:"true"`
	// The region where your cvm will be launch. You should
	// reference Region and Zone
	//  for parameter taking.
	Region string `mapstructure:"region" required:"true"`
	// The zone where your cvm will be launch. You should
	// reference Region and Zone
	//  for parameter taking.
	Zone string `mapstructure:"zone" required:"true"`
	// Do not check region and zone when validate.
	SkipValidation bool `mapstructure:"skip_region_validation" required:"false"`
	// The endpoint you want to reach the cloud endpoint,
	// if tce cloud you should set a tce cvm endpoint.
	CvmEndpoint string `mapstructure:"cvm_endpoint" required:"false"`
	// The endpoint you want to reach the cloud endpoint,
	// if tce cloud you should set a tce vpc endpoint.
	VpcEndpoint string `mapstructure:"vpc_endpoint" required:"false"`
}

func (cf *TencentCloudAccessConfig) Client() (*cvm.Client, *vpc.Client, error) {
	var (
		err        error
		cvm_client *cvm.Client
		vpc_client *vpc.Client
		resp       *cvm.DescribeZonesResponse
	)

	if err = cf.validateRegion(); err != nil {
		return nil, nil, err
	}

	if cf.Zone == "" {
		return nil, nil, fmt.Errorf("parameter zone must be set")
	}

	if cvm_client, err = NewCvmClient(cf.SecretId, cf.SecretKey, cf.Region, cf.CvmEndpoint); err != nil {
		return nil, nil, err
	}

	if vpc_client, err = NewVpcClient(cf.SecretId, cf.SecretKey, cf.Region, cf.VpcEndpoint); err != nil {
		return nil, nil, err
	}

	ctx := context.TODO()
	err = Retry(ctx, func(ctx context.Context) error {
		var e error
		resp, e = cvm_client.DescribeZones(nil)
		return e
	})
	if err != nil {
		return nil, nil, err
	}

	for _, zone := range resp.Response.ZoneSet {
		if cf.Zone == *zone.Zone {
			return cvm_client, vpc_client, nil
		}
	}

	return nil, nil, fmt.Errorf("unknown zone: %s", cf.Zone)
}

func (cf *TencentCloudAccessConfig) Prepare(ctx *interpolate.Context) []error {
	var errs []error

	if err := cf.Config(); err != nil {
		errs = append(errs, err)
	}

	if (cf.CvmEndpoint != "" && cf.VpcEndpoint == "") ||
		(cf.CvmEndpoint == "" && cf.VpcEndpoint != "") {
		errs = append(errs, fmt.Errorf("parameter cvm_endpoint and vpc_endpoint must be set simultaneously"))
	}

	if cf.Region == "" {
		errs = append(errs, fmt.Errorf("parameter region must be set"))
	} else if !cf.SkipValidation {
		if err := cf.validateRegion(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (cf *TencentCloudAccessConfig) Config() error {
	if cf.SecretId == "" {
		cf.SecretId = os.Getenv("TENCENTCLOUD_SECRET_ID")
	}

	if cf.SecretKey == "" {
		cf.SecretKey = os.Getenv("TENCENTCLOUD_SECRET_KEY")
	}

	if cf.SecretId == "" || cf.SecretKey == "" {
		return fmt.Errorf("parameter secret_id and secret_key must be set")
	}

	return nil
}

func (cf *TencentCloudAccessConfig) validateRegion() error {
	// if set cvm endpoint, do not validate region
	if cf.CvmEndpoint != "" {
		return nil
	}
	return validRegion(cf.Region)
}

func validRegion(region string) error {
	for _, valid := range ValidRegions {
		if Region(region) == valid {
			return nil
		}
	}

	return fmt.Errorf("unknown region: %s", region)
}
