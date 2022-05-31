package main

import (
	"encoding/xml"
	"os"

	"github.com/accrescent/apkstat"
)

func main() {
	apk, err := apkstat.OpenAPK("accrescent.apk")
	if err != nil {
		panic(err)
	}

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "\t")
	if err := enc.Encode(apk.Manifest()); err != nil {
		panic(err)
	}
}
