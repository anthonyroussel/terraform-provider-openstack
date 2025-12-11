package verify

import (
	"fmt"
	"time"
)

func ValidateDate(v any, k string) (ws []string, errors []error) {
	value := v.(string)

	if _, err := time.Parse(time.RFC3339, value); err != nil {
		errors = append(errors, fmt.Errorf(
			"%q is not a valid RFC3339 timestamp for %q: %w",
			value,
			k,
			err,
		))
	}

	return
}
