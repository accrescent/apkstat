package apkstat

type Manifest struct {
	Package         string           `xml:"package,attr"`
	VersionCode     int32            `xml:"http://schemas.android.com/apk/res/android versionCode,attr"`
	VersionName     string           `xml:"http://schemas.android.com/apk/res/android versionName,attr"`
	Application     Application      `xml:"application"`
	UsesPermissions []UsesPermission `xml:"uses-permission"`
	UsesSDK         UsesSDK          `xml:"uses-sdk"`
}

type Application struct {
	AllowTaskReparenting         bool       `xml:"http://schemas.android.com/apk/res/android allowTaskReparenting,attr"`
	AllowBackup                  bool       `xml:"http://schemas.android.com/apk/res/android allowBackup,attr"`
	BackupAgent                  string     `xml:"http://schemas.android.com/apk/res/android backupAgent,attr"`
	BackupInForeground           bool       `xml:"http://schemas.android.com/apk/res/android backupInForeground,attr"`
	DataExtractionRules          string     `xml:"http://schemas.android.com/apk/res/android dataExtractionRules,attr"`
	Debuggable                   bool       `xml:"http://schemas.android.com/apk/res/android debuggable,attr"`
	Label                        string     `xml:"http://schemas.android.com/apk/res/android label,attr"`
	ManageSpaceActivity          string     `xml:"http://schemas.android.com/apk/res/android manageSpaceActivity,attr"`
	Name                         string     `xml:"http://schemas.android.com/apk/res/android name,attr"`
	NetworkSecurityConfig        string     `xml:"http://schemas.android.com/apk/res/android networkSecurityConfig,attr"`
	RequestLegacyExternalStorage bool       `xml:"http://schemas.android.com/apk/res/android requestLegacyExternalStorage,attr"`
	SupportsRTL                  bool       `xml:"http://schemas.android.com/apk/res/android supportsRtl,attr"`
	TestOnly                     bool       `xml:"http://schemas.android.com/apk/res/android testOnly,attr"`
	UsesCleartextTraffic         bool       `xml:"http://schemas.android.com/apk/res/android usesCleartextTraffic,attr"`
	Activities                   []Activity `xml:"activity"`
}

type Activity struct {
	Exported      string         `xml:"http://schemas.android.com/apk/res/android exported,attr"`
	Label         string         `xml:"http://schemas.android.com/apk/res/android label,attr"`
	Name          string         `xml:"http://schemas.android.com/apk/res/android name,attr"`
	IntentFilters []IntentFilter `xml:"intent-filter"`
}

type IntentFilter struct {
	Priority   int        `xml:"http://schemas.android.com/apk/res/android priority,attr"`
	Order      int        `xml:"http://schemas.android.com/apk/res/android order,attr"`
	AutoVerify bool       `xml:"http://schemas.android.com/apk/res/android autoVerify,attr"`
	Actions    []Action   `xml:"action"`
	Categories []Category `xml:"category"`
	Data       []Data     `xml:"data"`
}

type Action struct {
	Name string `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

type Category struct {
	Name string `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

type Data struct {
	Scheme      string `xml:"http://schemas.android.com/apk/res/android scheme,attr"`
	Host        string `xml:"http://schemas.android.com/apk/res/android host,attr"`
	Port        string `xml:"http://schemas.android.com/apk/res/android port,attr"`
	Path        string `xml:"http://schemas.android.com/apk/res/android path,attr"`
	PathPattern string `xml:"http://schemas.android.com/apk/res/android pathPattern,attr"`
	PathPrefix  string `xml:"http://schemas.android.com/apk/res/android pathPrefix,attr"`
	MimeType    string `xml:"http://schemas.android.com/apk/res/android mimeType,attr"`
}

type UsesPermission struct {
	Name string `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

type UsesSDK struct {
	MinSDKVersion    uint `xml:"http://schemas.android.com/apk/res/android minSdkVersion,attr"`
	TargetSDKVersion uint `xml:"http://schemas.android.com/apk/res/android targetSdkVersion,attr"`
	MaxSDKVersion    uint `xml:"http://schemas.android.com/apk/res/android maxSdkVersion,attr"`
}
