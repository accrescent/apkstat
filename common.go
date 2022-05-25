package apkstat

import (
	"encoding/binary"
	"io"
	"unicode/utf16"
)

const maxReadBytes = 1 << 26 // 64 MiB

type resChunkHeader struct {
	Type       uint16
	HeaderSize uint16
	Size       uint32
}

const (
	resNullType       = 0x0
	resStringPoolType = 0x1
	resTableChunkType = 0x2
	resXMLType        = 0x3

	resXMLFirstChunkType     = 0x100
	resXMLStartNamespaceType = 0x100
	resXMLEndNamespaceType   = 0x101
	resXMLStartElementType   = 0x102
	resXMLEndElementType     = 0x103
	resXMLCDataType          = 0x104
	resXMLLastChunkType      = 0x17f
	resXMLResourceMapType    = 0x180

	resTablePackageType           = 0x200
	resTableTypeType              = 0x201
	resTableTypeSpecType          = 0x202
	resTableLibraryType           = 0x203
	resTableOverlayableType       = 0x204
	resTableOverlayablePolicyType = 0x205
	resTableStagedAliasType       = 0x206
)

type dataType = uint8

const (
	typeNull             dataType = 0x00
	typeReference        dataType = 0x01
	typeAttribute        dataType = 0x02
	typeString           dataType = 0x03
	typeFloat            dataType = 0x04
	typeDimension        dataType = 0x05
	typeFraction         dataType = 0x06
	typeDynamicReference dataType = 0x07
	typeDynamicAttribute dataType = 0x08

	typeFirstInt   dataType = 0x10
	typeIntDec     dataType = 0x10
	typeIntHex     dataType = 0x11
	typeIntBoolean dataType = 0x12

	typeFirstColorInt dataType = 0x1c
	typeIntColorARGB8 dataType = 0x1c
	typeIntColorRGB8  dataType = 0x1d
	typeIntColorARGB4 dataType = 0x1e
	typeIntColorRGB4  dataType = 0x1f
	typeLastColorInt  dataType = 0x1f
	typeLastInt       dataType = 0x1f
)

type resValue struct {
	Size     uint16
	Res0     uint8
	DataType dataType
	Data     uint32
}

type resStringPoolRef struct {
	Index uint32
}

type flags uint32

const (
	sortedFlag flags = 1 << 0
	utf8Flag   flags = 1 << 8
)

type resStringPoolHeader struct {
	Header       resChunkHeader
	StringCount  uint32
	StyleCount   uint32
	Flags        flags
	StringsStart uint32
	StylesStart  uint32
}

func parseStringPool(sr *io.SectionReader) (map[resStringPoolRef]string, error) {
	stringPool := make(map[resStringPoolRef]string)

	sp := new(resStringPoolHeader)
	if err := binary.Read(sr, binary.LittleEndian, sp); err != nil {
		return nil, err
	}

	sIndices := make([]uint32, sp.StringCount)
	if err := binary.Read(sr, binary.LittleEndian, sIndices); err != nil {
		return nil, err
	}

	if sp.Flags&utf8Flag == utf8Flag { // UTF-8
		for i, sStart := range sIndices {
			if _, err := sr.Seek(int64(sp.StringsStart+sStart), io.SeekStart); err != nil {
				return nil, err
			}

			if _, err := parseVar8Len(sr); err != nil {
				return nil, err
			}
			size, err := parseVar8Len(sr)
			if err != nil {
				return nil, err
			}

			buf := make([]uint8, size)
			if err := binary.Read(sr, binary.LittleEndian, buf); err != nil {
				return nil, err
			}

			spRef := resStringPoolRef{Index: uint32(i)}
			stringPool[spRef] = string(buf)
		}
	} else { // UTF-16
		for i, sStart := range sIndices {
			if _, err := sr.Seek(int64(sp.StringsStart+sStart), io.SeekStart); err != nil {
				return nil, err
			}

			size, err := parseVar16Len(sr)
			if err != nil {
				return nil, err
			}

			buf := make([]uint16, size)
			if err := binary.Read(sr, binary.LittleEndian, buf); err != nil {
				return nil, err
			}

			spRef := resStringPoolRef{Index: uint32(i)}
			stringPool[spRef] = string(utf16.Decode(buf))
		}
	}

	return stringPool, nil
}

func parseVar8Len(sr *io.SectionReader) (int, error) {
	var size int
	var first, second uint8
	if err := binary.Read(sr, binary.LittleEndian, &first); err != nil {
		return 0, err
	}
	if (first & 0x80) != 0 { // high bit is set, read next byte
		if err := binary.Read(sr, binary.LittleEndian, &second); err != nil {
			return 0, err
		}
		size = (int(first&0x7F) << 8) | int(second)
	} else {
		size = int(first)
	}

	return size, nil
}

func parseVar16Len(sr *io.SectionReader) (int, error) {
	var size int
	var first, second uint16
	if err := binary.Read(sr, binary.LittleEndian, &first); err != nil {
		return 0, err
	}
	if first&0x8000 != 0 { // high bit is set, read next byte
		if err := binary.Read(sr, binary.LittleEndian, &second); err != nil {
			return 0, err
		}
		size = ((int(first) & 0x7FFF) << 16) | int(second)
	} else {
		size = int(first)
	}

	return size, nil
}

type resStringPoolSpan struct {
	name                resStringPoolRef
	FirstChar, LastChar uint32
}
