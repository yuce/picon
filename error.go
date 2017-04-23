package picon

import "errors"

var (
	errNotConnected = errors.New("You must :connect to a server")
	errNoIndex      = errors.New("You must :use an index")
)
