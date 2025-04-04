package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetPodmanNetworkInfo(t *testing.T) {
	// Create a temporary directory for test
	tmpDir, err := os.MkdirTemp("", "podman-network-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test case 1: Valid network status file
	networkStatus := `{
		"podman": {
			"interfaces": {
				"eth0": {
					"subnets": [
						{
							"ipnet": "10.88.0.9/16",
							"gateway": "10.88.0.1"
						}
					],
					"mac_address": "f2:99:8d:fb:5a:57"
				}
			}
		}
	}`

	err = os.WriteFile(filepath.Join(tmpDir, "network.status"), []byte(networkStatus), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	ip, mac, err := getPodmanNetworkInfo(tmpDir)
	if err != nil {
		t.Errorf("getPodmanNetworkInfo failed: %v", err)
	}

	expectedIP := "10.88.0.9/16"
	expectedMAC := "f2:99:8d:fb:5a:57"

	if ip != expectedIP {
		t.Errorf("Expected IP %s, got %s", expectedIP, ip)
	}
	if mac != expectedMAC {
		t.Errorf("Expected MAC %s, got %s", expectedMAC, mac)
	}

	// Test case 2: Missing network status file
	emptyDir, err := os.MkdirTemp("", "empty-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(emptyDir)

	ip, mac, err = getPodmanNetworkInfo(emptyDir)
	if err != nil {
		t.Errorf("getPodmanNetworkInfo with missing file should not return error, got: %v", err)
	}
	if ip != "" || mac != "" {
		t.Errorf("Expected empty IP and MAC for missing file, got IP=%s, MAC=%s", ip, mac)
	}

	// Test case 3: Invalid JSON
	err = os.WriteFile(filepath.Join(tmpDir, "network.status"), []byte("invalid json"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	ip, mac, err = getPodmanNetworkInfo(tmpDir)
	if err == nil {
		t.Error("getPodmanNetworkInfo should fail with invalid JSON")
	}
}