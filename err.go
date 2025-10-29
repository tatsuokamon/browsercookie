package browsercookie

import "errors"

var (
	ErrNotImplemented error = errors.New("Not Implemented yet")

	// Safari
	ErrSafariOnlyOnOSX error = errors.New("Safari is only available on OSX")
)
