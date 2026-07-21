package eventstream

import (
	"github.com/geniusrabbit/eventstream/stream"
)

// Streamer is an alias of stream.Streamer for the root package API.
//
//go:generate mockgen -source $GOFILE -package mocks -destination internal/mocks/stream.go
type Streamer = stream.Streamer
