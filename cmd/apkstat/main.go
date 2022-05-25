package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/accrescent/apkstat"
)

func main() {
	fmt.Println("Reading binary resources.arsc...")
	res, err := os.Open("resources.arsc")
	if err != nil {
		panic(err)
	}
	defer res.Close()
	r, err := apkstat.NewResTable(res)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Reading binary AndroidManifest.xml...\n\n")
	bxml, err := os.Open("AndroidManifest.xml")
	if err != nil {
		panic(err)
	}
	defer bxml.Close()
	x, err := apkstat.NewXMLFile(bxml, r, nil)
	if err != nil {
		panic(err)
	}

	var m apkstat.Manifest
	if err := xml.Unmarshal([]byte(x.String()), &m); err != nil {
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
	fmt.Println("label:", m.Application.Label)
}
