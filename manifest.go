package apkstat

type Manifest struct {
	Package         string           `xml:"package,attr"`
	VersionCode     int32            `xml:"http://schemas.android.com/apk/res/android versionCode,attr"`
	VersionName     string           `xml:"http://schemas.android.com/apk/res/android versionName,attr"`
	UsesPermissions []UsesPermission `xml:"uses-permission"`
}

type UsesPermission struct {
	Name string `xml:"http://schemas.android.com/apk/res/android name,attr"`
}
