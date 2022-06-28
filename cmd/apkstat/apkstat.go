package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"

	"github.com/accrescent/apkstat"
	"github.com/accrescent/apkstat/schemas"
)

func main() {
	apkFlag := flag.String("apk", "", "APK to print binary XML from")
	xmlFlag := flag.String("xml", "", "binary XML file to print (Android manifest is default)")
	xmlResFlag := flag.String(
		"xmlres",
		"",
		"well-known XML resource to print. Must be one of 'network-security'",
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
	defer apk.Close()
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "    ")

	if *xmlFlag != "" {
		xmlFile, err := apk.OpenXML(*xmlFlag)
		if err != nil {
			fatal(err.Error())
		}

		fmt.Println(xmlFile.String())
	} else if *xmlResFlag != "" {
		manifest := apk.Manifest()
		switch *xmlResFlag {
		case "network-security":
			xmlRes := manifest.Application.NetworkSecurityConfig
			if xmlRes == nil {
				fatal("network security config not present")
			} else {
				xmlFile, err := apk.OpenXML(*xmlRes)
				if err != nil {
					fatal(err.Error())
				}

				var nsConfig schemas.NetworkSecurityConfig
				if err := xml.Unmarshal([]byte(xmlFile.String()), &nsConfig); err != nil {
					fatal(err.Error())
				}
				if err := enc.Encode(nsConfig); err != nil {
					fatal(err.Error())
				}
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
