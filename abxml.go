package apkstat

import (
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
)

type ResXMLTreeHeader struct {
	Header ResChunkHeader
}

type ResXMLTreeNode struct {
	Header     ResChunkHeader
	LineNumber uint32
	Comment    ResStringPoolRef
}

type ResXMLTreeNamespaceExt struct {
	Prefix ResStringPoolRef
	URI    ResStringPoolRef
}

type ResXMLTreeEndElementExt struct {
	NS   ResStringPoolRef
	Name ResStringPoolRef
}

type ResXMLTreeAttrExt struct {
	NS             ResStringPoolRef
	Name           ResStringPoolRef
	AttributeStart uint16
	AttributeSize  uint16
	AttributeCount uint16
	IDIndex        uint16
	ClassIndex     uint16
	StyleIndex     uint16
}

type ResXMLTreeAttribute struct {
	NS         ResStringPoolRef
	Name       ResStringPoolRef
	RawValue   ResStringPoolRef
	TypedValue ResValue
}

func parseXMLElement(sr *io.SectionReader, sp map[ResStringPoolRef]string) (*xml.StartElement, error) {
	var e xml.StartElement

	node := new(ResXMLTreeNode)
	if err := binary.Read(sr, binary.LittleEndian, node); err != nil {
		return nil, err
	}

	element := new(ResXMLTreeAttrExt)
	if err := binary.Read(sr, binary.LittleEndian, element); err != nil {
		return nil, err
	}
	e.Name = xml.Name{Space: sp[element.NS], Local: sp[element.Name]}

	for i := 0; i < int(element.AttributeCount); i++ {
		attr := new(ResXMLTreeAttribute)
		if err := binary.Read(sr, binary.LittleEndian, attr); err != nil {
			return nil, err
		}

		var value string
		switch attr.TypedValue.DataType {
		case TypeNull:
			value = ""
		case TypeReference:
			value = fmt.Sprintf("@0x%08X", attr.TypedValue.Data)
		case TypeString:
			value = sp[attr.RawValue]
		case TypeFloat:
			value = fmt.Sprintf("%f", float32(attr.TypedValue.Data))
		case TypeIntDec:
			value = fmt.Sprintf("%d", attr.TypedValue.Data)
		case TypeIntHex:
			value = fmt.Sprintf("0x%08X", attr.TypedValue.Data)
		case TypeIntBoolean:
			if attr.TypedValue.Data == 1 {
				value = "true"
			} else {
				value = "false"
			}
		}

		e.Attr = append(e.Attr, xml.Attr{
			Name:  xml.Name{Space: sp[attr.NS], Local: sp[attr.Name]},
			Value: value,
		})
	}

	return &e, nil
}
