package apkstat

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
