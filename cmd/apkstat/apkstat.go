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
	xmlResFlag := flag.String(
		"xmlres",
		"",
		"well-known XML resource to print. Must be one of 'network-security' or 'extraction-rules'",
	)
	flag.Parse()

	if *apkFlag == "" {
		fatal("must supply APK parameter")
	}
	if *xmlFlag != "" && *xmlResFlag != "" {
		fatal("-xml and -xmlres are mutually exclusive")
	}

	apk, err := apk.Open(*apkFlag)
	if err != nil {
		fatal(err.Error())
	}
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "    ")

	if *xmlFlag != "" {
		xmlFile, err := apk.OpenXML(*xmlFlag)
		if err != nil {
			fatal(err.Error())
		}

		fmt.Println(xmlFile.String())
	} else if *xmlResFlag != "" {
		switch *xmlResFlag {
		case "extraction-rules":
			rules, err := apk.DataExtractionRules()
			if err != nil {
				fatal(err.Error())
			}
			if err := enc.Encode(rules); err != nil {
				fatal(err.Error())
			}
		case "network-security":
			nsConfig, err := apk.NetworkSecurityConfig()
			if err != nil {
				fatal(err.Error())
			}
			if err := enc.Encode(nsConfig); err != nil {
				fatal(err.Error())
			}
		default:
			fatal("xmlres '" + *xmlResFlag + "' not valid")
		}
		fmt.Println()
	} else {
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
