package apkstat

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type ResParser struct {
	stringPool map[resStringPoolRef]string
	packages   map[uint32]*tablePackage
}

type tablePackage struct {
	typeStrings map[resStringPoolRef]string
	keyStrings  map[resStringPoolRef]string
	tableTypes  []*tableType
}

type tableType struct {
	entries []tableEntry
	header  resTableType
}

type tableEntry struct {
	key   resTableEntry
	value resValue
}

func NewResParser(r io.ReaderAt) (*ResParser, error) {
	p := new(ResParser)

	header := new(resTableHeader)
	sr := io.NewSectionReader(r, 0, maxReadBytes)
	if err := binary.Read(sr, binary.LittleEndian, header); err != nil {
		return nil, err
	}
	if header.Header.Type != resTableChunkType {
		return nil, errors.New("malformed header")
	}

	p.packages = make(map[uint32]*tablePackage)

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
		case resTablePackageType:
			err = p.parseTablePackage(io.NewSectionReader(sr, offset, maxReadBytes))
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

type resID uint32

func (id resID) pkg() uint32 {
	return uint32(id) >> 24
}

func (id resID) type_() int {
	return (int(id) >> 16) & 0xFF
}

func (id resID) entry() int {
	return int(id) & 0xFFFF
}

func (f *ResParser) getResource(id resID, config *ResTableConfig) (string, error) {
	pkg := id.pkg()
	entry := id.entry()

	p := f.packages[pkg]
	if p == nil {
		return "", fmt.Errorf("package 0x%02X not found", id.pkg())
	}

	var best *tableType

	for _, t := range p.tableTypes {
		switch {
		case int(t.header.ID) != id.type_():
			continue
		case !t.header.Config.match(config):
			continue
		case entry >= len(t.entries):
			continue
		case best == nil || t.header.Config.isBetterThan(&best.header.Config, config):
			best = t
		}
	}

	if best == nil || entry >= len(best.entries) {
		return "", fmt.Errorf("entry 0x%04X not found", entry)
	}

	e := best.entries[entry]
	v := e.value

	switch v.DataType {
	case typeNull:
		return "", nil
	case typeString:
		return f.stringPool[resStringPoolRef{v.Data}], nil
	case typeFloat:
		return fmt.Sprintf("%f", float32(v.Data)), nil
	case typeIntDec:
		return fmt.Sprintf("%d", v.Data), nil
	case typeIntHex:
		return fmt.Sprintf("0x%08X", v.Data), nil
	case typeIntBoolean:
		if v.Data == 1 {
			return "true", nil
		} else {
			return "false", nil
		}
	}

	return "", nil
}

func (p *ResParser) parseTablePackage(sr *io.SectionReader) error {
	pkg := new(tablePackage)

	header := new(resTablePackage)
	if err := binary.Read(sr, binary.LittleEndian, header); err != nil {
		return err
	}

	typeSR := io.NewSectionReader(sr, int64(header.TypeStrings), maxReadBytes)
	if typeStrings, err := parseStringPool(typeSR); err != nil {
		return err
	} else {
		pkg.typeStrings = typeStrings
	}

	keySR := io.NewSectionReader(sr, int64(header.KeyStrings), maxReadBytes)
	if keyStrings, err := parseStringPool(keySR); err != nil {
		return err
	} else {
		pkg.keyStrings = keyStrings
	}

	offset := int64(header.Header.HeaderSize)
	for offset < int64(header.Header.Size) {
		chunk := new(resChunkHeader)
		if _, err := sr.Seek(offset, io.SeekStart); err != nil {
			return err
		}
		if err := binary.Read(sr, binary.LittleEndian, chunk); err != nil {
			return err
		}

		if _, err := sr.Seek(offset, io.SeekStart); err != nil {
			return err
		}

		var err error
		switch chunk.Type {
		case resStringPoolType: // skip typestrings and keystrings
		case resTableTypeType:
			var tt *tableType
			tt, err = p.parseTableType(io.NewSectionReader(sr, offset, maxReadBytes))
			pkg.tableTypes = append(pkg.tableTypes, tt)
		case resTableTypeSpecType:
			// unimplemented
		default:
			return errors.New("encountered invalid chunk type")
		}
		if err != nil {
			return err
		}

		offset += int64(chunk.Size)
	}

	p.packages[header.ID] = pkg

	return nil
}

func (p *ResParser) parseTableType(sr *io.SectionReader) (*tableType, error) {
	header := new(resTableType)
	if err := binary.Read(sr, binary.LittleEndian, header); err != nil {
		return nil, err
	}

	entryIndices := make([]uint32, header.EntryCount)
	if _, err := sr.Seek(int64(header.Header.HeaderSize), io.SeekStart); err != nil {
		return nil, err
	}
	if err := binary.Read(sr, binary.LittleEndian, entryIndices); err != nil {
		return nil, err
	}

	entries := make([]tableEntry, header.EntryCount)
	for i, idx := range entryIndices {
		if idx == noEntry {
			continue
		}
		if _, err := sr.Seek(int64(header.EntriesStart+idx), io.SeekStart); err != nil {
			return nil, err
		}

		var key resTableEntry
		if err := binary.Read(sr, binary.LittleEndian, &key); err != nil {
			return nil, err
		}
		entries[i].key = key

		var value resValue
		if err := binary.Read(sr, binary.LittleEndian, &value); err != nil {
			return nil, err
		}
		entries[i].value = value
	}

	return &tableType{entries, *header}, nil
}
