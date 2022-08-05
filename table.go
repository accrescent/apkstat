package apk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// ResTable is a representation of an Android resource table. It can be referenced from XMLFile
// attributes to resolve resource table references.
type ResTable struct {
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

// NewResTable creates a new ResTable instance from a reader of an Android resource table.
func NewResTable(r io.ReaderAt) (*ResTable, error) {
	f := new(ResTable)

	header := new(resTableHeader)
	sr := io.NewSectionReader(r, 0, maxReadBytes)
	if err := binary.Read(sr, binary.LittleEndian, header); err != nil {
		return nil, err
	}
	if header.Header.Type != resTableChunkType {
		return nil, ErrMalformedHeader
	}

	f.packages = make(map[uint32]*tablePackage)

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
		case resTablePackageType:
			err = f.parseTablePackage(io.NewSectionReader(sr, offset, maxReadBytes-offset))
		default:
			return nil, ErrInvalidChunkType
		}
		if err != nil {
			return nil, err
		}

		offset += int64(chunk.Size)
	}

	return f, nil
}

type resID uint32

// pkg returns the package index of the given resource ID.
func (id resID) pkg() uint32 {
	return uint32(id) >> 24
}

// type_ returns the type index of the given resource ID.
func (id resID) type_() int {
	return (int(id) >> 16) & 0xFF
}

// entry returns the entry index of the given resource ID
func (id resID) entry() int {
	return int(id) & 0xFFFF
}

const sysPackageID = 0x01

// getResource retrieves the value of a resource, resolving string pool references and converting
// other types it encounters to strings as necessary.
func (f *ResTable) getResource(id resID, config *ResTableConfig) (string, error) {
	pkg := id.pkg()
	type_ := id.type_()
	entry := id.entry()

	if type_ < 0 {
		return "", ErrBadIndex
	}

	if pkg == sysPackageID {
		return "", nil
	}

	p := f.packages[pkg]
	if p == nil {
		return "", fmt.Errorf("package 0x%02X not found", id.pkg())
	}

	var best *tableType

	for _, t := range p.tableTypes {
		switch {
		case int(t.header.ID) != type_:
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
		if v.Data != 0 {
			return "true", nil
		} else {
			return "false", nil
		}
	}

	return "", nil
}

// parseTablePackage parses a tablePackage starting at sr and updates the parsing state of f as
// necessary.
func (f *ResTable) parseTablePackage(sr *io.SectionReader) error {
	pkg := new(tablePackage)

	header := new(resTablePackage)
	if err := binary.Read(sr, binary.LittleEndian, header); err != nil {
		return err
	}

	typeSR := io.NewSectionReader(
		sr,
		int64(header.TypeStrings),
		int64(header.Header.Size-header.TypeStrings),
	)
	if typeStrings, err := parseStringPool(typeSR); err != nil {
		return err
	} else {
		pkg.typeStrings = typeStrings
	}

	keySR := io.NewSectionReader(
		sr,
		int64(header.KeyStrings),
		int64(header.Header.Size-header.KeyStrings),
	)
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
			tt, err = f.parseTableType(io.NewSectionReader(sr, offset, int64(chunk.Size)))
			pkg.tableTypes = append(pkg.tableTypes, tt)
		case resTableTypeSpecType:
			// unimplemented
		default:
			return ErrInvalidChunkType
		}
		if err != nil {
			return err
		}

		offset += int64(chunk.Size)
	}

	f.packages[header.ID] = pkg

	return nil
}

// parseTableType parses a resTableType starting at sr and updates the parsing state of f as
// necessary.
func (f *ResTable) parseTableType(sr *io.SectionReader) (*tableType, error) {
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

type resTableHeader struct {
	Header       resChunkHeader
	PackageCount uint32
}

type resTablePackage struct {
	Header         resChunkHeader
	ID             uint32
	Name           [128]uint16
	TypeStrings    uint32
	LastPublicType uint32
	KeyStrings     uint32
	LastPublicKey  uint32
	TypeIDOffset   uint32
}

const (
	densityMedium      = 160
	densityAny         = 0xfffe
	maskKeysHidden     = 0x0003
	keysHiddenNo       = 0x0001
	keysHiddenSoft     = 0x0003
	maskNavHidden      = 0x000c
	maskScreenSize     = 0x0f
	screenSizeNormal   = 0x02
	maskScreenLong     = 0x30
	maskLayoutDir      = 0xc0
	maskUIModeType     = 0x0f
	maskUIModeNight    = 0x30
	maskScreenRound    = 0x03
	maskWideColorGamut = 0x03
	maskHDR            = 0x0c
)

// ResTableConfig describes a particular resource configuration.
type ResTableConfig struct {
	Size                    uint32
	MCC                     uint16
	MNC                     uint16
	Language                [2]uint8
	Country                 [2]uint8
	Orientation             uint8
	Touchscreen             uint8
	Density                 uint16
	Keyboard                uint8
	Navigation              uint8
	InputFlags              uint8
	InputPad0               uint8
	ScreenWidth             uint16
	ScreenHeight            uint16
	SDKVersion              uint16
	MinorVersion            uint16
	ScreenLayout            uint8
	UIMode                  uint8
	SmallestScreenWidthDP   uint16
	ScreenWidthDP           uint16
	ScreenHeightDP          uint16
	LocaleScript            [4]uint8
	LocaleVariant           [8]uint8
	ScreenLayout2           uint8
	ColorMode               uint8
	ScreenConfigPad2        uint16
	LocaleScriptWasComputed bool
	LocaleNumberingSystem   [8]uint8
}

// getImportanceScoleOfLocale returns an integer representing the importance score of the
// configuration locale. Since there isn't a well-specified "importance" order between variants or
// scripts (e.g. we can't easily tell whether "en-Latn-US" is more or less specific than
// "en-US-POSIX"), we arbitrarily decide to give priority to variants over scripts since it seems
// useful to do so. We will consider "en-US-POSIX" more specific than "en-Latn-US."
//
// Unicode extension keywords are considered to be less important than scripts and variants.
func (c ResTableConfig) getImportanceScoreOfLocale() int {
	var x, y, z int
	if c.LocaleVariant[0] != 0 {
		x = 4
	} else {
		x = 0
	}
	if c.LocaleScript[0] != 0 && !c.LocaleScriptWasComputed {
		y = 2
	} else {
		y = 0
	}
	if c.LocaleNumberingSystem[0] != 0 {
		z = 1
	} else {
		z = 0
	}

	return x + y + z
}

// isLocaleMoreSpecificThan returns a positive integer if this config is more specific than o with
// respect to their locales, a negative integer if o is more specific, and 0 if they're equally
// specific.
func (c ResTableConfig) isLocaleMoreSpecificThan(o *ResTableConfig) int {
	if c.Language != [2]uint8{} || c.Country != [2]uint8{} ||
		o.Language != [2]uint8{} || o.Country != [2]uint8{} {
		if c.Language[0] != o.Language[0] {
			if c.Language[0] == 0 {
				return -1
			}
			if o.Language[0] == 0 {
				return -1
			}
		}

		if c.Country[0] != o.Country[0] {
			if c.Country[0] == 0 {
				return -1
			}
			if o.Country[0] == 0 {
				return -1
			}
		}
	}

	return c.getImportanceScoreOfLocale() - o.getImportanceScoreOfLocale()
}

// isMoreSpecificThan returns whether c is more specific than o.
func (c ResTableConfig) isMoreSpecificThan(o *ResTableConfig) bool {
	if o == nil {
		return false
	}

	if c.MCC != 0 || c.MNC != 0 || o.MCC != 0 || o.MNC != 0 {
		if c.MCC != o.MCC {
			if c.MCC == 0 {
				return false
			} else if o.MCC == 0 {
				return true
			}
		}

		if c.MNC != o.MNC {
			if c.MNC == 0 {
				return false
			} else if o.MNC == 0 {
				return true
			}
		}
	}

	if c.Language != [2]uint8{} || c.Country != [2]uint8{} ||
		o.Language != [2]uint8{} || o.Country != [2]uint8{} {
		diff := c.isLocaleMoreSpecificThan(o)
		if diff < 0 {
			return false
		} else if diff > 0 {
			return true
		}
	}

	if c.ScreenLayout != 0 || o.ScreenLayout != 0 {
		if (c.ScreenLayout^o.ScreenLayout)&maskLayoutDir != 0 {
			if c.ScreenLayout&maskLayoutDir == 0 {
				return false
			} else if o.ScreenLayout&maskLayoutDir == 0 {
				return true
			}
		}
	}

	if c.SmallestScreenWidthDP != 0 || o.SmallestScreenWidthDP != 0 {
		if c.SmallestScreenWidthDP != o.SmallestScreenWidthDP {
			if c.SmallestScreenWidthDP == 0 {
				return false
			} else if o.SmallestScreenWidthDP == 0 {
				return true
			}
		}
	}

	if c.ScreenWidthDP != 0 || c.ScreenHeightDP != 0 ||
		o.ScreenWidthDP != 0 || o.ScreenHeightDP != 0 {
		if c.ScreenWidthDP != o.ScreenWidthDP {
			if c.ScreenWidthDP == 0 {
				return false
			} else if o.ScreenWidthDP == 0 {
				return true
			}
		}

		if c.ScreenHeightDP != o.ScreenHeightDP {
			if c.ScreenHeightDP == 0 {
				return false
			} else if o.ScreenHeightDP == 0 {
				return true
			}
		}
	}

	if c.ScreenLayout != 0 || o.ScreenLayout != 0 {
		if (c.ScreenLayout^o.ScreenLayout)&maskScreenSize != 0 {
			if c.ScreenLayout&maskScreenSize == 0 {
				return false
			} else if o.ScreenLayout&maskScreenSize == 0 {
				return true
			}
		}
		if (c.ScreenLayout^o.ScreenLayout)&maskScreenLong != 0 {
			if c.ScreenLayout&maskScreenLong == 0 {
				return false
			} else if o.ScreenLayout&maskScreenLong == 0 {
				return true
			}
		}
	}

	if c.ScreenLayout2 != 0 || o.ScreenLayout2 != 0 {
		if (c.ScreenLayout2^o.ScreenLayout2)&maskScreenRound != 0 {
			if c.ScreenLayout2&maskScreenRound == 0 {
				return false
			} else if o.ScreenLayout2&maskScreenRound == 0 {
				return true
			}
		}
	}

	if c.ColorMode != 0 || o.ColorMode != 0 {
		if (c.ColorMode^o.ColorMode)&maskHDR != 0 {
			if c.ColorMode&maskHDR == 0 {
				return false
			} else if o.ColorMode&maskHDR == 0 {
				return true
			}
		}
		if (c.ColorMode^o.ColorMode)&maskWideColorGamut != 0 {
			if c.ColorMode&maskWideColorGamut == 0 {
				return false
			} else if o.ColorMode&maskWideColorGamut == 0 {
				return true
			}
		}
	}

	if c.Orientation != o.Orientation {
		if c.Orientation == 0 {
			return false
		} else if o.Orientation == 0 {
			return true
		}
	}

	if c.UIMode != 0 || o.UIMode != 0 {
		if (c.UIMode^o.UIMode)&maskUIModeType != 0 {
			if c.UIMode&maskUIModeType == 0 {
				return false
			} else if o.UIMode&maskUIModeType == 0 {
				return true
			}
		}
		if (c.UIMode^o.UIMode)&maskUIModeNight != 0 {
			if c.UIMode&maskUIModeNight == 0 {
				return false
			} else if o.UIMode&maskUIModeNight == 0 {
				return true
			}
		}
	}

	if c.Touchscreen != o.Touchscreen {
		if c.Touchscreen == 0 {
			return false
		} else if o.Touchscreen == 0 {
			return true
		}
	}

	if c.Keyboard != 0 || c.Navigation != 0 || c.InputFlags != 0 || c.InputPad0 != 0 ||
		o.Keyboard != 0 || o.Navigation != 0 || o.InputFlags != 0 || o.InputPad0 != 0 {
		if (c.InputFlags&o.InputFlags)&maskKeysHidden != 0 {
			if c.InputFlags&maskKeysHidden == 0 {
				return false
			} else if o.InputFlags&maskKeysHidden == 0 {
				return true
			}
		}

		if (c.InputFlags&o.InputFlags)&maskNavHidden != 0 {
			if c.InputFlags&maskNavHidden == 0 {
				return false
			} else if o.InputFlags&maskNavHidden == 0 {
				return true
			}
		}

		if c.Keyboard != o.Keyboard {
			if c.Keyboard == 0 {
				return false
			} else if o.Keyboard == 0 {
				return true
			}
		}

		if c.Navigation != o.Navigation {
			if c.Navigation == 0 {
				return false
			} else if o.Navigation == 0 {
				return true
			}
		}
	}

	if c.ScreenWidth != 0 || c.ScreenHeight != 0 || o.ScreenWidth != 0 || o.ScreenHeight != 0 {
		if c.ScreenWidth != o.ScreenWidth {
			if c.ScreenWidth == 0 {
				return false
			} else if o.ScreenWidth == 0 {
				return true
			}
		} else if c.ScreenHeight != o.ScreenHeight {
			if c.ScreenHeight == 0 {
				return false
			} else if o.ScreenHeight == 0 {
				return true
			}
		}
	}

	if c.SDKVersion != 0 || c.MinorVersion != 0 || o.SDKVersion != 0 || o.MinorVersion != 0 {
		if c.SDKVersion != o.SDKVersion {
			if c.SDKVersion == 0 {
				return false
			} else if o.SDKVersion == 0 {
				return true
			}
		} else if c.MinorVersion != o.MinorVersion {
			if c.MinorVersion == 0 {
				return false
			} else if o.MinorVersion == 0 {
				return true
			}
		}
	}

	return false
}

// Codes for specially handles languages and regions

func english() [2]uint8 {
	return [2]uint8{'e', 'n'} // packed version of "en"
}

func unitedStates() [2]uint8 {
	return [2]uint8{'U', 'S'} // packed version of "US"
}

func filipino() [2]uint8 {
	return [2]uint8{'\xAD', '\x05'} // packed version of "fil"
}

func tagalog() [2]uint8 {
	return [2]uint8{'t', 'l'} // packed version of "tl"
}

func langsAreEquivalent(lang1 [2]uint8, lang2 [2]uint8) bool {
	return lang1 == lang2 ||
		lang1 == tagalog() && lang2 == filipino() ||
		lang1 == filipino() && lang2 == tagalog()
}

// isLocaleBetterThan returns whether c is a better locale match than o for the requested
// configuration r. Similar to isBetterThan, this assumes that match has already been used to remove
// any configurations that don't match the requested configuration at all.
func (c ResTableConfig) isLocaleBetterThan(o, r *ResTableConfig) bool {
	if r.Language == [2]uint8{} && r.Country == [2]uint8{} {
		return false
	}

	if r.Language == [2]uint8{} && r.Country == [2]uint8{} &&
		o.Language == [2]uint8{} && o.Country == [2]uint8{} {
		return false
	}

	if !langsAreEquivalent(c.Language, o.Language) {
		if r.Language == english() {
			if r.Country == unitedStates() {
				if c.Language[0] != 0 {
					return c.Country[0] == 0 || c.Country == unitedStates()
				} else {
					return !(o.Country[0] == 0 || o.Country == unitedStates())
				}
			}
		} else if localeDataIsCloseToUSEnglish(r.Country[:]) {
			if c.Language[0] != 0 {
				return localeDataIsCloseToUSEnglish(c.Country[:])
			} else {
				return !localeDataIsCloseToUSEnglish(o.Country[:])
			}
		}
		return c.Language[0] != 0
	}

	// If we are here, both the resources have an equivalent non-empty language
	// to the request.
	//
	// Because the languages are equivalent, computeScript() always returns a
	// non-empty script for languages it knows about, and we have passed the
	// script checks in match(), the scripts are either all unknown or are all
	// the same. So we can't gain anything by checking the scripts. We need to
	// check the region and variant.

	// See if any of the regions is better than the other.
	regionComparison := localeDataCompareRegions(
		c.Country[:],
		o.Country[:],
		r.Language[:],
		r.LocaleScript[:],
		r.Country[:],
	)
	if regionComparison != 0 {
		return regionComparison > 0
	}

	// The regions are the same. Try the variant.
	localeMatches := c.LocaleVariant == r.LocaleVariant
	otherMatches := o.LocaleVariant == r.LocaleVariant
	if localeMatches != otherMatches {
		return localeMatches
	}

	// The variants are the same. Try the numbering system.
	localeNumsysMatches := c.LocaleNumberingSystem == r.LocaleNumberingSystem
	otherNumsysMatches := o.LocaleNumberingSystem == r.LocaleNumberingSystem
	if localeNumsysMatches != otherNumsysMatches {
		return localeNumsysMatches
	}

	// Finally, the languages, although equivalent, may still be different (like for Tagalog and
	// Filipino). Identical is better than just equivalent.
	if c.Language == r.Language && o.Language != r.Language {
		return true
	}

	return false
}

// isBetterThan returns whether c is a better match than o for the requested configuration r. It
// assumes that match has already been used to remove any configurations that don't match the
// requested configuration at all; if they are not first filtered, non-matching results can be
// considered better than matching ones.
//
// The general rule per attribute is as follows: if the request cares about an attribute (it
// normally does), it's a tie if c and o are equal. If they are not equal then one must be generic
// because only generic and '== r' will pass the match() call. So if c is not generic, it wins. If c
// _is_ generic, o wins and isBetterThan returns false.
func (c ResTableConfig) isBetterThan(o *ResTableConfig, r *ResTableConfig) bool {
	switch {
	case r == nil:
		return c.isMoreSpecificThan(o)
	case o == nil:
		return false
	}

	if c.MCC != 0 || c.MNC != 0 || o.MCC != 0 || o.MNC != 0 {
		if c.MCC != o.MCC && r.MCC != 0 {
			return c.MCC != 0
		}

		if c.MNC != o.MNC && r.MNC != 0 {
			return c.MNC != 0
		}
	}

	if c.isLocaleBetterThan(o, r) {
		return true
	}

	if c.ScreenLayout != 0 || r.ScreenLayout != 0 {
		if (c.ScreenLayout^o.ScreenLayout)&maskLayoutDir != 0 &&
			r.ScreenLayout&maskLayoutDir != 0 {
			myLayoutDir := c.ScreenLayout & maskLayoutDir
			oLayoutDir := o.ScreenLayout & maskLayoutDir
			return myLayoutDir > oLayoutDir
		}
	}

	if c.SmallestScreenWidthDP != 0 || o.SmallestScreenWidthDP != 0 {
		if c.SmallestScreenWidthDP != o.SmallestScreenWidthDP {
			return c.SmallestScreenWidthDP > o.SmallestScreenWidthDP
		}
	}

	if c.ScreenWidthDP != 0 || c.ScreenHeightDP != 0 || o.ScreenWidthDP != 0 || o.ScreenHeightDP != 0 {
		myDelta, otherDelta := 0, 0
		if r.ScreenWidthDP != 0 {
			myDelta += int(r.ScreenWidthDP - c.ScreenWidthDP)
			otherDelta += int(r.ScreenWidthDP - o.ScreenWidthDP)
		}
		if r.ScreenHeightDP != 0 {
			myDelta += int(r.ScreenHeightDP - c.ScreenHeightDP)
			otherDelta += int(r.ScreenHeightDP - o.ScreenHeightDP)
		}
		if myDelta != otherDelta {
			return myDelta < otherDelta
		}
	}

	if c.ScreenLayout != 0 || o.ScreenLayout != 0 {
		if (c.ScreenLayout^o.ScreenLayout)&maskScreenSize != 0 &&
			r.ScreenLayout&maskScreenSize != 0 {
			mySL := c.ScreenLayout & maskScreenSize
			oSL := o.ScreenLayout & maskScreenSize
			fixedMySL := mySL
			fixedOSL := oSL
			if r.ScreenLayout&maskScreenSize >= screenSizeNormal {
				if fixedMySL == 0 {
					fixedMySL = screenSizeNormal
				}
				if fixedOSL == 0 {
					fixedOSL = screenSizeNormal
				}
			}

			if fixedMySL == fixedOSL {
				return mySL != 0
			} else {
				return fixedMySL > fixedOSL
			}
		}
		if (c.ScreenLayout^o.ScreenLayout)&maskScreenLong != 0 &&
			r.ScreenLayout&maskScreenLong != 0 {
			return c.ScreenLayout&maskScreenLong != 0
		}
	}

	if c.ScreenLayout2 != 0 || o.ScreenLayout2 != 0 {
		if (c.ScreenLayout2^o.ScreenLayout2)&maskScreenRound != 0 &&
			r.ScreenLayout2&maskScreenRound != 0 {
			return c.ScreenLayout2&maskScreenRound != 0
		}
	}

	if c.ColorMode != 0 || o.ColorMode != 0 {
		if (c.ColorMode^o.ColorMode)&maskWideColorGamut != 0 &&
			r.ColorMode&maskWideColorGamut != 0 {
			return c.ColorMode&maskWideColorGamut != 0
		}
		if (c.ColorMode^o.ColorMode)&maskHDR != 0 &&
			r.ColorMode&maskHDR != 0 {
			return c.ColorMode&maskHDR != 0
		}
	}

	if c.Orientation != o.Orientation && r.Orientation != 0 {
		return c.Orientation != 0
	}

	if c.UIMode != 0 && o.UIMode != 0 {
		if (c.UIMode^o.UIMode)&maskUIModeType != 0 &&
			r.UIMode&maskUIModeType != 0 {
			return c.UIMode&maskUIModeType != 0
		}
		if (c.UIMode^o.UIMode)&maskUIModeNight != 0 &&
			r.UIMode&maskUIModeNight != 0 {
			return c.UIMode&maskUIModeNight != 0
		}
	}

	if c.Orientation != 0 || c.Touchscreen != 0 || c.Density != 0 ||
		o.Orientation != 0 || o.Touchscreen != 0 || o.Density != 0 {
		if c.Density != o.Density {
			var thisDensity int
			if c.Density != 0 {
				thisDensity = int(c.Density)
			} else {
				thisDensity = densityMedium
			}
			var otherDensity int
			if o.Density != 0 {
				otherDensity = int(o.Density)
			} else {
				otherDensity = densityMedium
			}

			if thisDensity == densityAny {
				return true
			} else if otherDensity == densityAny {
				return false
			}

			requestedDensity := int(r.Density)
			if r.Density == 0 || r.Density == densityAny {
				requestedDensity = densityMedium
			}

			h := thisDensity
			l := otherDensity
			imBigger := true
			if l > h {
				t := h
				h = l
				l = t
				imBigger = false
			}

			if requestedDensity >= h {
				return imBigger
			}
			if l >= requestedDensity {
				return !imBigger
			}
			if ((2*l)-requestedDensity)*h > requestedDensity*requestedDensity {
				return !imBigger
			} else {
				return imBigger
			}
		}

		if c.Touchscreen != o.Touchscreen && r.Touchscreen != 0 {
			return c.Touchscreen != 0
		}
	}

	if c.Keyboard != 0 || c.Navigation != 0 || c.InputFlags != 0 || c.InputPad0 != 0 ||
		o.Keyboard != 0 || o.Navigation != 0 || o.InputFlags != 0 || o.InputPad0 != 0 {
		keysHidden := c.InputFlags & maskKeysHidden
		oKeysHidden := o.InputFlags & maskKeysHidden
		if keysHidden != oKeysHidden {
			reqKeysHidden := r.InputFlags & maskKeysHidden
			if reqKeysHidden != 0 {
				switch {
				case keysHidden == 0:
					return false
				case oKeysHidden == 0:
					return true
				case reqKeysHidden == keysHidden:
					return true
				case reqKeysHidden == oKeysHidden:
					return false
				}
			}
		}

		navHidden := c.InputFlags & maskNavHidden
		oNavHidden := o.InputFlags & maskNavHidden
		if navHidden != oNavHidden {
			reqNavHidden := r.InputFlags & maskNavHidden
			if reqNavHidden != 0 {
				if navHidden == 0 {
					return false
				} else if oNavHidden == 0 {
					return true
				}
			}
		}

		if c.Keyboard != o.Keyboard && r.Keyboard != 0 {
			return c.Keyboard != 0
		}

		if c.Navigation != o.Navigation && r.Navigation != 0 {
			return c.Navigation != 0
		}
	}

	if c.ScreenWidth != 0 || c.ScreenHeight != 0 || o.ScreenWidth != 0 || o.ScreenHeight != 0 {
		myDelta, otherDelta := 0, 0
		if r.ScreenWidth != 0 {
			myDelta += int(r.ScreenWidth - c.ScreenWidth)
			otherDelta += int(r.ScreenWidth - o.ScreenWidth)
		}
		if r.ScreenHeight != 0 {
			myDelta += int(r.ScreenHeight - c.ScreenHeight)
			otherDelta += int(r.ScreenHeight - o.ScreenHeight)
		}
		if myDelta != otherDelta {
			return myDelta < otherDelta
		}
	}

	if c.SDKVersion != 0 || c.MinorVersion != 0 || o.SDKVersion != 0 || o.MinorVersion != 0 {
		if c.SDKVersion != o.SDKVersion && r.SDKVersion != 0 {
			return c.SDKVersion > o.SDKVersion
		}

		if c.MinorVersion != o.MinorVersion && r.MinorVersion != 0 {
			return c.MinorVersion != 0
		}
	}

	return false
}

// match returns whether c can be considered a match for the request parameters in settings.
//
// Note this is asymetric. A default piece of data will match every request, but a request for the
// default should not match odd specifics (i.e. a request with no MCC should not match a particular
// MCC's data).
func (c ResTableConfig) match(settings *ResTableConfig) bool {
	if settings == nil {
		return true
	}

	if c.MCC != 0 || c.MNC != 0 {
		if c.MCC != 0 && c.MCC != settings.MCC {
			return false
		}
		if c.MNC != 0 && c.MNC != settings.MNC {
			return false
		}
	}

	if c.Language != [2]uint8{0, 0} {
		if !langsAreEquivalent(c.Language, settings.Language) {
			return false
		}

		// For backward compatibility and supporting private-use locales, we fall back to
		// old behavior if we couldn't determine the script for either of the desired locale
		// or the provided locale. But if if we could determine the scripts, they should be
		// the same for the locales to match.
		countriesMustMatch := false
		computedScript := [4]uint8{}
		script := []uint8{}
		if settings.LocaleScript[0] == 0 { // could not determine the request's script
			countriesMustMatch = true
		} else {
			if c.LocaleScript[0] == 0 && !c.LocaleScriptWasComputed {
				// script was not provided or computed, so we try to compute it
				localeDataComputeScript(&computedScript, c.Language[:], c.Country[:])
				if computedScript[0] == 0 { // we could not compute the script
					countriesMustMatch = true
				} else {
					script = computedScript[:]
				}
			} else { // script was provided, so just use it
				script = c.LocaleScript[:]
			}
		}

		if countriesMustMatch {
			if c.Country[0] == 0 && c.Country != settings.Country {
				return false
			}
		} else {
			if !bytes.Equal(script, settings.LocaleScript[:]) {
				return false
			}
		}
	}

	if c.ScreenLayout != 0 || c.UIMode != 0 || c.SmallestScreenWidthDP != 0 {
		layoutDir := c.ScreenLayout & maskLayoutDir
		setLayoutDir := settings.ScreenLayout & maskLayoutDir
		if layoutDir != 0 && layoutDir != setLayoutDir {
			return false
		}

		screenSize := c.ScreenLayout & maskScreenSize
		setScreenSize := settings.ScreenLayout & maskScreenSize
		if screenSize != 0 && screenSize > setScreenSize {
			return false
		}

		screenLong := c.ScreenLayout & maskScreenLong
		setScreenLong := settings.ScreenLayout & maskScreenLong
		if screenLong != 0 && screenLong != setScreenLong {
			return false
		}

		uiModeType := c.UIMode & maskUIModeType
		setUIModeType := settings.UIMode & maskUIModeType
		if uiModeType != 0 && uiModeType != setUIModeType {
			return false
		}

		uiModeNight := c.UIMode & maskUIModeNight
		setUIModeNight := settings.UIMode & maskUIModeNight
		if uiModeNight != 0 && uiModeNight != setUIModeNight {
			return false
		}

		if c.SmallestScreenWidthDP != 0 && c.SmallestScreenWidthDP > settings.SmallestScreenWidthDP {
			return false
		}
	}

	if c.ScreenLayout2 != 0 || c.ColorMode != 0 || c.ScreenConfigPad2 != 0 {
		screenRound := c.ScreenLayout2 & maskScreenRound
		setScreenRound := settings.ScreenLayout2 & maskScreenRound
		if screenRound != 0 && screenRound != setScreenRound {
			return false
		}

		hdr := c.ColorMode & maskHDR
		setHDR := settings.ColorMode & maskHDR
		if hdr != 0 && hdr != setHDR {
			return false
		}

		wideColorGamut := c.ColorMode & maskWideColorGamut
		setWideColorGamut := settings.ColorMode & maskWideColorGamut
		if wideColorGamut != 0 && wideColorGamut != setWideColorGamut {
			return false
		}
	}

	if c.ScreenWidthDP != 0 || c.ScreenHeightDP != 0 {
		if c.ScreenWidthDP != 0 && c.ScreenWidthDP > settings.ScreenWidthDP {
			return false
		}
		if c.ScreenHeightDP != 0 && c.ScreenHeightDP > settings.ScreenHeightDP {
			return false
		}
	}
	if c.Orientation != 0 || c.Touchscreen != 0 || c.Density != 0 { // screen type
		if c.Orientation != 0 && c.Orientation != settings.Orientation {
			return false
		}
		if c.Touchscreen != 0 && c.Touchscreen != settings.Touchscreen {
			return false
		}
	}
	if c.Keyboard != 0 || c.Navigation != 0 || c.InputFlags != 0 || c.InputPad0 != 0 { // input
		keysHidden := c.InputFlags & maskKeysHidden
		setKeysHidden := settings.InputFlags & maskKeysHidden
		if keysHidden != 0 && keysHidden != setKeysHidden {
			if keysHidden != keysHiddenNo || setKeysHidden != keysHiddenSoft {
				return false
			}
		}
		navHidden := c.InputFlags & maskNavHidden
		setNavHidden := settings.InputFlags & maskNavHidden
		if navHidden != 0 && navHidden != setNavHidden {
			return false
		}
		if c.Keyboard != 0 && c.Keyboard != settings.Keyboard {
			return false
		}
		if c.Navigation != 0 && c.Navigation != settings.Navigation {
			return false
		}
	}
	if c.ScreenWidth != 0 || c.ScreenHeight != 0 { // screen size
		if c.ScreenWidth != 0 && c.ScreenWidth > settings.ScreenWidth {
			return false
		}
		if c.ScreenHeight != 0 && c.ScreenHeight > settings.ScreenHeight {
			return false
		}
	}
	if c.SDKVersion != 0 || c.MinorVersion != 0 { // version
		if c.SDKVersion != 0 && c.SDKVersion > settings.SDKVersion {
			return false
		}
		if c.MinorVersion != 0 && c.MinorVersion != settings.MinorVersion {
			return false
		}
	}

	return true
}

const noEntry = 0xFFFFFFFF

type resTableType struct {
	Header       resChunkHeader
	ID           uint8
	Flags        uint8
	Reserved     uint16
	EntryCount   uint32
	EntriesStart uint32
	Config       ResTableConfig
}

type resTableEntry struct {
	Size  uint16
	Flags uint16
	Key   resStringPoolRef
}
