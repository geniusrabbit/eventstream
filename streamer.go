//
// @project geniusrabbit::eventstream 2017, 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2020
//

package eventstream

import (
	"github.com/geniusrabbit/eventstream/stream"
)

// StreamConfig of the stream
type StreamConfig = stream.Config

// Streamer interface of data processing describes
// basic methods of data pipeline
//go:generate mockgen -source $GOFILE -package mocks -destination internal/mocks/stream.go
type Streamer = stream.Streamer
