package helpers

import (
	"strings"
)

// A tag pode conter "name,omitempty" ou apenas "name". A função remove o omitempty
func ExtractJSONTag(tag string) string {
	
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}
	return tag 
}