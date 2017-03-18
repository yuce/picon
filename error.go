package picon

import "errors"

var (
	errNotConnected = errors.New("You must :connect to a server")
	errNoDatabase   = errors.New("You must :use a database")
)
