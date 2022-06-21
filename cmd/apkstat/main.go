package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"

	"github.com/accrescent/apkstat"
)

func main() {
	apkFlag := flag.String("apk", "", "APK to print manifest of")
	xmlFlag := flag.String("xml", "", "binary XML to print as text")
	flag.Parse()

	if *apkFlag == "" && *xmlFlag == "" || *apkFlag != "" && *xmlFlag != "" {
		fatal("must supply either APK or binary XML")
	}

	if *apkFlag != "" {
		apk, err := apkstat.OpenAPK(*apkFlag)
		if err != nil {
			fatal(err.Error())
		}

		enc := xml.NewEncoder(os.Stdout)
		enc.Indent("", "\t")
		if err := enc.Encode(apk.Manifest()); err != nil {
			fatal(err.Error())
		}
		fmt.Println()
	} else {
		file, err := os.Open(*xmlFlag)
		if err != nil {
			fatal(err.Error())
		}
		xmlFile, err := apkstat.NewXMLFile(file, nil, nil)
		if err != nil {
			fatal(err.Error())
		}

		fmt.Println(xmlFile.String())
	}
}

func fatal(err string) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
	flag.Usage()
	os.Exit(1)
}
