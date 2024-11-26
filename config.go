package main

type Gateways struct {
	Gateways []Gateway `json:"gateways"`
}

type Gateway struct {
	URL         string   `json:"url"`
	DcName      string   `json:"dc_name"`
	CidrRange   []string `json:"cidr_range"`
	Environment string   `json:"environment"`
}
