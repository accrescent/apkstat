// Package apk implements a parser for Android APKs.
//
// The APK type represents an APK file and is the API most users should use. You can open an APK
// file with apk.Open, or if you want more control over resource resolution, apk.OpenWithConfig.
//
// Most information about an APK is contained in its manifest. The Manifest() method will return an
// APK's Manifest.
//
// BUG(lberrymage): Some resource table references in binary XML are incorrectly parsed as empty
// strings.
package apk

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"io"
	"os"

	"github.com/accrescent/apkstat/schemas"
)

// APK is a representation of an Android APK file.
type APK struct {
	zipReader *zip.Reader
	config    *ResTableConfig
	table     *ResTable
	manifest  *Manifest
}

// Open opens an APK at path name and returns a new APK if successful. It automatically parses the
// app's Android manifest and resource table, resolving resource table references from the manifest
// as necessary.
func Open(name string) (*APK, error) {
	return OpenWithConfig(name, nil)
}

// OpenWithConfig opns an APK an path name and returns a new APK if successful. It automatically
// parses the app's Android manifest and resource table, using config to resolve resource tables
// from the manifest as necessary.
func OpenWithConfig(name string, config *ResTableConfig) (*APK, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return FromReaderWithConfig(f, info.Size(), config)
}

// FromReader opens an APK of the given size from reader r and returns a new APK if successful. It
// automatically parses the app's Android manifest and resource table, resolving table references
// from the manifest as necessary.
func FromReader(r io.ReaderAt, size int64) (*APK, error) {
	return FromReaderWithConfig(r, size, nil)
}

// FromReaderWithConfig opens an APK of the given size from reader r and returns a new APK if
// successful. It automatically parses the app's Android manifest and resource table, using config
// to resolve resource table references from the manifest as necessary.
func FromReaderWithConfig(r io.ReaderAt, size int64, config *ResTableConfig) (*APK, error) {
	apk := new(APK)
	apk.config = config

	z, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	apk.zipReader = z

	rawTable, err := apk.zipReader.Open("resources.arsc")
	if err != nil {
		return nil, err
	}
	defer rawTable.Close()
	rawTableBytes, err := io.ReadAll(rawTable)
	if err != nil {
		return nil, err
	}
	table, err := NewResTable(bytes.NewReader(rawTableBytes))
	if err != nil {
		return nil, err
	}
	apk.table = table

	xmlFile, err := apk.OpenXML("AndroidManifest.xml")
	if err != nil {
		return nil, err
	}

	if err := xml.Unmarshal([]byte(xmlFile.String()), &apk.manifest); err != nil {
		return nil, err
	}

	return apk, nil
}

// OpenXML is a utility function for opening an arbitrary Android binary XML file within an APK. If
// a resource table config was specified when opening the APK with apk.OpenWithConfig, it will be
// used.
func (a *APK) OpenXML(name string) (*XMLFile, error) {
	return a.OpenXMLWithConfig(name, nil)
}

// OpenXMLWithConfig is like OpenXML, but allows for specifying a ResTableConfig after opening the
// APK. It overrides the APK's ResTableConfig for this function call but doesn't modify the APK's
// ResTableConfig that was specified at open time.
func (a *APK) OpenXMLWithConfig(name string, config *ResTableConfig) (*XMLFile, error) {
	rawXML, err := a.zipReader.Open(name)
	if err != nil {
		return nil, err
	}
	defer rawXML.Close()
	rawXMLBytes, err := io.ReadAll(rawXML)
	if err != nil {
		return nil, err
	}
	xmlFile, err := NewXMLFile(bytes.NewReader(rawXMLBytes), a.table, config)
	if err != nil {
		return nil, err
	}

	return xmlFile, nil
}

// SetConfig sets the APK's ResTableConfig for future operations.
func (a *APK) SetConfig(config *ResTableConfig) {
	a.config = config
}

// Manifest returns an APK's Manifest.
func (a *APK) Manifest() Manifest {
	return *a.manifest
}

// DataExtractionRules is a utility function which parses an APK's data extraction rules into a
// struct.
func (a *APK) DataExtractionRules() (*schemas.DataExtractionRules, error) {
	manifestRules := a.manifest.Application.DataExtractionRules
	if manifestRules == nil {
		return nil, XMLResourceNotFound
	}

	xmlFile, err := a.OpenXML(*manifestRules)
	if err != nil {
		return nil, err
	}

	var rules schemas.DataExtractionRules
	if err := xml.Unmarshal([]byte(xmlFile.String()), &rules); err != nil {
		return nil, err
	}

	return &rules, nil
}

// NetworkSecurityConfig is a utility function which parses an APK's network security config into a
// struct.
func (a *APK) NetworkSecurityConfig() (*schemas.NetworkSecurityConfig, error) {
	manifestNSConfig := a.manifest.Application.NetworkSecurityConfig
	if manifestNSConfig == nil {
		return nil, XMLResourceNotFound
	}

	xmlFile, err := a.OpenXML(*manifestNSConfig)
	if err != nil {
		return nil, err
	}

	var nsConfig schemas.NetworkSecurityConfig
	if err := xml.Unmarshal([]byte(xmlFile.String()), &nsConfig); err != nil {
		return nil, err
	}

	return &nsConfig, nil
}
