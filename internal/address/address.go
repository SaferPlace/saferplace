package address

import "errors"

var ErrUnresolved = errors.New("address: unresolved")

type Resolver interface {
	// Resolve converts the provided address to extended address version, x and
	// y coordinate, or returns an error if there was a problem resolving
	// the address.
	Resolve(address string) (string, float64, float64, error)
}
