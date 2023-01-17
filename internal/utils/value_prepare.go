package utils

import (
	"os"
	"regexp"
	"strings"
)

var reVariables = regexp.MustCompile(`\{\{\s*([^} \t\n]+)\s*\}\}`)

// PrepareValue from value formated string
func PrepareValue(val string, vars ...map[string]string) string {
	for _, varNameExp := range reVariables.FindAllStringSubmatch(val, -1) {
		var (
			varName = varNameExp[1]
			varVal  = varNameExp[0]
		)
		if strings.HasPrefix(varName, "@env:") {
			varVal = os.Getenv(varName[5:])
		} else {
			for _, mVars := range vars {
				if len(mVars) != 0 {
					varVal = mVars[varName]
				}
			}
		}
		val = strings.ReplaceAll(val, varNameExp[0], varVal)
	}
	return val
}
