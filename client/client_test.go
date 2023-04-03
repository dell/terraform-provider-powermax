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
