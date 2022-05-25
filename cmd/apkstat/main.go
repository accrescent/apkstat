package main

import (
	"fmt"

	"github.com/accrescent/apkstat"
)

func main() {
	apk, err := apkstat.OpenAPK("accrescent.apk")
	if err != nil {
		panic(err)
	}
	m := apk.Manifest()

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
