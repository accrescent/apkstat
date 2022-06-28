package apk

func packLocale(lang []uint8, region []uint8) uint32 {
	return uint32(lang[0])<<24 | uint32(lang[1])<<16 | uint32(region[0])<<8 | uint32(region[1])
}

func dropRegion(packedLocale uint32) uint32 {
	return packedLocale & 0xFFFF0000
}

func hasRegion(packedLocale uint32) bool {
	return packedLocale&0x0000FFFF != 0
}

const (
	scriptLength       = 4
	scriptParentsCount = 5
	packedRoot         = 0 // to represent the root locale
)

func findParent(packedLocale uint32, script []uint8) uint32 {
	if hasRegion(packedLocale) {
		for i := 0; i < scriptParentsCount; i++ {
			// The joys of using Go.
			// https://github.com/golang/go/issues/46505
			if *(*[scriptLength]uint8)(script) == scriptParents()[i].script {
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

// findAncestors finds the ancestors of a locale and fills out with it (assuming out is large enough
// in the process). If any of the members of stopList are seen, they are written to the output but
// the function immediately stops.
//
// findAncestors also outputs the index of the last written ancestor in stopList to stopListIndex,
// which will be -1 if it is not found in stopList.
//
// Returns the number of ancestors written in the output, which is always at least one.
func findAncestors(
	out *[]uint32, stopListIndex *int,
	packedLocale uint32, script []uint8,
	stopList []uint32, stopSetLength int,
) int {
	ancestor := packedLocale
	var count int
	for {
		if out != nil {
			(*out)[count] = ancestor
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

func findDistance(
	supported uint32,
	script []uint8,
	requestAncestors []uint32,
	requestAncestorsCount int,
) int {
	var requestAncestorsIndex int
	supportedAncestorCount := findAncestors(
		nil,
		&requestAncestorsIndex,
		supported,
		script,
		requestAncestors[:],
		requestAncestorsCount,
	)
	// Since both locales share the same root, there will always be a shared ancestor, so the
	// distance in the parent tree is the sum of the distance of 'supported' to the lowest
	// common ancestor (number of ancestors written for 'supported' minus 1) plus the distance
	// of 'request' to the lowest common ancestor (the index of the ancestor in
	// request_ancestors).
	return supportedAncestorCount + requestAncestorsIndex - 1
}

func isRepresentative(languageAndRegion uint32, script []uint8) bool {
	packedLocale :=
		uint64(languageAndRegion)<<32 |
			uint64(script[0])<<24 |
			uint64(script[1])<<16 |
			uint64(script[2])<<8 |
			uint64(script[3])
	_, exists := representativeLocales()[packedLocale]
	return exists
}

const (
	usSpanish            = 0x65735553 // es-US
	mexicanSpanish       = 0x65734D58 // es-MX
	latinAmericanSpanish = 0x6573A424 // es-419
)

// isSpecialSpanish returns whether the locale is a special fallback for es-419. es-US and es-MX are
// considered its equivalent if there is no es-419.
func isSpecialSpanish(languageAndRegion uint32) bool {
	return languageAndRegion == usSpanish || languageAndRegion == mexicanSpanish
}

func localeDataCompareRegions(
	leftRegion []uint8,
	rightRegion []uint8,
	requestedLanguage []uint8,
	requestedScript []uint8,
	requestedRegion []uint8,
) int {
	if leftRegion[0] == rightRegion[0] && leftRegion[1] == rightRegion[1] {
		return 0
	}
	left := packLocale(requestedLanguage, leftRegion)
	right := packLocale(requestedLanguage, rightRegion)
	request := packLocale(requestedLanguage, requestedRegion)

	// If one and only one of the two locales is a special Spanish locale, we replace it with
	// es-419. We don't do the replacement if the other locale is already es-419, or both
	// locales are special Spanish locales (when es-US is being compared to es-MX).
	leftIsSpecialSpanish := isSpecialSpanish(left)
	rightIsSpecialSpanish := isSpecialSpanish(right)
	if leftIsSpecialSpanish && !rightIsSpecialSpanish && right != latinAmericanSpanish {
		left = latinAmericanSpanish
	} else if rightIsSpecialSpanish && !leftIsSpecialSpanish && left != latinAmericanSpanish {
		right = latinAmericanSpanish
	}

	var requestAncestors [maxParentDepth + 1]uint32
	requestAncestorsSlice := requestAncestors[:]
	var leftRightIndex int
	// Find the parents of the request, but stop as soon as we saw left or right.
	leftAndRight := [2]uint32{left, right}
	ancestorCount := findAncestors(
		&requestAncestorsSlice,
		&leftRightIndex,
		request,
		requestedScript,
		leftAndRight[:],
		len(leftAndRight),
	)
	if leftRightIndex == 0 { // We saw left earlier
		return 1
	}
	if leftRightIndex == 1 { // We saw right earlier
		return -1
	}

	// If we are here, neither left nor right are an ancestor of the request. This means that
	// all the ancestors have been computed and the last ancestor is just the language by
	// itself. we will use the distance in the parent tree for determining the better match.
	leftDistance := findDistance(left, requestedScript, requestAncestors[:], ancestorCount)
	rightDistance := findDistance(right, requestedScript, requestAncestors[:], ancestorCount)
	if leftDistance != rightDistance {
		return rightDistance - leftDistance // smaller distance is better
	}

	// If we are here, left and right are equidistant from the request. We will try and see if
	// any of them is a representative locale.
	leftIsRepresentative := isRepresentative(left, requestedScript)
	rightIsRepresentative := isRepresentative(right, requestedScript)
	if leftIsRepresentative != rightIsRepresentative {
		var leftIsRepresentativeVal int
		var rightIsRepresentativeVal int
		if leftIsRepresentative {
			leftIsRepresentativeVal = 1
		} else {
			leftIsRepresentativeVal = 0
		}
		if rightIsRepresentative {
			rightIsRepresentativeVal = 1
		} else {
			rightIsRepresentativeVal = 0
		}

		return leftIsRepresentativeVal - rightIsRepresentativeVal
	}

	// We have no way of figuring out which locale is a better match. For the sake of stability,
	// we consider the locale lwith the lower region code (in dictionary order) better, with
	// two-letter codes before three-digit codes (singe two-letter codes are more specific).
	return int(right - left)
}

func localeDataComputeScript(out *[4]uint8, language []uint8, region []uint8) {
	likelyScripts := likelyScripts()
	scriptCodes := scriptCodes()

	if language[0] == 0 {
		*out = [scriptLength]uint8{}
		return
	}
	lookupKey := packLocale(language, region)
	lookupResult, ok := likelyScripts[lookupKey]
	if !ok {
		// We couldn't find the locale. Let's try without the region.
		if region[0] != 0 {
			lookupKey = dropRegion(lookupKey)
			lookupResult, ok = likelyScripts[lookupKey]
			if ok {
				*out = scriptCodes[lookupResult]
				return
			}
		}
		// We don't know anything about the locale.
		*out = [scriptLength]uint8{}
		return
	} else {
		// We found the locale.
		*out = scriptCodes[lookupResult]
	}
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

func localeDataIsCloseToUSEnglish(region []uint8) bool {
	englishChars := englishChars()
	latinChars := latinChars()
	englishStopList := englishStopList()

	locale := packLocale(englishChars[:], region)
	var stopListIndex int
	findAncestors(nil, &stopListIndex, locale, latinChars[:], englishStopList[:], 2)

	// A locale is like US English if we see "en" before "en-001" in its ancestor list.
	return stopListIndex == 0 // 'en' is first in englishStopList
}
