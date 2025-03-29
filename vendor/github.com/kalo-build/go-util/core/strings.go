package core

import "strings"

func LinesToString(allLines []string) (string, error) {
	var builder strings.Builder
	for _, structLine := range allLines {
		_, writeErr := builder.WriteString(structLine + "\n")
		if writeErr != nil {
			return "", writeErr
		}
	}
	return builder.String(), nil
}
