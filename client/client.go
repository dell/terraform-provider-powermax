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

// NewClient returns the gopowermax client
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
	pmaxClientWithID := pmaxClient.WithSymmetrixID(serialNumber).(*pmax.Client)

	client := Client{
		SymmetrixID: serialNumber,
		PmaxClient:  pmaxClientWithID,
	}
	return &client, nil
}
