package apkstat

import (
	"encoding/binary"
	"encoding/xml"
	"errors"
	"io"
	"strconv"
)

type Manifest struct {
	Package         string           `xml:"package,attr"`
	VersionCode     int32            `xml:"http://schemas.android.com/apk/res/android versionCode,attr"`
	VersionName     string           `xml:"http://schemas.android.com/apk/res/android versionName,attr"`
	UsesPermissions []UsesPermission `xml:"uses-permission"`
}

type UsesPermission struct {
	Name string `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

const androidNS = "http://schemas.android.com/apk/res/android"
const maxReadBytes = 1 << 58 // 64 MiB

func ParseManifest(r io.ReaderAt) (*Manifest, error) {
	var m Manifest

	header := new(ResXMLTreeHeader)
	sr := io.NewSectionReader(r, 0, maxReadBytes)
	if err := binary.Read(sr, binary.LittleEndian, header); err != nil {
		return nil, err
	}
	if header.Header.Type != ResXMLType {
		return nil, errors.New("malformed header")
	}

	var stringPool map[ResStringPoolRef]string

	offset := int64(header.Header.HeaderSize)
	for offset < int64(header.Header.Size) {
		chunk := new(ResChunkHeader)
		if _, err := sr.Seek(offset, io.SeekStart); err != nil {
			return nil, err
		}
		if err := binary.Read(sr, binary.LittleEndian, chunk); err != nil {
			return nil, err
		}

		if _, err := sr.Seek(offset, io.SeekStart); err != nil {
			return nil, err
		}

		var err error
		switch chunk.Type {
		case ResStringPoolType:
			stringPool, err = parseStringPool(io.NewSectionReader(sr, offset, maxReadBytes))
		case ResXMLResourceMapType:
		case ResXMLStartNamespaceType:
		case ResXMLEndNamespaceType:
		case ResXMLStartElementType:
			e, err := parseXMLElement(io.NewSectionReader(sr, offset, maxReadBytes), stringPool)
			if err != nil {
				return nil, err
			}
			err = m.populateElement(e)
		case ResXMLEndElementType:
		default:
			return nil, errors.New("encountered invalid chunk type")
		}
		if err != nil {
			return nil, err
		}

		offset += int64(chunk.Size)
	}

	return &m, nil
}

func (m *Manifest) populateElement(e *xml.StartElement) error {
	switch e.Name.Local {
	case "manifest":
		for i := 0; i < len(e.Attr); i++ {
			ns := e.Attr[i].Name.Space
			name := e.Attr[i].Name.Local
			switch {
			case ns == "" && name == "package":
				m.Package = e.Attr[i].Value
			case ns == androidNS && name == "versionCode":
				v, err := strconv.ParseInt(e.Attr[i].Value, 10, 32)
				if err != nil {
					return err
				}
				m.VersionCode = int32(v)
			case ns == androidNS && name == "versionName":
				m.VersionName = e.Attr[i].Value
			}
		}
	case "uses-permission":
		for i := 0; i < len(e.Attr); i++ {
			ns := e.Attr[i].Name.Space
			name := e.Attr[i].Name.Local
			switch {
			case ns == androidNS && name == "name":
				permission := UsesPermission{Name: e.Attr[i].Value}
				m.UsesPermissions = append(m.UsesPermissions, permission)
			}
		}
	}

	return nil
}
