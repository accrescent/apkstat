package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"

	"github.com/accrescent/apkstat"
)

func main() {
	apkFlag := flag.String("apk", "", "APK to print binary XML from")
	xmlFlag := flag.String("xml", "", "binary XML file to print (Android manifest is default)")
	flag.Parse()

	if *apkFlag == "" {
		fatal("must supply APK parameter")
	}

	apk, err := apk.Open(*apkFlag)
	if err != nil {
		fatal(err.Error())
	}
	defer apk.Close()

	if *xmlFlag != "" {
		xmlFile, err := apk.OpenXML(*xmlFlag)
		if err != nil {
			fatal(err.Error())
		}

		fmt.Println(xmlFile.String())
	} else {
		enc := xml.NewEncoder(os.Stdout)
		enc.Indent("", "    ")
		if err := enc.Encode(apk.Manifest()); err != nil {
			fatal(err.Error())
		}
		fmt.Println()
	}
}

func fatal(err string) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
	flag.Usage()
	os.Exit(1)
}
