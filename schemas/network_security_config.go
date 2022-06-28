package schemas

type NetworkSecurityConfig struct {
	BaseConfig     *BaseConfig     `xml:"base-config"`
	DomainConfigs  *[]DomainConfig `xml:"domain-config"`
	DebugOverrides *DebugOverrides `xml:"debug-overrides"`
}

type BaseConfig struct {
	CleartextTrafficPermitted *bool         `xml:"cleartextTrafficPermitted,attr"`
	TrustAnchors              *TrustAnchors `xml:"trust-anchors"`
}

type DomainConfig struct {
	CleartextTrafficPermitted *bool           `xml:"cleartextTrafficPermitted,attr"`
	Domains                   []Domain        `xml:"domain"`
	TrustAnchors              *TrustAnchors   `xml:"trust-anchors"`
	PinSet                    *PinSet         `xml:"pin-set"`
	DomainConfigs             *[]DomainConfig `xml:"domain-config"`
}

type Domain struct {
	IncludeSubdomains *bool  `xml:"includeSubdomains,attr"`
	Domain            string `xml:",chardata"`
}

type DebugOverrides struct {
	TrustAnchors *TrustAnchors `xml:"trust-anchors"`
}

type TrustAnchors struct {
	Certificates *[]Certificates `xml:"certificates"`
}

type Certificates struct {
	Source       string `xml:"src,attr"`
	OverridePins *bool  `xml:"overridePins,attr"`
}

type PinSet struct {
	Expiration *string `xml:"expiration,attr"`
	Pins       *[]Pin  `xml:"pin"`
}

type Pin struct {
	Digest     string `xml:"digest,attr"`
	CertDigest string `xml:",chardata"`
}
