package escape

import "strings"

var escapedCharters = []string{"\\", "_", "%"}

func Like(value string) string {
	for _, charter := range escapedCharters {
		if strings.Count(value, charter) > 0 {
			value = strings.ReplaceAll(value, charter, "\\"+charter)
		}
	}
	return value
}

func LikeWithEscape(value string, escape string) string {
	for _, charter := range escapedCharters {
		if strings.Count(value, charter) > 0 {
			value = strings.ReplaceAll(value, charter, escape+charter)
		}
	}
	return value
}
