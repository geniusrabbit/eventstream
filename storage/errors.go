// Package storage contains shared publish-storage helpers and errors used by
// queue drivers (storage/nats, storage/kafka, …) and by specialty storages.
package storage

import "github.com/pkg/errors"

// ErrInvalidOption is returned when Stream is called with an unsupported option.
var ErrInvalidOption = errors.New(`invalid option`)
