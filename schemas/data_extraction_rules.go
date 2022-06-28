package schemas

type DataExtractionRules struct {
	CloudBackup    *CloudBackup    `xml:"cloud-backup"`
	DeviceTransfer *DeviceTransfer `xml:"device-transfer"`
}

type CloudBackup struct {
	DisableIfNoEncryptionCapabilities *bool      `xml:"disableIfNoEncryptionCapabilities,attr"`
	Includes                          *[]Include `xml:"include"`
	Excludes                          *[]Exclude `xml:"exclude"`
}

type DeviceTransfer struct {
	Includes *[]Include `xml:"include"`
	Excludes *[]Exclude `xml:"exclude"`
}

type Include struct {
	Domain string `xml:"domain,attr"`
	Path   string `xml:"path,attr"`
}

type Exclude struct {
	Domain string `xml:"domain,attr"`
	Path   string `xml:"path,attr"`
}
