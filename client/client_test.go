package client

import (
	"context"
	"fmt"
	"testing"

	pmax "github.com/dell/gopowermax/v2"
	"github.com/stretchr/testify/assert"
)

var endpoint string = "powermax_endpoint"
var username string = "username"
var password string = "password"
var serialNumber string = "123"
var pmaxVersion string = "100"
var insecure bool = true
var pmaxClient = &pmax.Client{}
var oldNewClientWithArgs = newClientWithArgs
var oldAuthenticate = authenticate
var oldWithSymmetrixID = withSymmetrixID

func init() {
	defer func() {
		newClientWithArgs = oldNewClientWithArgs
		authenticate = oldAuthenticate
		withSymmetrixID = oldWithSymmetrixID
	}()
}

func TestNewClient(t *testing.T) {
	newClientWithArgs = func(endpoint string,
		applicationName string,
		insecure,
		useCerts bool) (client pmax.Pmax, err error) {
		return pmaxClient, nil
	}
	authenticate = func(pmax pmax.Pmax, ctx context.Context, configConnect *pmax.ConfigConnect) error {
		return nil
	}
	withSymmetrixID = func(pmaxIn pmax.Pmax, symmetrixID string) pmax.Pmax {
		return pmaxClient
	}

	actualClient, actualErr := NewClient(endpoint, username, password, serialNumber, pmaxVersion, insecure)
	assert.Nil(t, actualErr)
	assert.NotNil(t, actualClient)
	assert.Equal(t, pmaxClient, actualClient.PmaxClient)
	assert.Equal(t, serialNumber, actualClient.SymmetrixID)
}

func TestNewClientWithCertError(t *testing.T) {
	newClientWithArgs = func(endpoint string,
		applicationName string,
		insecure,
		useCerts bool) (client pmax.Pmax, err error) {
		return nil, fmt.Errorf("invalid certificate")
	}
	actualClient, actualErr := NewClient(endpoint, username, password, serialNumber, pmaxVersion, insecure)
	assert.Nil(t, actualClient)
	assert.NotNil(t, actualErr)
	assert.Equal(t, "invalid certificate", actualErr.Error())
}

func TestNewClientWithAuthError(t *testing.T) {

	newClientWithArgs = func(endpoint string,
		applicationName string,
		insecure,
		useCerts bool) (client pmax.Pmax, err error) {
		return pmaxClient, nil
	}
	authenticate = func(pmax pmax.Pmax, ctx context.Context, configConnect *pmax.ConfigConnect) error {
		return fmt.Errorf("authentication failed, invalid username/password")
	}

	actualClient, actualErr := NewClient(endpoint, username, password, serialNumber, pmaxVersion, insecure)
	assert.Nil(t, actualClient)
	assert.NotNil(t, actualErr)
	assert.Equal(t, "authentication failed, invalid username/password", actualErr.Error())
}
