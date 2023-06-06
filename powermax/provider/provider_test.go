/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"fmt"
	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"powermax": providerserver.NewProtocol6WithError(New("test")()),
}

var ProviderConfig = ""
var FunctionMocker *Mocker

// for acc test, avoid conflict of existing resources.
var ResourceSuffix = RandResNameSuffix(5)

func init() {
	err := godotenv.Load("powermax.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return
	}

	username := os.Getenv("POWERMAX_USERNAME")
	password := os.Getenv("POWERMAX_PASSWORD")
	endpoint := os.Getenv("POWERMAX_ENDPOINT")
	serialNumber := os.Getenv("POWERMAX_SERIAL_NUMBER")
	pmaxVersion := os.Getenv("POWERMAX_VERSION")

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
	if os.Getenv("POWERMAX_ENDPOINT") == "" {
		t.Fatal("POWERMAX_ENDPOINT environment variable not set")
	}
	if os.Getenv("POWERMAX_USERNAME") == "" {
		t.Fatal("POWERMAX_USERNAME environment variable not set")
	}
	if os.Getenv("POWERMAX_PASSWORD") == "" {
		t.Fatal("POWERMAX_PASSWORD environment variable not set")
	}
	if os.Getenv("POWERMAX_SERIAL_NUMBER") == "" {
		t.Fatal("POWERMAX_SERIAL_NUMBER environment variable not set")
	}
	if os.Getenv("POWERMAX_VERSION") == "" {
		t.Fatal("POWERMAX_VERSION environment variable not set")
	}

	t.Log(ProviderConfig)
}

func RandResNameSuffix(length int) string {
	const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())
	// generate arr of bytes for ascii characters
	b := make([]byte, length)
	for i := range b {
		b[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(b)
}
