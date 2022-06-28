package apk

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"io"
)

// APK is a representation of an Android APK file.
type APK struct {
	manifest Manifest
}

// Open opens an APK at path name and returns a new APK if successful. It automatically parses the
// app's Android manifest and resource table, resolving resource table references from the manifest
// as necessary.
func Open(name string) (*APK, error) {
	z, err := zip.OpenReader(name)
	if err != nil {
		return nil, err
	}
	defer z.Close()

	rawTable, err := z.Open("resources.arsc")
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

	rawManifest, err := z.Open("AndroidManifest.xml")
	if err != nil {
		return nil, err
	}
	defer rawManifest.Close()
	rawManifestBytes, err := io.ReadAll(rawManifest)
	if err != nil {
		return nil, err
	}
	xmlFile, err := NewXMLFile(bytes.NewReader(rawManifestBytes), table, nil)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	if err := xml.Unmarshal([]byte(xmlFile.String()), &manifest); err != nil {
		return nil, err
	}

	apk := new(APK)
	apk.manifest = manifest

	return apk, nil
}

// Manifest returns an APKs Manifest.
func (a *APK) Manifest() Manifest {
	return a.manifest
}
