package apk

// Manifest is the root element of an AndroidManifest.xml file.
type Manifest struct {
	Package         string            `xml:"package,attr"`
	VersionCode     int32             `xml:"http://schemas.android.com/apk/res/android versionCode,attr"`
	VersionName     string            `xml:"http://schemas.android.com/apk/res/android versionName,attr"`
	Application     Application       `xml:"application"`
	Queries         *[]Query          `xml:"queries"`
	SupportsScreens *[]SupportsScreen `xml:"supports-screens"`
	UsesPermissions *[]UsesPermission `xml:"uses-permission"`
	UsesSDK         *UsesSDK          `xml:"uses-sdk"`
}

// Application is a declaration of an application.
type Application struct {
	AllowTaskReparenting         *bool            `xml:"http://schemas.android.com/apk/res/android allowTaskReparenting,attr"`
	AllowBackup                  *bool            `xml:"http://schemas.android.com/apk/res/android allowBackup,attr"`
	BackupAgent                  *string          `xml:"http://schemas.android.com/apk/res/android backupAgent,attr"`
	BackupInForeground           *bool            `xml:"http://schemas.android.com/apk/res/android backupInForeground,attr"`
	DataExtractionRules          *string          `xml:"http://schemas.android.com/apk/res/android dataExtractionRules,attr"`
	Debuggable                   *bool            `xml:"http://schemas.android.com/apk/res/android debuggable,attr"`
	Label                        *string          `xml:"http://schemas.android.com/apk/res/android label,attr"`
	ManageSpaceActivity          *string          `xml:"http://schemas.android.com/apk/res/android manageSpaceActivity,attr"`
	Name                         *string          `xml:"http://schemas.android.com/apk/res/android name,attr"`
	NetworkSecurityConfig        *string          `xml:"http://schemas.android.com/apk/res/android networkSecurityConfig,attr"`
	RequestLegacyExternalStorage *bool            `xml:"http://schemas.android.com/apk/res/android requestLegacyExternalStorage,attr"`
	SupportsRTL                  *bool            `xml:"http://schemas.android.com/apk/res/android supportsRtl,attr"`
	TestOnly                     *bool            `xml:"http://schemas.android.com/apk/res/android testOnly,attr"`
	UsesCleartextTraffic         *bool            `xml:"http://schemas.android.com/apk/res/android usesCleartextTraffic,attr"`
	Activities                   *[]Activity      `xml:"activity"`
	ActivityAliases              *[]ActivityAlias `xml:"activity-alias"`
	MetaData                     *[]MetaData      `xml:"meta-data"`
	Services                     *[]Service       `xml:"service"`
	Receivers                    *[]Receiver      `xml:"receiver"`
	Profileable                  *Profileable     `xml:"profileable"`
	Providers                    *[]Provider      `xml:"provider"`
}

// Activity is an activity that implements part of an application's visual user interface.
type Activity struct {
	Exported      *string         `xml:"http://schemas.android.com/apk/res/android exported,attr"`
	Label         *string         `xml:"http://schemas.android.com/apk/res/android label,attr"`
	Name          string          `xml:"http://schemas.android.com/apk/res/android name,attr"`
	IntentFilters *[]IntentFilter `xml:"intent-filter"`
	MetaData      *[]MetaData     `xml:"meta-data"`
}

// IntentFilter specifies the types of intents that an activity, service, or broadcast receiver can
// respond to.
type IntentFilter struct {
	Priority   *uint32     `xml:"http://schemas.android.com/apk/res/android priority,attr"`
	Order      *int32      `xml:"http://schemas.android.com/apk/res/android order,attr"`
	AutoVerify *bool       `xml:"http://schemas.android.com/apk/res/android autoVerify,attr"`
	Actions    []Action    `xml:"action"`
	Categories *[]Category `xml:"category"`
	Data       *[]Data     `xml:"data"`
}

// Action adds an action to an intent filter.
type Action struct {
	Name string `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

// Category adds a category name to an intent filter.
type Category struct {
	Name string `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

// Data adds a data specification to an intent filter. A specification can be just a data type, just
// a URI, or both a data type and a URI.
type Data struct {
	Scheme      *string `xml:"http://schemas.android.com/apk/res/android scheme,attr"`
	Host        *string `xml:"http://schemas.android.com/apk/res/android host,attr"`
	Port        *string `xml:"http://schemas.android.com/apk/res/android port,attr"`
	Path        *string `xml:"http://schemas.android.com/apk/res/android path,attr"`
	PathPattern *string `xml:"http://schemas.android.com/apk/res/android pathPattern,attr"`
	PathPrefix  *string `xml:"http://schemas.android.com/apk/res/android pathPrefix,attr"`
	MimeType    *string `xml:"http://schemas.android.com/apk/res/android mimeType,attr"`
}

// MetaData is a name-value pair for an item of additional, arbitrary data that can be supplied to
// the parent component. A component element can contain any number of MetaData subelements.
type MetaData struct {
	Name     string  `xml:"http://schemas.android.com/apk/res/android name,attr"`
	Resource *string `xml:"http://schemas.android.com/apk/res/android resource,attr"`
	Value    *string `xml:"http://schemas.android.com/apk/res/android value,attr"`
}

// ActivityAlias is an alias for an activity named by the TargetActivity attribute.
type ActivityAlias struct {
	Enabled        *bool           `xml:"http://schemas.android.com/apk/res/android enabled,attr"`
	Exported       bool            `xml:"http://schemas.android.com/apk/res/android exported,attr"`
	Label          *string         `xml:"http://schemas.android.com/apk/res/android label,attr"`
	Name           string          `xml:"http://schemas.android.com/apk/res/android name,attr"`
	Permission     *string         `xml:"http://schemas.android.com/apk/res/android permission,attr"`
	TargetActivity string          `xml:"http://schemas.android.com/apk/res/android targetActivity,attr"`
	IntentFilters  *[]IntentFilter `xml:"intent-filter"`
	MetaData       *[]MetaData     `xml:"meta-data"`
}

// Service declares a service as one of the application's components.
type Service struct {
	Description     *string         `xml:"http://schemas.android.com/apk/res/android description,attr"`
	DirectBootAware *bool           `xml:"http://schemas.android.com/apk/res/android directBootAware,attr"`
	Enabled         *bool           `xml:"http://schemas.android.com/apk/res/android enabled,attr"`
	Exported        bool            `xml:"http://schemas.android.com/apk/res/android exported,attr"`
	IsolatedProcess *bool           `xml:"http://schemas.android.com/apk/res/android isolatedProcess,attr"`
	Label           *string         `xml:"http://schemas.android.com/apk/res/android label,attr"`
	Name            string          `xml:"http://schemas.android.com/apk/res/android name,attr"`
	Permission      *string         `xml:"http://schemas.android.com/apk/res/android permission,attr"`
	IntentFilters   *[]IntentFilter `xml:"intent-filter"`
	MetaData        *[]MetaData     `xml:"meta-data"`
}

// Receiver declares a broadcast receiver as one of the application's components.
type Receiver struct {
	DirectBootAware *bool           `xml:"http://schemas.android.com/apk/res/android directBootAware,attr"`
	Enabled         bool            `xml:"http://schemas.android.com/apk/res/android enabled,attr"`
	Exported        *bool           `xml:"http://schemas.android.com/apk/res/android exported,attr"`
	Label           *string         `xml:"http://schemas.android.com/apk/res/android label,attr"`
	Name            string          `xml:"http://schemas.android.com/apk/res/android name,attr"`
	IntentFilters   *[]IntentFilter `xml:"intent-filter"`
	MetaData        *[]MetaData     `xml:"meta-data"`
}

// Profileable specifies how profilies can access the application.
type Profileable struct {
	Shell   *bool `xml:"http://schemas.android.com/apk/res/android shell,attr"`
	Enabled *bool `xml:"http://schemas.android.com/apk/res/android enabled,attr"`
}

// Provider declares a content provider component which supplies structured access to data managed
// by the application.
type Provider struct {
	Authorities         string          `xml:"http://schemas.android.com/apk/res/android authorities,attr"`
	Enabled             *bool           `xml:"http://schemas.android.com/apk/res/android enabled,attr"`
	DirectBootAware     *bool           `xml:"http://schemas.android.com/apk/res/android directBootAware,attr"`
	Exported            bool            `xml:"http://schemas.android.com/apk/res/android exported,attr"`
	GrantURIPermissions *bool           `xml:"http://schemas.android.com/apk/res/android grantUriPermissions,attr"`
	InitOrder           *int32          `xml:"http://schemas.android.com/apk/res/android initOrder,attr"`
	Lable               *string         `xml:"http://schemas.android.com/apk/res/android label,attr"`
	MultiProcess        *bool           `xml:"http://schemas.android.com/apk/res/android multiprocess,attr"`
	Name                string          `xml:"http://schemas.android.com/apk/res/android name,attr"`
	Permission          *string         `xml:"http://schemas.android.com/apk/res/android permission,attr"`
	Process             *string         `xml:"http://schemas.android.com/apk/res/android process,attr"`
	ReadPermission      *string         `xml:"http://schemas.android.com/apk/res/android readPermission,attr"`
	Syncable            *bool           `xml:"http://schemas.android.com/apk/res/android syncable,attr"`
	WritePermission     *string         `xml:"http://schemas.android.com/apk/res/android writePermission,attr"`
	MetaData            *[]MetaData     `xml:"meta-data"`
	IntentFilters       *[]IntentFilter `xml:"intent-filter"`
}

// Query specifies the set of other apps that an app intends to interact with. These other apps can
// be specified by package name, by intent signature, or by provider authority.
type Query struct {
	Packages  *[]Package      `xml:"package"`
	Intents   *[]IntentFilter `xml:"intent"`
	Providers *[]Provider     `xml:"provider"`
}

// Package specifies a single app that an app intends to access. This other app might integrate with
// said app, or said app might use services the other app provides.
type Package struct {
	Name string `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

// SupportsScreen specifies the screen sizes an application supports and whether screen
// compatibility mode is enabled for screens larger than what the application supports.
type SupportsScreen struct {
	Resizeable              *bool  `xml:"http://schemas.android.com/apk/res/android resizeable,attr"`
	SmallScreens            *bool  `xml:"http://schemas.android.com/apk/res/android smallScreens,attr"`
	NormalScreens           *bool  `xml:"http://schemas.android.com/apk/res/android normalScreens,attr"`
	LargeScreens            *bool  `xml:"http://schemas.android.com/apk/res/android largeScreens,attr"`
	XLargeScreens           *bool  `xml:"http://schemas.android.com/apk/res/android xlargeScreens,attr"`
	AnyDensity              *bool  `xml:"http://schemas.android.com/apk/res/android anyDensity,attr"`
	RequiresSmallestWidthDP *int32 `xml:"http://schemas.android.com/apk/res/android requiresSmallestWidthDp,attr"`
	CompatibleWidthLimitDP  *int32 `xml:"http://schemas.android.com/apk/res/android compatibleWidthLimitDp,attr"`
	LargestWidthLimitDP     *int32 `xml:"http://schemas.android.com/apk/res/android largestWidthLimitDp,attr"`
}

// UsesPermission specifies a system permission that the user must grant for the application to
// operate correctly.
type UsesPermission struct {
	Name          string `xml:"http://schemas.android.com/apk/res/android name,attr"`
	MaxSDKVersion *int32 `xml:"http://schemas.android.com/apk/res/android maxSdkVersion,attr"`
}

// UsesSDK is an expression of an application's compatibility with one or more versions of the
// Android platform by means of an API level integer.
type UsesSDK struct {
	MinSDKVersion    *uint `xml:"http://schemas.android.com/apk/res/android minSdkVersion,attr"`
	TargetSDKVersion *uint `xml:"http://schemas.android.com/apk/res/android targetSdkVersion,attr"`
	MaxSDKVersion    *uint `xml:"http://schemas.android.com/apk/res/android maxSdkVersion,attr"`
}
