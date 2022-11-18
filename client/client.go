package client

import (
	"context"

	pmax "github.com/dell/gopowermax/v2"
)

// Client type is to hold powermax client and symmetrix ID
type Client struct {
	PmaxClient  *pmax.Client
	SymmetrixID string
}

var (
	newClientWithArgs = pmax.NewClientWithArgs
	authenticate      = pmax.Pmax.Authenticate
	withSymmetrixID   = pmax.Pmax.WithSymmetrixID
)

// NewClient returns the gopowermax client
func NewClient(endpoint, username, password, serialNumber, pmaxVersion string, insecure bool) (*Client, error) {
	cc := pmax.ConfigConnect{
		Endpoint: endpoint,
		Version:  pmaxVersion,
		Username: username,
		Password: password,
	}
	pmaxClient, err := newClientWithArgs(endpoint, "Terraform Provider for PowerMax", insecure, false)
	if err != nil {
		return nil, err
	}
	err = authenticate(pmaxClient, context.Background(), &cc)
	if err != nil {
		return nil, err
	}
	pmaxClientWithSymID := withSymmetrixID(pmaxClient, serialNumber).(*pmax.Client)

	client := Client{
		SymmetrixID: serialNumber,
		PmaxClient:  pmaxClientWithSymID,
	}
	return &client, nil
}
