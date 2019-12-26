package upstream

import (
	"fmt"
	"regexp"
	"strings"
)

func replaceRegexp(toMatch string, toReplaceIn string, re *regexp.Regexp) string {
	newPath := toReplaceIn

	if re == nil {
		return newPath
	}

	matched := re.FindStringSubmatch(toMatch)
	if len(matched) > 0 {
		for index, part := range matched[1:] {
			newPath = strings.Replace(newPath, fmt.Sprintf("$%d", index+1), part, -1)
		}
	}
	return newPath
}
