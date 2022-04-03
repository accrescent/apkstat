package apkstat

import (
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

const maxReadBytes = 1 << 58 // 64 MiB

type Parser struct {
	stringPool map[resStringPoolRef]string
	nsToPrefix map[resStringPoolRef]resStringPoolRef
	namespaces map[resStringPoolRef]resStringPoolRef
	xml        strings.Builder
}

func NewParser(r io.ReaderAt) (*Parser, error) {
	p := new(Parser)

	fmt.Fprintf(&p.xml, xml.Header)

	header := new(resXMLTreeHeader)
	sr := io.NewSectionReader(r, 0, maxReadBytes)
	if err := binary.Read(sr, binary.LittleEndian, header); err != nil {
		return nil, err
	}
	if header.Header.Type != resXMLType {
		return nil, errors.New("malformed header")
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
			p.stringPool, err = parseStringPool(io.NewSectionReader(sr, offset, maxReadBytes))
		case resXMLResourceMapType:
		case resXMLStartNamespaceType:
			err = p.parseStartNamespace(sr)
		case resXMLEndNamespaceType:
			err = p.parseEndNamespace(sr)
		case resXMLStartElementType:
			err = p.parseXMLStartElement(sr)
		case resXMLEndElementType:
			err = p.parseXMLEndElement(sr)
		default:
			return nil, errors.New("encountered invalid chunk type")
		}
		if err != nil {
			return nil, err
		}

		offset += int64(chunk.Size)
	}

	return p, nil
}

func (p *Parser) String() string {
	return p.xml.String()
}

func (p *Parser) parseStartNamespace(sr *io.SectionReader) error {
	node := new(resXMLTreeNode)
	if err := binary.Read(sr, binary.LittleEndian, node); err != nil {
		return err
	}

	ns := new(resXMLTreeNamespaceExt)
	if err := binary.Read(sr, binary.LittleEndian, ns); err != nil {
		return err
	}

	if p.nsToPrefix == nil {
		p.nsToPrefix = make(map[resStringPoolRef]resStringPoolRef)
	}
	if p.namespaces == nil {
		p.namespaces = make(map[resStringPoolRef]resStringPoolRef)
	}
	p.nsToPrefix[ns.URI] = ns.Prefix
	p.namespaces[ns.URI] = ns.Prefix
	return nil
}

func (p *Parser) parseEndNamespace(sr *io.SectionReader) error {
	node := new(resXMLTreeNode)
	if err := binary.Read(sr, binary.LittleEndian, node); err != nil {
		return err
	}

	ns := new(resXMLTreeNamespaceExt)
	if err := binary.Read(sr, binary.LittleEndian, ns); err != nil {
		return err
	}

	delete(p.namespaces, ns.URI)

	return nil
}

func (p *Parser) parseXMLStartElement(sr *io.SectionReader) error {
	node := new(resXMLTreeNode)
	if err := binary.Read(sr, binary.LittleEndian, node); err != nil {
		return err
	}

	element := new(resXMLTreeAttrExt)
	if err := binary.Read(sr, binary.LittleEndian, element); err != nil {
		return err
	}

	fmt.Fprintf(&p.xml, "<%s", p.nsPrefix(element.NS, element.Name))

	for uri, prefix := range p.nsToPrefix {
		fmt.Fprintf(&p.xml, " xmlns:%s=\"", p.stringPool[prefix])
		if err := xml.EscapeText(&p.xml, []byte(p.stringPool[uri])); err != nil {
			return err
		}
		fmt.Fprintf(&p.xml, "\"")
	}
	p.nsToPrefix = nil

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
			value = fmt.Sprintf("@0x%08X", attr.TypedValue.Data)
		case typeString:
			value = p.stringPool[attr.RawValue]
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

		fmt.Fprintf(&p.xml, " %s=\"", p.nsPrefix(attr.NS, attr.Name))
		if err := xml.EscapeText(&p.xml, []byte(value)); err != nil {
			return err
		}
		fmt.Fprintf(&p.xml, "\"")
	}
	fmt.Fprintf(&p.xml, ">")

	return nil
}

func (p *Parser) parseXMLEndElement(sr *io.SectionReader) error {
	node := new(resXMLTreeNode)
	if err := binary.Read(sr, binary.LittleEndian, node); err != nil {
		return err
	}

	element := new(resXMLTreeEndElementExt)
	if err := binary.Read(sr, binary.LittleEndian, element); err != nil {
		return err
	}

	fmt.Fprintf(&p.xml, "</%s>", p.nsPrefix(element.NS, element.Name))

	return nil
}

func (p *Parser) nsPrefix(ns resStringPoolRef, name resStringPoolRef) string {
	if ns.Index == 0xFFFFFFFF {
		return fmt.Sprintf("%s", p.stringPool[name])
	} else {
		return fmt.Sprintf("%s:%s", p.stringPool[p.namespaces[ns]], p.stringPool[name])
	}
}
