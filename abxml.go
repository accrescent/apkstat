package apkstat

import (
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
)

type resXMLTreeHeader struct {
	Header resChunkHeader
}

type resXMLTreeNode struct {
	Header     resChunkHeader
	LineNumber uint32
	Comment    resStringPoolRef
}

type resXMLTreeNamespaceExt struct {
	Prefix resStringPoolRef
	URI    resStringPoolRef
}

type resXMLTreeEndElementExt struct {
	NS   resStringPoolRef
	Name resStringPoolRef
}

type resXMLTreeAttrExt struct {
	NS             resStringPoolRef
	Name           resStringPoolRef
	AttributeStart uint16
	AttributeSize  uint16
	AttributeCount uint16
	IDIndex        uint16
	ClassIndex     uint16
	StyleIndex     uint16
}

type resXMLTreeAttribute struct {
	NS         resStringPoolRef
	Name       resStringPoolRef
	RawValue   resStringPoolRef
	TypedValue resValue
}

func parseXMLElement(sr *io.SectionReader, sp map[resStringPoolRef]string) (*xml.StartElement, error) {
	var e xml.StartElement

	node := new(resXMLTreeNode)
	if err := binary.Read(sr, binary.LittleEndian, node); err != nil {
		return nil, err
	}

	element := new(resXMLTreeAttrExt)
	if err := binary.Read(sr, binary.LittleEndian, element); err != nil {
		return nil, err
	}
	e.Name = xml.Name{Space: sp[element.NS], Local: sp[element.Name]}

	for i := 0; i < int(element.AttributeCount); i++ {
		attr := new(resXMLTreeAttribute)
		if err := binary.Read(sr, binary.LittleEndian, attr); err != nil {
			return nil, err
		}

		var value string
		switch attr.TypedValue.DataType {
		case typeNull:
			value = ""
		case typeReference:
			value = fmt.Sprintf("@0x%08X", attr.TypedValue.Data)
		case typeString:
			value = sp[attr.RawValue]
		case typeFloat:
			value = fmt.Sprintf("%f", float32(attr.TypedValue.Data))
		case typeIntDec:
			value = fmt.Sprintf("%d", attr.TypedValue.Data)
		case typeIntHex:
			value = fmt.Sprintf("0x%08X", attr.TypedValue.Data)
		case typeIntBoolean:
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
