package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"observability.dtci.technology/opentelemetry-collector/cmd/opentelemetry-installer/utils"
)

type Installer struct {
	otelVersion  string
	os           string
	arch         string
	otelEndpoint string
	serviceName  string
}

func newInstaller(otelVersion string, otelEndpoint string, serviceName string) *Installer {
	os, arch := utils.GetOsDetails()
	return &Installer{otelVersion, os, arch, otelEndpoint, serviceName}
}

func (i Installer) download() (string, error) {
	fileName := fmt.Sprintf("otelcol-contrib_%s_%s_%s.rpm", i.otelVersion, i.os, i.arch)
	url := fmt.Sprintf("https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v%s/%s", i.otelVersion, fileName)
	fmt.Printf("Downloading Opentelemetry Collector from %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return fileName, err
	}
	defer resp.Body.Close()

	out, err := os.Create(fileName)
	if err != nil {
		fmt.Println("x")
		return fileName, err
	}
	defer out.Close()
	fmt.Println("xx")
	_, err = io.Copy(out, resp.Body)
	fmt.Println("Download finished")
	return fileName, err
}

func (i Installer) Install() error {
	fileName, err := i.download()
	if err != nil {
		log.Fatal(err)
	}

	if i.os == "windows" {
		return fmt.Errorf("your OS (%s/%s) is not yet supported", i.os, i.arch)
	}
	fmt.Println("Installing Opentelemetry Collector...")
	fmt.Printf("sudo rpm -ivh %s\n", fileName)
	out, _ := exec.Command("sudo", "rpm", "-ivh", fileName).Output()
	fmt.Printf("%s", string(out))

	utils.DeleteFileIfExists("/etc/otelcol/otelcol.conf")
	i.Yamltool()

	out, err = exec.Command("sudo", "systemctl", "restart", "otelcol").Output()
	fmt.Printf("%s", string(out))

	return err
}

func (i Installer) Yamltool() {
	filePath := "./config/otel_config.yaml"

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	modifiedData := string(data)
	modifiedData = strings.ReplaceAll(modifiedData, "%SERVICE_NAME%", i.serviceName)
	modifiedData = strings.ReplaceAll(modifiedData, "%OTEL_ENDPOINT%", i.otelEndpoint)

	err = os.WriteFile("/etc/otelcol/otelcol.conf", []byte(modifiedData), 0644)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Opentelemetry config file is created")
}
