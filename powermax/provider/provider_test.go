/*
Copyright (c) 2025 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package provider

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"powermax": providerserver.NewProtocol6WithError(New("test")()),
}

var ProviderConfig = ""
var FunctionMocker *mockey.Mocker
var globalEnvMap = getEnvMap()

func init() {
	username := globalEnvMap["POWERMAX_USERNAME"]
	password := globalEnvMap["POWERMAX_PASSWORD"]
	endpoint := globalEnvMap["POWERMAX_ENDPOINT"]
	serialNumber := globalEnvMap["POWERMAX_SERIAL_NUMBER"]
	pmaxVersion := globalEnvMap["POWERMAX_VERSION"]

	ProviderConfig = fmt.Sprintf(`
		provider "powermax" {
			username      = "%s"
			password      = "%s"
  			endpoint      = "%s"
  			serial_number = "%s"
  			pmax_version  = "%s"
  			insecure      = true
		}
	`, username, password, endpoint, serialNumber, pmaxVersion)
}

func testAccPreCheck(t *testing.T) {
	// Check that the required environment variables are set.
	if globalEnvMap["POWERMAX_ENDPOINT"] == "" {
		t.Fatal("POWERMAX_ENDPOINT environment variable not set")
	}
	if globalEnvMap["POWERMAX_USERNAME"] == "" {
		t.Fatal("POWERMAX_USERNAME environment variable not set")
	}
	if globalEnvMap["POWERMAX_PASSWORD"] == "" {
		t.Fatal("POWERMAX_PASSWORD environment variable not set")
	}
	if globalEnvMap["POWERMAX_SERIAL_NUMBER"] == "" {
		t.Fatal("POWERMAX_SERIAL_NUMBER environment variable not set")
	}
	if globalEnvMap["POWERMAX_VERSION"] == "" {
		t.Fatal("POWERMAX_VERSION environment variable not set")
	}

	t.Log(ProviderConfig)
	// Before each test clear out the mocker
	if FunctionMocker != nil {
		FunctionMocker.UnPatch()
	}
}

func getEnvMap() map[string]string {
	envMap, err := loadEnvFile("powermax.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return envMap
	}
	return envMap
}

func loadEnvFile(path string) (map[string]string, error) {
	envMap := make(map[string]string)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		envMap[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envMap, nil
}
