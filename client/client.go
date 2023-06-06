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
	"errors"

	pmax "github.com/dell/gopowermax/v2"
)

// Client type is to hold powermax client and symmetrix ID.
type Client struct {
	PmaxClient  *pmax.Client
	SymmetrixID string
}

// NewClient returns the gopowermax client.
func NewClient(endpoint, username, password, serialNumber, pmaxVersion string, insecure bool) (*Client, error) {
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
	client := Client{
		SymmetrixID: serialNumber,
		PmaxClient:  pmaxClientWithID,
	}
	return &client, nil
}
