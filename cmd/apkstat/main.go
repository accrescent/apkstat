package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/accrescent/apkstat"
)

func main() {
	fmt.Printf("Reading binary AndroidManifest.xml...\n\n")
	file, err := os.Open("AndroidManifest.xml")
	if err != nil {
		panic(err)
	}

	p, err := apkstat.NewParser(file)
	if err != nil {
		panic(err)
	}

	var m apkstat.Manifest
	if err := xml.Unmarshal([]byte(p.String()), &m); err != nil {
		panic(err)
	}

	fmt.Println("package:", m.Package)
	fmt.Println("versionCode:", m.VersionCode)
	fmt.Println("versionName:", m.VersionName)

	for i := 0; i < len(m.UsesPermissions); i++ {
		fmt.Println("permission:", m.UsesPermissions[i].Name)
	}

	fmt.Println("minSdkVersion:", m.UsesSDK.MinSDKVersion)
	fmt.Println("targetSdkVersion:", m.UsesSDK.TargetSDKVersion)
	fmt.Println("maxSdkVersion:", m.UsesSDK.MaxSDKVersion)
}
