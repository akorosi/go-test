package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

func GetOsDetails() (string, string) {
	os := runtime.GOOS
	arch := runtime.GOARCH

	var distribution, architecture string

	switch os {
	case "linux":
		// Check for common Linux distributions
		if _, err := exec.LookPath("dpkg"); err == nil {
			distribution = "Debian-based"
		} else if _, err := exec.LookPath("rpm"); err == nil {
			distribution = "RPM-based"
		} else {
			distribution = "Unknown Linux distribution"
		}
	case "darwin":
		distribution = "macOS"
	case "windows":
		distribution = "Windows"
	default:
		distribution = "Unknown OS"
	}

	switch arch {
	case "i386":
		architecture = "386"
	default:
		architecture = arch
	}

	fmt.Printf("OS: %s\nArchitecture: %s\n", distribution, architecture)
	return os, architecture
}

func GetPublicIP() (string, error) {
	url := "https://api.ipify.org?format=json"

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to get IP address: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var data map[string]string
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Extract the IP address from the "ip" key
	ip, ok := data["ip"]
	if !ok {
		return "", fmt.Errorf("missing 'ip' key in response")
	}

	return ip, nil

}

func DeleteFileIfExists(filename string) error {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	err = os.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}
