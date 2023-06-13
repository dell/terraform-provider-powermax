/*
Copyright (c) 2022-2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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

package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

	pmaxop "dell/powermax-go-client"

	pmax "github.com/dell/gopowermax/v2"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Client type is to hold powermax client and symmetrix ID.
type Client struct {
	PmaxOpenapiClient *pmaxop.APIClient
	PmaxClient        *pmax.Client
	SymmetrixID       string
}

// NewClient returns the client.
func NewClient(ctx context.Context, endpoint, username, password, serialNumber, pmaxVersion string, insecure bool) (*Client, error) {

	gomaxclient, _ := NewGopmClient(endpoint, username, password, serialNumber, pmaxVersion, insecure)
	openapiClient, _ := NewOpenApiClient(ctx, endpoint, username, password, serialNumber, pmaxVersion, insecure)
	client := Client{
		SymmetrixID:       serialNumber,
		PmaxOpenapiClient: openapiClient,
		PmaxClient:        gomaxclient,
	}
	return &client, nil
}

// NewClient returns the OpenAPI client.
func NewOpenApiClient(ctx context.Context, endpoint, username, password, serialNumber, pmaxVersion string, insecure bool) (*pmaxop.APIClient, error) {
	// Setup a User-Agent for your API client (replace the provider name for yours):
	userAgent := "terraform-powermax-provider/1.0.0"
	jar, err := cookiejar.New(nil)
	if err != nil {
		tflog.Error(ctx, "Got error while creating cookie jar")
	}

	httpclient := &http.Client{
		Timeout: (2000 * time.Second),
		Jar:     jar,
	}
	if insecure {
		httpclient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	} else {
		// Loading system certs by default if insecure is set to false
		// TODO: Check if we need to remove references to UseCerts from the code
		pool, err := x509.SystemCertPool()
		if err != nil {
			errSysCerts := errors.New("unable to initialize cert pool from system")
			return nil, errSysCerts
		}
		httpclient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            pool,
				InsecureSkipVerify: false,
			},
		}
	}

	url := fmt.Sprintf("%s/univmax/restapi", endpoint)
	basicAuthString := basicAuth(username, password)

	cfg := &pmaxop.Configuration{
		HTTPClient:    httpclient,
		DefaultHeader: make(map[string]string),
		UserAgent:     userAgent,
		Debug:         false,
		Servers: pmaxop.ServerConfigurations{
			{
				URL:         url,
				Description: url,
			},
		},
		OperationServers: map[string]pmaxop.ServerConfigurations{},
	}
	cfg.DefaultHeader = getHeaders()
	cfg.AddDefaultHeader("Authorization", "Basic "+basicAuthString)
	if serialNumber != "" {
		cfg.AddDefaultHeader("symid", serialNumber)
	}
	fmt.Printf("config %+v header %+v", cfg, cfg.DefaultHeader)

	apiClient := pmaxop.NewAPIClient(cfg)
	return apiClient, nil

}

// Generate the base 64 Authorization string from username / password.
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func getHeaders() map[string]string {
	header := make(map[string]string)

	header["Content-Type"] = "application/json; charset=utf-8"
	header["Accept"] = "application/json; charset=utf-8"
	return header

}

// NewClient returns the gopowermax client.
func NewGopmClient(endpoint, username, password, serialNumber, pmaxVersion string, insecure bool) (*pmax.Client, error) {
	cc := pmax.ConfigConnect{
		Endpoint: endpoint,
		Version:  pmaxVersion,
		Username: username,
		Password: password,
	}
	pmaxClient, err := pmax.NewClientWithArgs(endpoint, "Terraform Provider for PowerMax", insecure, false)
	if err != nil {
		return nil, err
	}
	err = pmax.Pmax.Authenticate(pmaxClient, context.Background(), &cc)
	if err != nil {
		return nil, err
	}
	pmaxClientWithID, ok := pmaxClient.WithSymmetrixID(serialNumber).(*pmax.Client)
	if !ok { // type assertion failed
		err := errors.New("creating client returned error")
		return nil, err
	}
	return pmaxClientWithID, nil
}
