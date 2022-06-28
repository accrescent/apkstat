package apk

import (
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type XMLFile struct {
	stringPool map[resStringPoolRef]string
	nsToPrefix map[resStringPoolRef]resStringPoolRef
	namespaces map[resStringPoolRef]resStringPoolRef
	xml        strings.Builder
	table      *ResTable
	cfg        *ResTableConfig
}

func NewXMLFile(r io.ReaderAt, t *ResTable, cfg *ResTableConfig) (*XMLFile, error) {
	f := new(XMLFile)
	f.table = t
	f.cfg = cfg

	fmt.Fprintf(&f.xml, xml.Header)

	header := new(resXMLTreeHeader)
	sr := io.NewSectionReader(r, 0, maxReadBytes)
	if err := binary.Read(sr, binary.LittleEndian, header); err != nil {
		return nil, err
	}
	if header.Header.Type != resXMLType {
		return nil, MalformedHeader
	}

	offset := int64(header.Header.HeaderSize)
	for offset < int64(header.Header.Size) {
		chunk := new(resChunkHeader)
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
		case resStringPoolType:
			f.stringPool, err = parseStringPool(io.NewSectionReader(
				sr,
				offset,
				maxReadBytes-offset,
			))
		case resXMLResourceMapType:
		case resXMLStartNamespaceType:
			err = f.parseStartNamespace(sr)
		case resXMLEndNamespaceType:
			err = f.parseEndNamespace(sr)
		case resXMLStartElementType:
			err = f.parseXMLStartElement(sr)
		case resXMLEndElementType:
			err = f.parseXMLEndElement(sr)
		default:
			return nil, InvalidChunkType
		}
		if err != nil {
			return nil, err
		}

		offset += int64(chunk.Size)
	}

	return f, nil
}

func (f *XMLFile) String() string {
	return f.xml.String()
}

// parseStartNamespace parses a resXMLTreeNamespaceExt as a namespace start and updates the parsing
// state of f as necessary.
func (f *XMLFile) parseStartNamespace(sr *io.SectionReader) error {
	node := new(resXMLTreeNode)
	if err := binary.Read(sr, binary.LittleEndian, node); err != nil {
		return err
	}

	ns := new(resXMLTreeNamespaceExt)
	if err := binary.Read(sr, binary.LittleEndian, ns); err != nil {
		return err
	}

	if f.nsToPrefix == nil {
		f.nsToPrefix = make(map[resStringPoolRef]resStringPoolRef)
	}
	if f.namespaces == nil {
		f.namespaces = make(map[resStringPoolRef]resStringPoolRef)
	}
	f.nsToPrefix[ns.URI] = ns.Prefix
	f.namespaces[ns.URI] = ns.Prefix

	return nil
}

// parseEndNamespace parses a resXMLTreeNamespaceExt as a namespace end and updates the parsing
// state of f as necessary.
func (f *XMLFile) parseEndNamespace(sr *io.SectionReader) error {
	node := new(resXMLTreeNode)
	if err := binary.Read(sr, binary.LittleEndian, node); err != nil {
		return err
	}

	ns := new(resXMLTreeNamespaceExt)
	if err := binary.Read(sr, binary.LittleEndian, ns); err != nil {
		return err
	}

	delete(f.namespaces, ns.URI)

	return nil
}

// parseXMLStartElement parses an XML start element along with its attributes and updates the
// parsing state of f as necessary.
func (f *XMLFile) parseXMLStartElement(sr *io.SectionReader) error {
	node := new(resXMLTreeNode)
	if err := binary.Read(sr, binary.LittleEndian, node); err != nil {
		return err
	}

	element := new(resXMLTreeAttrExt)
	if err := binary.Read(sr, binary.LittleEndian, element); err != nil {
		return err
	}

	fmt.Fprintf(&f.xml, "<%s", f.nsPrefix(element.NS, element.Name))

	for uri, prefix := range f.nsToPrefix {
		fmt.Fprintf(&f.xml, " xmlns:%s=\"", f.stringPool[prefix])
		if err := xml.EscapeText(&f.xml, []byte(f.stringPool[uri])); err != nil {
			return err
		}
		fmt.Fprintf(&f.xml, "\"")
	}
	f.nsToPrefix = nil

	for i := 0; i < int(element.AttributeCount); i++ {
		attr := new(resXMLTreeAttribute)
		if err := binary.Read(sr, binary.LittleEndian, attr); err != nil {
			return err
		}

		var value string
		switch attr.TypedValue.DataType {
		case typeNull:
			value = ""
		case typeReference:
			if f.table != nil {
				r, err := f.table.getResource(resID(attr.TypedValue.Data), f.cfg)
				if err != nil {
					return err
				}
				value = r
			} else {
				value = fmt.Sprintf("@0x%08X", attr.TypedValue.Data)
			}
		case typeString:
			value = f.stringPool[attr.RawValue]
		case typeFloat:
			value = fmt.Sprintf("%f", float32(attr.TypedValue.Data))
		case typeIntDec:
			value = fmt.Sprintf("%d", attr.TypedValue.Data)
		case typeIntHex:
			value = fmt.Sprintf("0x%08X", attr.TypedValue.Data)
		case typeIntBoolean:
			if attr.TypedValue.Data != 0 {
				value = "true"
			} else {
				value = "false"
			}
		}

		fmt.Fprintf(&f.xml, " %s=\"", f.nsPrefix(attr.NS, attr.Name))
		if err := xml.EscapeText(&f.xml, []byte(value)); err != nil {
			return err
		}
		fmt.Fprintf(&f.xml, "\"")
	}
	fmt.Fprintf(&f.xml, ">")

	return nil
}

// parseXMLEndElement parses an XML end element and updates the parsing state of f as necessary.
func (f *XMLFile) parseXMLEndElement(sr *io.SectionReader) error {
	node := new(resXMLTreeNode)
	if err := binary.Read(sr, binary.LittleEndian, node); err != nil {
		return err
	}

	element := new(resXMLTreeEndElementExt)
	if err := binary.Read(sr, binary.LittleEndian, element); err != nil {
		return err
	}

	fmt.Fprintf(&f.xml, "</%s>", f.nsPrefix(element.NS, element.Name))

	return nil
}

// nsPrefix takes a namespace and an XML attribute name as string pool references and returns the
// XML attribute prefixed with the namespace if the namespace string pool reference is not empty.
func (f *XMLFile) nsPrefix(ns resStringPoolRef, name resStringPoolRef) string {
	if ns.Index == 0xFFFFFFFF {
		return fmt.Sprintf("%s", f.stringPool[name])
	} else {
		return fmt.Sprintf("%s:%s", f.stringPool[f.namespaces[ns]], f.stringPool[name])
	}
}

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
