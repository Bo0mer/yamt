package internal

import "strconv"

// ErrParser extends strconv functions for integers, adding support for lazy
// error checking.
// All conversions are done using base 10, assuming 64 bit integers.
type ErrParser struct {
	err error
}

// ParseInt extends strconv.ParseInt by preserving last occurred error.
func (p *ErrParser) ParseInt(s string) int {
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		p.err = err
	}
	return int(i64)
}

// ParseUint64 extends strconv.ParseUint64 by preserving last occurred error.
func (p *ErrParser) ParseUint64(s string) uint64 {
	u, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		p.err = err
	}
	return u
}

// Err returns the last error encountered by the parser.
func (p *ErrParser) Err() error {
	return p.err
}
