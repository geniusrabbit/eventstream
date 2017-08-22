//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package errors

import "errors"

// List of common system errors
var (
	ErrConnectionIsNotDefined = errors.New(`"connect" field is not defined`)
)
