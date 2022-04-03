package apkstat

import (
	"encoding/binary"
	"io"
	"unicode/utf16"
)

type ResChunkHeader struct {
	Type       uint16
	HeaderSize uint16
	Size       uint32
}

const (
	ResNullType       = 0x0
	ResStringPoolType = 0x1
	ResTableType      = 0x2
	ResXMLType        = 0x3

	ResXMLFirstChunkType     = 0x100
	ResXMLStartNamespaceType = 0x100
	ResXMLEndNamespaceType   = 0x101
	ResXMLStartElementType   = 0x102
	ResXMLEndElementType     = 0x103
	ResXMLCDataType          = 0x104
	ResXMLLastChunkType      = 0x17f
	ResXMLResourceMapType    = 0x180

	ResTablePackageType           = 0x200
	ResTableTypeType              = 0x201
	ResTableTypeSpecType          = 0x202
	ResTableLibraryType           = 0x203
	ResTableOverlayableType       = 0x204
	ResTableOverlayablePolicyType = 0x205
	ResTableStagedAliasType       = 0x206
)

type DataType = uint8

const (
	TypeNull             DataType = 0x00
	TypeReference        DataType = 0x01
	TypeAttribute        DataType = 0x02
	TypeString           DataType = 0x03
	TypeFloat            DataType = 0x04
	TypeDimension        DataType = 0x05
	TypeFraction         DataType = 0x06
	TypeDynamicReference DataType = 0x07
	TypeDynamicAttribute DataType = 0x08

	TypeFirstInt   DataType = 0x10
	TypeIntDec     DataType = 0x10
	TypeIntHex     DataType = 0x11
	TypeIntBoolean DataType = 0x12

	TypeFirstColorInt DataType = 0x1c
	TypeIntColorARGB8 DataType = 0x1c
	TypeIntColorRGB8  DataType = 0x1d
	TypeIntColorARGB4 DataType = 0x1e
	TypeIntColorRGB4  DataType = 0x1f
	TypeLastColorInt  DataType = 0x1f
	TypLastInt        DataType = 0x1f
)

type ResValue struct {
	Size     uint16
	Res0     uint8
	DataType DataType
	Data     uint32
}

type ResStringPoolRef struct {
	Index uint32
}

type Flags uint32

const (
	SortedFlag Flags = 1 << 0
	UTF8Flag   Flags = 1 << 8
)

type ResStringPoolHeader struct {
	Header       ResChunkHeader
	StringCount  uint32
	StyleCount   uint32
	Flags        Flags
	StringsStart uint32
	StylesStart  uint32
}

func parseStringPool(sr *io.SectionReader) (map[ResStringPoolRef]string, error) {
	stringPool := make(map[ResStringPoolRef]string)

	sp := new(ResStringPoolHeader)
	if err := binary.Read(sr, binary.LittleEndian, sp); err != nil {
		return nil, err
	}

	sIndices := make([]uint32, sp.StringCount)
	if err := binary.Read(sr, binary.LittleEndian, sIndices); err != nil {
		return nil, err
	}

	if sp.Flags&UTF8Flag != UTF8Flag {
		for i, sStart := range sIndices {
			if _, err := sr.Seek(int64(sp.StringsStart+sStart), io.SeekStart); err != nil {
				return nil, err
			}
			var strlen uint16
			if err := binary.Read(sr, binary.LittleEndian, &strlen); err != nil {
				return nil, err
			}
			buf := make([]uint16, strlen)
			if err := binary.Read(sr, binary.LittleEndian, buf); err != nil {
				return nil, err
			}

			spRef := ResStringPoolRef{Index: uint32(i)}
			stringPool[spRef] = string(utf16.Decode(buf))
		}
	}

	return stringPool, nil
}

type ResStringPoolSpan struct {
	Name                ResStringPoolRef
	FirstChar, LastChar uint32
}
