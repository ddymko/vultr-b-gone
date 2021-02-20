package cmd

import (
	"errors"
	"fmt"
	"strings"
)

//todo we should get a set of regions that we want to delete if `region` is passed on regardless if its region or omit

// SchemeInterface
type SubscriptionSchemeInterface interface {
	Validate() error
}

// Scheme
type SubscriptionScheme struct {
	Method      string
	Regions     []string
	OmitRegions []string
}

// Validate
func (s *SubscriptionScheme) Validate() error {
	method := s.Method
	regions := s.Regions
	omitRegions := s.OmitRegions

	if method != "region" && method != "all" {
		return fmt.Errorf("%s is invalid please provide 'region' or 'all'", method)
	}

	if method == "region" {
		var idk [][]string
		if len(regions) > 0 {
			idk = append(idk, regions)
		}

		if len(omitRegions) > 0 {
			idk = append(idk, omitRegions)
		}

		if len(idk) == 2 || len(idk) == 0 {
			return errors.New("please provide either regions or omit-regions")
		}

		regions := []string{
			"EWR",
			"ORD",
			"DFW",
			"SEA",
			"LAX",
			"ATL",
			"AMS",
			"LHR",
			"FRA",
			"SJC",
			"SYD",
			"YTO",
			"CDG",
			"NRT",
			"ICN",
			"MIA",
			"SGP",
		}

		exists := struct {
			Name   string
			Exists bool
		}{}

		for _, v := range idk[0] {
			for _, region := range regions {
				if strings.ToUpper(v) == region {
					exists.Exists = true
					break
				}
				exists.Exists = false
				exists.Name = v
			}
			if !exists.Exists {
				break
			}
		}

		if !exists.Exists {
			return fmt.Errorf("%s is a invalid region", exists.Name)
		}
	}

	return nil
}
