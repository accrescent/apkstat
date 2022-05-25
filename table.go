package apkstat

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
		}
		return c.Language[0] != 0
	}

	return false
}

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

	return 0
}

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