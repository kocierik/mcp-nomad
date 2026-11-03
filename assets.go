package main

import _ "embed"

var (
	//go:embed README.md
	embeddedReadme string
	//go:embed LICENSE
	embeddedLicense string
)

func readmeAsset() string {
	return embeddedReadme
}

func licenseAsset() string {
	return embeddedLicense
}
