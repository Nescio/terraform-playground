package sdk

import (
	"regexp"
)

// A regular expression used to validate common string ID patterns.
var reID = regexp.MustCompile(`^[a-zA-Z0-9\-\._]+$`)

// validID checks it the given input is present and matches the ID pattern.
func validID(v *string) bool {
	return v != nil && reID.MatchString(*v)
}

// validString checks if the given input is present and non-empty.
func validString(v *string) bool {
	return v != nil && *v != ""
}

// validInteger checks if the given input is present and greater than zero
func validInt(v *int) bool {
	return v != nil && *v > 0
}
