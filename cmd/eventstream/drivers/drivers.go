// Package drivers blank-imports cmd-side source/storage driver registration.
// Which drivers are linked is controlled by build tags on
// cmd/eventstream/{source,storage}/register_*.go (-tags all, nats, …).
// Library driver packages themselves are always available via import.
package drivers

import (
	_ "github.com/geniusrabbit/eventstream/cmd/eventstream/source"
	_ "github.com/geniusrabbit/eventstream/cmd/eventstream/storage"
)
