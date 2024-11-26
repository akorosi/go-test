package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"observability.dtci.technology/opentelemetry-collector/cmd/opentelemetry-installer/utils"
)

func getOptentelemetryEndpoint() (string, error) {
	ipAddr, error := utils.GetPublicIP()
	if error != nil {
		log.Fatal(error)
	}
	fmt.Println("Your Public IP: ", ipAddr)

	jsonData, _ := os.ReadFile("./config/gateways.json")

	var gateways Gateways
	error = json.Unmarshal(jsonData, &gateways)
	if error != nil {
		log.Fatal(error)
	}

	for _, gateway := range gateways.Gateways {
		for _, cidrString := range gateway.CidrRange {
			_, ipnet, err := net.ParseCIDR(cidrString)
			if err != nil {
				continue
			}

			if ipnet.Contains(net.ParseIP(ipAddr)) {
				fmt.Println("Matching Configuration: ")
				fmt.Printf("Datacenter: %s, Opentelemetry Gateway: %s\n", gateway.DcName, gateway.URL)
				return gateway.URL, nil
			}
		}
	}

	return "", fmt.Errorf("no matching Gateway found")
}

func main() {
	otelVersion := *flag.String("o", "0.114.0", "[Optional] Opentelemetry Collector Version")
	serviceName := flag.String("s", "", "[Required] Service name")
	flag.Parse()
	if *serviceName == "" {
		flag.Usage()
		os.Exit(1)
	}
	otelEndpoint, _ := getOptentelemetryEndpoint()
	installer := newInstaller(otelVersion, otelEndpoint, *serviceName)
	err := installer.Install()

	if err != nil {
		log.Fatal(err)
	}
}
