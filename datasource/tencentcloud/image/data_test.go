package image

import (
	cvm "github.com/hashicorp/packer-plugin-tencentcloud/builder/tencentcloud/cvm"
	"testing"
)

var tencentCloudAccessConfig = cvm.TencentCloudAccessConfig{
	Region:    "na-ashburn",
	SecretId:  "secret",
	SecretKey: "key",
}

func TestDatasourceConfigure_NoOptionsSpecified(t *testing.T) {
	ds := Datasource{
		config: Config{
			TencentCloudAccessConfig: tencentCloudAccessConfig,
			ImageFilterOptions: ImageFilterOptions{
				Filters:     map[string]string{},
				ImageFamily: "",
				MostRecent:  false,
			},
		},
	}

	if err := ds.Configure(); err == nil {
		t.Fatal("Should fail since at least one option must be specified")
	} else {
		t.Log(err)
	}
}

func TestDatasourceConfigure_BothFiltersAndImageFamilySpecified(t *testing.T) {
	ds := Datasource{
		config: Config{
			TencentCloudAccessConfig: tencentCloudAccessConfig,
			ImageFilterOptions: ImageFilterOptions{
				Filters: map[string]string{
					"foo": "bar",
				},
				ImageFamily: "foo",
				MostRecent:  false,
			},
		},
	}

	if err := ds.Configure(); err == nil {
		t.Fatal("Should fail since options are mutually exclusive")
	} else {
		t.Log(err)
	}
}

func TestDatasourceConfigure_FiltersSpecified(t *testing.T) {
	ds := Datasource{
		config: Config{
			TencentCloudAccessConfig: tencentCloudAccessConfig,
			ImageFilterOptions: ImageFilterOptions{
				Filters: map[string]string{
					"foo": "bar",
				},
				ImageFamily: "",
				MostRecent:  false,
			},
		},
	}

	if err := ds.Configure(); err != nil {
		t.Fatal("Should not fail")
	}
}

func TestDatasourceConfigure_ImageFamilySpecified(t *testing.T) {
	ds := Datasource{
		config: Config{
			TencentCloudAccessConfig: tencentCloudAccessConfig,
			ImageFilterOptions: ImageFilterOptions{
				Filters:     map[string]string{},
				ImageFamily: "foo",
				MostRecent:  false,
			},
		},
	}

	if err := ds.Configure(); err != nil {
		t.Fatal("Should not fail")
	}
}
