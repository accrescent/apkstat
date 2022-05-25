package apkstat

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"io"
)

type APK struct {
	manifest Manifest
}

func OpenAPK(name string) (*APK, error) {
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

func (a *APK) Manifest() Manifest {
	return a.manifest
}
