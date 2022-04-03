package main

import (
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

	m, err := apkstat.ParseManifest(file)
	if err != nil {
		panic(err)
	}

	fmt.Println("package:", m.Package)
	fmt.Println("versionCode:", m.VersionCode)
	fmt.Println("versionName:", m.VersionName)

	for i := 0; i < len(m.UsesPermissions); i++ {
		fmt.Println("permission:", m.UsesPermissions[i].Name)
	}
}
