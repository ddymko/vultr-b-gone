package util

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// OptionsInterface
type OptionsInterface interface {
	Validate() error
}

// OptionsScheme
type OptionsScheme struct {
	Method      string
	Regions     []string
	OmitRegions []string
}

// SetOptions
func (o *OptionsScheme) SetOptions(cobra *cobra.Command) error {
	var err error

	if o.Method, err = cobra.Flags().GetString("method"); err != nil {
		return err
	}

	if o.Regions, err = cobra.Flags().GetStringSlice("regions"); err != nil {
		return err
	}

	if o.OmitRegions, err = cobra.Flags().GetStringSlice("omit-regions"); err != nil {
		return err
	}

	return nil
}

// Validate
func (o *OptionsScheme) Validate() error {
	method := o.Method
	regions := o.Regions
	omitRegions := o.OmitRegions

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

		if len(o.Regions) > 0 {
			temp := make(map[string]bool)
			for _, v := range idk[0] {
				if !RegionsMap[strings.ToLower(v)] {
					return fmt.Errorf("%s is a invalid region", v)
				}
				temp[v] = true
			}
			RegionsMap = temp
		} else {
			for _, v := range idk[0] {
				if !RegionsMap[strings.ToLower(v)] {
					return fmt.Errorf("%s is a invalid region", v)
				}
				delete(RegionsMap, v)
			}
		}
	}
	return nil
}
