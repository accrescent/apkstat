package apkstat

func packLocale(lang [2]uint8, region [2]uint8) uint32 {
	return uint32(lang[0])<<24 | uint32(lang[1])<<16 | uint32(region[0])<<8 | uint32(region[1])
}

func dropRegion(packedLocale uint32) uint32 {
	return packedLocale & 0xFFFF0000
}

func hasRegion(packedLocale uint32) bool {
	return packedLocale&0x0000FFFF != 0
}

const scriptLength = 4
const scriptParentsCount = 5
const packedRoot = 0

func findParent(packedLocale uint32, script []uint8) uint32 {
	if hasRegion(packedLocale) {
		for i := 0; i < scriptParentsCount; i++ {
			// The joys of using Go.
			// https://github.com/golang/go/issues/46505
			if *(*[4]uint8)(script) == scriptParents()[i].script {
				map_ := scriptParents()[i].map_
				lookupResult, exists := map_[packedLocale]
				if exists {
					return lookupResult
				}
				break
			}
		}
		return dropRegion(packedLocale)
	}
	return packedRoot
}

func findAncestors(out []uint32, stopListIndex *int,
	packedLocale uint32, script []uint8,
	stopList []uint32, stopSetLength int) int {
	ancestor := packedLocale
	var count int
	for {
		if out != nil {
			out[count] = ancestor
		}
		count++
		for i := 0; i < stopSetLength; i++ {
			if stopList[i] == ancestor {
				*stopListIndex = i
				return count
			}
		}
		ancestor = findParent(ancestor, script)

		if ancestor == packedRoot {
			break
		}

	}
	*stopListIndex = -1
	return count
}

func englishStopList() [2]uint32 {
	return [2]uint32{
		0x656E0000, // en
		0x656E8400, // en-001
	}
}

func englishChars() [2]uint8 {
	return [2]uint8{'e', 'n'}
}

func latinChars() [4]uint8 {
	return [4]uint8{'L', 'a', 't', 'n'}
}

func localeDataIsCloseToUSEnglish(region [2]uint8) bool {
	latinChars := latinChars()
	englishStopList := englishStopList()

	locale := packLocale(englishChars(), region)
	var stopListIndex int
	findAncestors(nil, &stopListIndex, locale, latinChars[:], englishStopList[:], 2)

	return stopListIndex == 0
}
