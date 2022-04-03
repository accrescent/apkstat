package apkstat

type Manifest struct {
	Package         string           `xml:"package,attr"`
	VersionCode     int32            `xml:"http://schemas.android.com/apk/res/android versionCode,attr"`
	VersionName     string           `xml:"http://schemas.android.com/apk/res/android versionName,attr"`
	UsesPermissions []UsesPermission `xml:"uses-permission"`
	UsesSDK         UsesSDK          `xml:"uses-sdk"`
}

type UsesPermission struct {
	Name string `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

type UsesSDK struct {
	MinSDKVersion    uint `xml:"http://schemas.android.com/apk/res/android minSdkVersion,attr"`
	TargetSDKVersion uint `xml:"http://schemas.android.com/apk/res/android targetSdkVersion,attr"`
	MaxSDKVersion    uint `xml:"http://schemas.android.com/apk/res/android maxSdkVersion,attr"`
}
