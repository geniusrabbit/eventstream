//go:build nats || kafka || natstream || redistream || allstorage || allsource || all
// +build nats kafka natstream redistream allstorage allsource all

package storage

import (
	"fmt"

	storenc "github.com/geniusrabbit/eventstream/storage"
)

func ncStreamOpts(streamName string, p StreamParams) ([]any, error) {
	if len(p.Targets) == 0 {
		return nil, fmt.Errorf("stream %q: targets required", streamName)
	}
	targets := make([]storenc.Target, 0, len(p.Targets))
	for _, t := range p.Targets {
		targets = append(targets, storenc.Target{Fields: t.Fields, Where: t.Where})
	}
	return []any{
		storenc.WithID(streamName),
		storenc.WithTargets(targets...),
	}, nil
}
