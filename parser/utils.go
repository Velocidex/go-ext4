package parser

import "strings"

func GetComponents(path string) []string {
	path = strings.ReplaceAll(path, "/", "\\")
	result := []string{}
	for _, c := range strings.Split(path, "\\") {
		if c != "" {
			result = append(result, c)
		}
	}
	return result
}

func DlvBreak() {
}
