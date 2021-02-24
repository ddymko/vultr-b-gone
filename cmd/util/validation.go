package util

import (
	"fmt"
	"os"
	"strings"
)

var RegionsMap = map[string]bool{
	"ewr": true,
	"ord": true,
	"dfw": true,
	"sea": true,
	"lax": true,
	"atl": true,
	"ams": true,
	"lhr": true,
	"fra": true,
	"sjc": true,
	"syd": true,
	"yto": true,
	"cdg": true,
	"nrt": true,
	"icn": true,
	"mia": true,
	"sgp": true,
}

func CheckError(err error) {
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}
}

func LocationCheck(location string) bool {
	return RegionsMap[strings.ToLower(location)]
}
