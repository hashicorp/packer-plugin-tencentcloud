// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package image

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
)

//go:embed test-fixtures/template.pkr.hcl
var testDatasourceHCL2Basic string

// Run with: PACKER_ACC=1 go test -count 1 -v ./datasource/image/data_acc_test.go  -timeout=120m
func TestAccImageDatasource(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name: "image_datasource_basic_test",
		Setup: func() error {
			return nil
		},
		Teardown: func() error {
			return nil
		},
		Template: testDatasourceHCL2Basic,
		Type:     "image-my-datasource",
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}

			logs, err := os.Open(logfile)
			if err != nil {
				return fmt.Errorf("Unable find %s", logfile)
			}
			defer logs.Close()

			logsBytes, err := ioutil.ReadAll(logs)
			if err != nil {
				return fmt.Errorf("Unable to read %s", logfile)
			}
			logsString := string(logsBytes)

			idLog := "null.basic-example: id: img-39ei7bw5"
			nameLog := "null.basic-example: name: Rocky Linux 9.4 64bit"

			if matched, _ := regexp.MatchString(idLog+".*", logsString); !matched {
				t.Fatalf("logs doesn't contain expected ID value %q", logsString)
			}
			if matched, _ := regexp.MatchString(nameLog+".*", logsString); !matched {
				t.Fatalf("logs doesn't contain expected name value %q", logsString)
			}

			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}
