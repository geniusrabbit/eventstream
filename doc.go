// Package eventstream defines the core interfaces of an event processing
// pipeline: sources that ingest messages, storages that open sinks, and
// streamers that process events.
//
// Concrete drivers and the wiring Engine live in subpackages
// (source/…, storage/…, pkg/pipeline).
package eventstream
