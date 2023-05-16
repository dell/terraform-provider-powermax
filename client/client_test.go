package client

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	// Test case 1: Valid arguments
	endpoint := "https://10.225.104.36:8443/"
	username := "smc"
	password := "smc"
	serialNumber := "000120000605"
	pmaxVersion := "100"
	insecure := true

	client, err := NewClient(endpoint, username, password, serialNumber, pmaxVersion, insecure)

	if err != nil {
		t.Errorf("Error creating new client: %v", err)
	}

	if client.PmaxClient == nil {
		t.Errorf("Error creating new client: PmaxClient is nil")
	}

	if client.SymmetrixID != serialNumber {
		t.Errorf("Error creating new client: SymmetrixID does not match")
	}
}

func TestNewClientFail(t *testing.T) {
	// Test case 2: Invalid arguments
	endpoint := "https://10.225.104.36:8443/"
	username := "bad"
	password := "reallybad"
	serialNumber := "0123"
	pmaxVersion := "100"
	insecure := true

	_, err := NewClient(endpoint, username, password, serialNumber, pmaxVersion, insecure)

	if err != nil {
		t.Log("Should show error when bad username or password")
		return
	}

	t.Errorf("There should be an error here since we gave bad username and password: %v", err)

}
