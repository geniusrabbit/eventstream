// Package converter provides message encode/decode formats (JSON, BSON, raw).
package converter

import (
	"encoding/json"

	"gopkg.in/mgo.v2/bson"

	"github.com/pkg/errors"
)

var (
	errUnsupportedDecodeType = errors.New(`unsupported converter decode type`)
	errUnsupportedEncodeType = errors.New(`unsupported converter encode type`)
)

// Converter marshals and unmarshals message payloads.
type Converter interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

// Fnk adapts encoder/decoder functions to the Converter interface.
type Fnk struct {
	name    string
	encoder func(v any) ([]byte, error)
	decoder func(data []byte, v any) error
}

func (f Fnk) String() string {
	return f.name
}

// Marshal encodes v to bytes.
func (f Fnk) Marshal(v any) ([]byte, error) {
	return f.encoder(v)
}

// Unmarshal decodes data into v.
func (f Fnk) Unmarshal(data []byte, v any) error {
	return f.decoder(data, v)
}

// Built-in converters.
var (
	JSON Converter = Fnk{name: "json", encoder: json.Marshal, decoder: json.Unmarshal}
	BSON Converter = Fnk{name: "bson", encoder: bson.Marshal, decoder: bson.Unmarshal}
	RAW  Converter = Fnk{
		name: "raw",
		encoder: func(v any) ([]byte, error) {
			switch b := v.(type) {
			case []byte:
				return b, nil
			case string:
				return []byte(b), nil
			}
			return nil, errors.Wrapf(errUnsupportedEncodeType, "[raw] %T", v)
		},
		decoder: func(data []byte, v any) error {
			switch pt := v.(type) {
			case *any:
				*pt = data
			case *[]byte:
				*pt = data
			case []byte:
				copy(pt, data)
			}
			return errors.Wrapf(errUnsupportedDecodeType, "[raw] %T", v)
		},
	}
)

// ByName returns a named converter (json, bson, or raw by default).
func ByName(name string) Converter {
	switch name {
	case "bson":
		return BSON
	case "json":
		return JSON
	}
	return RAW
}
