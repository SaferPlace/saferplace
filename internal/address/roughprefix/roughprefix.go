// Package roughprefix resolves the very rough address based on the eircode
// prefix code, as well as other well known prefixes (like city names)
package roughprefix

import (
	"strings"

	"safer.place/data"
	"safer.place/internal/address"
)

type Resolver struct {
	prefixes map[string]data.PrefixCoordinates
}

func New() *Resolver {
	return &Resolver{
		prefixes: data.RoughPrefixCoordinates(),
	}
}

// Roughly resolve the address
func (r *Resolver) Resolve(addr string) (string, float64, float64, error) {
	for prefix, data := range r.prefixes {
		if strings.HasPrefix(
			strings.ToLower(addr),
			strings.ToLower(prefix),
		) {
			return "~ " + strings.Join(data.Names, ", "), data.X, data.Y, nil
		}
	}

	return "", 0, 0, address.ErrUnresolved
}
