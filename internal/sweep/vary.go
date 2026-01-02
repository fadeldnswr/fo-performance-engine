package sweep

import (
	"errors"
	"strconv"
	"strings"
)

// Define struct to hold variation information
type Variation struct {
	Field  string
	Values []float64
}

// Define function to parse variations from a configuration (placeholder)
func ParseVariations(spec string) (Variation, error) {
	// Parse the specification string to extract field and values
	parts := strings.SplitN(spec, "=", 2)
	if len(parts) != 2 { return Variation{}, errors.New("")}

	// Separate key and value parts
	key := strings.TrimSpace(parts[0])
	raw := strings.Split(parts[1], ",")
	if len(raw) == 0 { return Variation{}, errors.New("No values in raw data") }

	// Split values and convert to float64
	vals := make([]float64, 0, len(raw))
	for _, s := range raw {
		s = strings.TrimSpace(s)
		value, err := strconv.ParseFloat(s,64)
		if err != nil {
			return Variation{}, errors.New("Raw data has bad value")
		}
		vals = append(vals, value)
	}
	return Variation{Field: key, Values: vals}, nil
}