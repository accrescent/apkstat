## apkstat

An APK parsing tool and library for Go.

## Usage

### CLI

`apkstat` is a basic CLI tool for printing APK manifests and binary XML files.

```
Usage of apkstat:
  -apk string
        APK to print binary XML from
  -xml string
        binary XML file to print (Android manifest is default)
```

`-apk` must be specified. If `-xml` is specified, apkstat will attempt to print
that file in the APK ZIP hierarchy. If it isn't, apkstat will pretty print the
Android manifest.

### Library

The main entry point for apkstat is the APK type, which you can create an
instance of with the `apk.Open` and `apk.OpenWithConfig` functions.

If you need to do lower-level parsing (which is usually unnecessary), you can
open resource tables and Android binary XML files directly with `NewResTable()`
and `NewXMLFile()` respectively.

Example usage:

```go
package main

import (
	"fmt"

	"github.com/accrescent/apkstat"
)

func main() {
	apk, err := apk.Open("accrescent.apk")
	if err != nil {
		panic(err)
	}

	manifest := apk.Manifest()

	fmt.Println("App ID:", manifest.Package)
	fmt.Println("App version code:", manifest.VersionCode)
	fmt.Println("App version name:", manifest.VersionName)
	for _, p := range *manifest.UsesPermissions {
		fmt.Println("Requested permission:", p.Name)
	}
}
```

## License

apkstat is licensed under the ISC license. However, parts of it are based on
code from the Android Open Source Project and the androidbinary library by
Ichinose Shogo which are licensed under the Apache 2.0 and MIT licenses
respectively.
